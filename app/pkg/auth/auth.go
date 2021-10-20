package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// Encrypt 使用 bcrypt 加密纯文本
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare 验证加密的文本是否与纯文本相同
func Compare(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

type JwtOptions jwt.StandardClaims

type JwtAuthClaims struct {
	Guard string `json:"guard"` // 授权守卫
	jwt.StandardClaims
}

// SignJwtToken 生成 JWT 令牌
func SignJwtToken(guard string, secret string, ops *JwtOptions) string {
	claims := JwtAuthClaims{
		Guard: guard,
		StandardClaims: jwt.StandardClaims{
			Audience:  ops.Audience,
			ExpiresAt: ops.ExpiresAt,
			Id:        ops.Id,
			IssuedAt:  ops.IssuedAt,
			Issuer:    ops.Issuer,
			NotBefore: ops.NotBefore,
			Subject:   ops.Subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString
}

// VerifyJwtToken 验证 Token
func VerifyJwtToken(token string, secret string) (*JwtAuthClaims, error) {
	claims := &JwtAuthClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	return claims, err
}

// GetAuthUserID 获取授权登录的用户ID
func GetAuthUserID(c *gin.Context) int {
	return c.GetInt(entity.LoginUserID)
}

// GetJwtToken 获取登录授权 token
func GetJwtToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	token = strings.TrimLeft(token, "Bearer")
	token = strings.TrimSpace(token)

	// Headers 中没有授权信息则读取 url 中的 token
	if len(token) == 0 {
		token = c.DefaultQuery("token", "")
	}

	if len(token) == 0 {
		token = c.DefaultPostForm("token", "")
	}

	return token
}