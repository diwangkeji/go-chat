package request

// 登录接口验证请求
type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// 注册接口验证请求
type RegisterRequest struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=30" label:"用户昵称"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,len=11,phone" label:"手机号"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16" label:"密码"`
	SmsCode  string `form:"sms_code" json:"sms_code" binding:"required,len=6,numeric" label:"验证码"`
	Platform string `form:"platform" json:"platform" binding:"required,oneof=h5 ios windows mac web" label:"登录平台"`
}