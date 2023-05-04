package router

import (
	"github.com/gin-gonic/gin"
	"mall/cmd/Account/config"
	myErrors "mall/cmd/Account/errors"
	"mall/cmd/Account/service"
	"mall/pkg/util/errors"
	"net/http"
	"sync"
)

// 前端请求参数
const (
	accountNameParam = "accountName"
	passwordParam    = "password"
)

type WebRouter struct {
	mu sync.RWMutex
	eg *gin.Engine
}

func NewWebRouter() *WebRouter {
	return &WebRouter{}
}

func (r *WebRouter) StartToServe() {
	eg := gin.Default()
	if eg == nil {
		panic("gin engine initial failed")
	}

	r.mu.Lock()
	r.eg = eg
	r.mu.Unlock()

	r.AddRoute(RegisterURL)
	r.AddRoute(UnRegisterURL)
	r.AddRoute(UpdatePassword)
	r.AddRoute(LoginURL)
	r.AddRoute(LogoutURL)

	eg.Run(config.GetGConfig().GetServiceHttpPort())
}

func (r *WebRouter) AddRoute(URLPath string) {
	switch URLPath {
	case RegisterURL:
		r.eg.POST(URLPath, RegisterHandler)
	case UnRegisterURL:
		r.eg.POST(URLPath, UnRegisterHandler)
	case UpdatePassword:
		r.eg.POST(URLPath, UpdatePasswordHandler)
	case LoginURL:
		r.eg.POST(URLPath, LoginHandler)
	case LogoutURL:
		r.eg.POST(URLPath, LogoutHandler)
	}
}

func BuildHttpResponse(context *gin.Context, httpStatusCode int, err errors.Code) {
	context.JSON(httpStatusCode, gin.H{
		"code": err.Code(),
		"msg":  err.Message(),
	})
}

func RegisterHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	password := context.PostForm(passwordParam)

	success := service.Register(accountName, password)
	if success {
		BuildHttpResponse(context, http.StatusOK, myErrors.RegisterSuccess)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.RegisterFail)
	}
}

func UnRegisterHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	success := service.UnRegister(accountName)
	if success {
		BuildHttpResponse(context, http.StatusOK, myErrors.UnRegisterSuccess)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UnRegisterFail)
	}
}

func UpdatePasswordHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	password := context.PostForm(passwordParam)
	success := service.ChangePassword(accountName, password)
	if success {
		BuildHttpResponse(context, http.StatusOK, myErrors.PasswordChangeSuccess)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.PasswordChangeFail)
	}
}

func LoginHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	password := context.PostForm(passwordParam)
	token, success := service.LoginUsingRpc(accountName, password)
	if success {
		context.JSON(http.StatusOK, gin.H{
			"code":  myErrors.LoginSuccess.Code(),
			"msg":   myErrors.LoginSuccess.Message(),
			"token": token,
		})
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.LoginFail)
	}
}

func LogoutHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	success := service.LogoutUsingRpc(accountName)
	if success {
		BuildHttpResponse(context, http.StatusOK, myErrors.LogoutSuccess)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.LogoutFail)
	}
}
