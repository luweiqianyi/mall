package errors

import (
	errorsUtil "mall/pkg/util/errors"
)

var (
	RegisterSuccess       = errorsUtil.New(1001, "register success")
	RegisterFail          = errorsUtil.New(1002, "register failed, account already exist")
	UnRegisterSuccess     = errorsUtil.New(1003, "unregister success")
	UnRegisterFail        = errorsUtil.New(1004, "unregister fail")
	PasswordChangeSuccess = errorsUtil.New(1005, "password change success")
	PasswordChangeFail    = errorsUtil.New(1006, "password change fail")
	LoginSuccess          = errorsUtil.New(1007, "login success")
	LoginFail             = errorsUtil.New(1008, "login fail")
	LogoutSuccess         = errorsUtil.New(1009, "Logout success")
	LogoutFail            = errorsUtil.New(1010, "Logout fail")
)
