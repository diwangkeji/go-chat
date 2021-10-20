package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-chat/app/cache"
	"go-chat/app/repository"
	"go-chat/app/service"
	"go-chat/provider"
	"go-chat/testutil"
	"net/url"
	"testing"
)

func testUser() *User {
	config := testutil.GetConfig()
	db := provider.MysqlConnect(config)
	redisClient := testutil.TestRedisClient()

	userRepo := repository.UserRepository{DB: db}
	smsService := service.NewSmsService(&cache.SmsCodeCache{Redis: redisClient})

	return NewUserHandler(&userRepo, smsService)
}

func TestUser_Detail(t *testing.T) {
	a := testUser()

	r := testutil.NewTestRequest("/user/detail", func(context *gin.Context) {
		context.Set("__user_id__", 2054)
	}, a.Detail)

	value := &url.Values{}

	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestUser_ChangeDetail(t *testing.T) {
	a := testUser()

	r := testutil.NewTestRequest("/user/change/detail", func(context *gin.Context) {
		context.Set("__user_id__", 2054)
	}, a.ChangeDetail)

	value := &url.Values{}
	value.Add("nickname", "返税款1")
	value.Add("gender", "1")

	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestUser_ChangePassword(t *testing.T) {
	a := testUser()

	r := testutil.NewTestRequest("/user/change/password", func(context *gin.Context) {
		context.Set("__user_id__", 2054)
	}, a.ChangePassword)

	value := &url.Values{}
	value.Add("old_password", "admin123")
	value.Add("new_password", "admin123")

	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}