package router

import (
	"github.com/gin-gonic/gin"
	"mall/cmd/UserInfo/config"
	"mall/cmd/UserInfo/entity"
	myErrors "mall/cmd/UserInfo/errors"
	"mall/cmd/UserInfo/service"
	"mall/pkg/log"
	"mall/pkg/util/errors"
	"net/http"
	"sync"
)

// 前端请求参数
const (
	accountNameParam = "accountName"
	nickNameParam    = "nickName"
	portraitURLParam = "portraitURL"
	birthdayParam    = "birthday"
	phoneParam       = "phone"
	genderParam      = "gender"
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

	r.AddRoute(RootURL)
	r.AddRoute(CreateUserInfoURL)
	r.AddRoute(UpdateUserInfoURL)
	r.AddRoute(UpdateNickNameURL)
	r.AddRoute(UpdatePortraitURL)
	r.AddRoute(UpdateBirthdayURL)
	r.AddRoute(UpdatePhoneURL)
	r.AddRoute(UpdateGenderURL)
	r.AddRoute(DeleteUserInfoURL)
	r.AddRoute(QueryUserInfoURL)
	eg.Run(config.GetGConfig().GetServiceHttpPort())
}

func (r *WebRouter) AddRoute(URLPath string) {
	switch URLPath {
	case RootURL:
		r.eg.GET(URLPath, RecordClientInfoHandler)
	case CreateUserInfoURL:
		r.eg.POST(URLPath, CreateUserInfoHandler)
	case UpdateUserInfoURL:
		r.eg.POST(URLPath, UpdateUserInfoHandler)
	case UpdateNickNameURL:
		r.eg.POST(URLPath, UpdateNickNameHandler)
	case UpdatePortraitURL:
		r.eg.POST(URLPath, UpdatePortraitURLHandler)
	case UpdateBirthdayURL:
		r.eg.POST(URLPath, UpdateBirthdayHandler)
	case UpdatePhoneURL:
		r.eg.POST(URLPath, UpdatePhoneHandler)
	case UpdateGenderURL:
		r.eg.POST(URLPath, UpdateGenderHandler)
	case DeleteUserInfoURL:
		r.eg.POST(URLPath, DeleteUserInfoHandler)
	case QueryUserInfoURL:
		r.eg.POST(URLPath, QueryUserInfoURLHandler)
	}
}

func BuildHttpResponse(context *gin.Context, httpStatusCode int, err errors.Code) {
	context.JSON(httpStatusCode, gin.H{
		"code": err.Code(),
		"msg":  err.Message(),
	})
}

func RecordClientInfoHandler(context *gin.Context) {
	// 打印peer的相关信息
	log.PrintLog("Client addr info:%s", context.Request.RemoteAddr)
}

func CreateUserInfoHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	nickName := context.PostForm(nickNameParam)
	portraitURL := context.PostForm(portraitURLParam)
	birthday := context.PostForm(birthdayParam)
	phone := context.PostForm(phoneParam)
	gender := context.PostForm(genderParam)

	userInfo := entity.TbUserInfo{
		AccountName: accountName,
		NickName:    nickName,
		PortraitURL: portraitURL,
		Birthday:    birthday,
		Phone:       phone,
		Gender:      gender,
	}
	err := service.CreateUserInfo(userInfo)
	if err != nil {
		log.PrintLog("CreateUserInfoHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.CreateUserInfoFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.CreateUserInfoSuccess)
	}
}

func UpdateUserInfoHandler(context *gin.Context) {
	// 暂不支持
}

func UpdateNickNameHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	nickName := context.PostForm(nickNameParam)

	err := service.UpdateNickName(accountName, nickName)
	if err != nil {
		log.PrintLog("UpdateNickNameHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateNickNameFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateNickNameSuccess)
	}
}

func UpdatePortraitURLHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	portraitURL := context.PostForm(portraitURLParam)

	err := service.UpdatePortraitURL(accountName, portraitURL)
	if err != nil {
		log.PrintLog("UpdatePortraitURLHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdatePortraitURLFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdatePortraitURLSuccess)
	}
}

func UpdateBirthdayHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	birthday := context.PostForm(birthdayParam)

	err := service.UpdateBirthday(accountName, birthday)
	if err != nil {
		log.PrintLog("UpdateBirthdayHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateBirthdayFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateBirthdaySuccess)
	}
}

func UpdatePhoneHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	phone := context.PostForm(phoneParam)

	err := service.UpdatePhone(accountName, phone)
	if err != nil {
		log.PrintLog("UpdatePhoneHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdatePhoneFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdatePhoneSuccess)
	}
}

func UpdateGenderHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)
	gender := context.PostForm(genderParam)

	err := service.UpdateGender(accountName, gender)
	if err != nil {
		log.PrintLog("UpdateGenderHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateGenderFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateGenderSuccess)
	}
}

func DeleteUserInfoHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)

	err := service.DeleteUserInfo(accountName)
	if err != nil {
		log.PrintLog("UpdateBirthdayHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateGenderFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateGenderSuccess)
	}
}

func QueryUserInfoURLHandler(context *gin.Context) {
	accountName := context.PostForm(accountNameParam)

	userInfo, err := service.QueryUserInfo(accountName)
	if err != nil {
		log.PrintLog("QueryUserInfoURLHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.QueryUserInfoFail)
	} else {
		context.JSON(http.StatusOK, gin.H{
			"code": myErrors.QueryUserInfoSuccess.Code(),
			"info": userInfo,
		})
	}
}
