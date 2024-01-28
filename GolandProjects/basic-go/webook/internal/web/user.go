// 与http打交道
package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

// 所有用户有关的路由
type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *service.UserService
}

func NewHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
	}

}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	{
		ug.POST("/login", h.Login)
		//server.PUT("/user/signup", h.Signup)
		ug.POST("/signup", h.SignUp)
		ug.GET("/profile", h.Profile)
		ug.POST("/edit", h.Edit)
	}

}

func (h *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := c.Bind(&req); err != nil {
		return
	}
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		c.String(http.StatusOK, "非法邮箱")
		return
	}
	if req.Password != req.ConfirmPassword {
		c.String(http.StatusOK, "两次输入密码不对")
		return
	}
	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		c.String(http.StatusOK, "密码必须包含数字，字母特殊字符，并且不少于8位")
		return
	}
	err = h.svc.Signup(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:
		c.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		c.String(http.StatusOK, "邮箱错误，请换一个")
	default:
		c.String(http.StatusOK, "系统错误")

	}

}

func (h *UserHandler) Login(c *gin.Context) {
	type Req struct {
		Email    string
		Password string
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(c, req.Email, req.Password)

	switch err {
	case nil:
		sess := sessions.Default(c)
		sess.Set("userId", u.ID)
		sess.Options(sessions.Options{
			MaxAge:   900,
			HttpOnly: true,
		})
		err = sess.Save()
		if err != nil {
			c.String(http.StatusOK, "系统错误")
			return
		}
		c.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		c.String(http.StatusOK, "用户名或者密码不对")

	default:
		c.String(http.StatusOK, "系统错误")

	}
	if err != nil {

		return
	}

}

func (h *UserHandler) Profile(c *gin.Context) {

	type ProfileReq struct {
		Id int64 `json:"id"`
	}
	var req ProfileReq
	if err := c.Bind(&req); err != nil {
		return
	}
	sess := sessions.Default(c)
	uc := sess.Get("userId")
	u, err := h.svc.Profile(c, domain.User{
		ID: uc.(int64),
	})
	println(err)
	switch err {
	case nil:
		c.JSON(http.StatusOK, &domain.User{
			NickName: u.NickName,
			Birthday: u.Birthday,
			AboutMe:  u.AboutMe,
		})
	default:
		c.String(http.StatusOK, "系统错误")

	}
}

func (h *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req EditReq
	if err := c.Bind(&req); err != nil {
		return
	}
	sess := sessions.Default(c)
	uc := sess.Get("userId")
	_, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		c.String(http.StatusOK, "生日格式不对")
		return
	}
	err = h.svc.Edit(c, domain.User{
		ID:       uc.(int64),
		NickName: req.Nickname,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	})
	switch err {
	case nil:
		c.String(http.StatusOK, "修改成功")
	default:
		c.String(http.StatusOK, "系统错误")

	}
}
