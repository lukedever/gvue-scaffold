package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lukedever/api"
)

var (
	errInvalidUser   = errors.New("用户不存在")
	errEmailExist    = errors.New("邮箱用户已存在")
	errEmailNotExist = errors.New("邮箱用户不存在")
	errWrongPassword = errors.New("密码错误")
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResp struct {
	Token string    `json:"token"`
	User  *api.User `json:"user"`
}

// HandleLogin handle '/login' route
func (s *Server) HandleLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.respondWithErr(c, err)
		return
	}

	user, _ := s.UserService.FindUserByKV("email", req.Email)
	if user.ID == 0 {
		s.respondWithErr(c, errEmailNotExist)
		return
	}

	if md5Str(req.Password) != user.Password {
		s.respondWithErr(c, errWrongPassword)
		return
	}

	token, err := s.genToken(user.ID, user.Name)
	if err != nil {
		s.respondWithServerErr(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResp{
		Token: token,
		User:  user,
	})
}

type registerRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6,max=12"`
	RePassword string `json:"repassword" binding:"required,eqfield=Password"`
}

// HandleRegister handle '/register' route
func (s *Server) HandleRegister(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.respondWithErr(c, err)
		return
	}

	user, _ := s.UserService.FindUserByKV("email", req.Email)
	if user.ID > 0 {
		s.respondWithErr(c, errEmailExist)
		return
	}

	ss := strings.Split(req.Email, "@")
	u := api.User{
		Name:     ss[0],
		Email:    req.Email,
		Password: md5Str(req.Password),
	}
	if err := s.UserService.CreateUser(&u); err != nil {
		s.respondWithServerErr(c, err)
		return
	}

	c.JSON(http.StatusOK, u)
}

// HandleProfile handle '/profile' route
func (s *Server) HandleProfile(c *gin.Context) {
	id, _ := c.Get("user_id")
	user, _ := s.UserService.FindUserByKV("id", id)
	if user.ID == 0 {
		s.respondWithAuthErr(c, errInvalidUser)
		return
	}

	c.JSON(http.StatusOK, user)
}
