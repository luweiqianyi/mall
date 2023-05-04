package service

import (
	"context"
	"mall/cmd/Account/config"
	"mall/cmd/Account/dao"
	authPb "mall/cmd/Auth/pb"
	authService "mall/cmd/Auth/service"
	"mall/pkg/log"
	"mall/pkg/util"
	"time"
)

func Register(accountName string, password string) bool {
	return dao.Register(accountName, password)
}

func UnRegister(accountName string) bool {
	success := dao.UnRegister(accountName)
	if !success {
		log.PrintLog("unregister failed!\n")
		return false
	}

	conn := config.GetGConfig().GetRpcClientConn()
	if conn == nil {
		log.PrintLog("UnRegister grpc connection invalid")
		return false
	}

	// 向远程AuthService发起RPC请求，删除用户token
	rpcClient := authPb.NewAuthServiceClient(conn)
	resp, _ := rpcClient.DelToken(context.Background(), &authPb.DelTokenRequest{
		AccountName: accountName,
	})
	if resp.Code == authService.DelTokenFail {
		log.PrintLog("login fail, save account[%s]'s token to remote failed, reason=%v", accountName, resp)
		return false
	}
	return true
}

func ChangePassword(accountName string, password string) bool {
	return dao.ChangePassword(accountName, password)
}

// LoginUsingRpc 登录，将登录token保存在Auth服务上
// TODO 优化：用户的账号和密码可能发生泄露，检测用户登录是否是常用地址登录，如果不是，向用户发送短信验证码进行验证，验证通过，才算登录成功
func LoginUsingRpc(accountName string, password string) (string, bool) {
	success := dao.Login(accountName, password)
	if !success {
		log.PrintLog("user not exist,please register first!\n")
		return "", false
	}

	conn := config.GetGConfig().GetRpcClientConn()
	if conn == nil {
		log.PrintLog("UnRegister grpc connection invalid")
		return "UnRegister grpc connection invalid", false
	}

	// 向远程AuthService发起RPC请求，获取用户token
	rpcClient := authPb.NewAuthServiceClient(conn)
	resp, _ := rpcClient.QueryToken(context.Background(), &authPb.QueryTokenRequest{
		AccountName: accountName,
	})
	var token string
	if resp.Code == authService.QueryTokenCodeFail {
		log.PrintLog("login, query account[%s]'s token failed, reason=%v", accountName, resp)
		token = util.GenerateToken()
	} else {
		token = resp.Token
	}

	// 向远程AuthService发起RPC请求：新增或者更新Token
	resp2, _ := rpcClient.SaveToken(context.Background(), &authPb.SaveTokenRequest{
		AccountName:     accountName,
		Token:           token,
		TokenExpireTime: int64(time.Hour * 24),
	})
	if resp2.Code == authService.SaveTokenCodeFail {
		log.PrintLog("login fail, save account[%s]'s token to remote failed, reason=%v", accountName, resp2)
		return "", false
	}
	return token, true
}

// LogoutUsingRpc 退出登录，调用Auth服务的方法将Token删除
func LogoutUsingRpc(accountName string) bool {
	conn := config.GetGConfig().GetRpcClientConn()
	if conn == nil {
		log.PrintLog("UnRegister grpc connection invalid")
		return false
	}

	// 向远程AuthService发起RPC请求，删除用户Token
	rpcClient := authPb.NewAuthServiceClient(conn)
	resp, _ := rpcClient.DelToken(context.Background(), &authPb.DelTokenRequest{
		AccountName: accountName,
	})
	if resp.Code == authService.DelTokenFail {
		log.PrintLog("logout fail, delete account[%s]'s remote token failed, reason=%v", accountName, resp)
		return false
	}
	return true
}
