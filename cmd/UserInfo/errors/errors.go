package errors

import (
	errorsUtil "mall/pkg/util/errors"
)

var (
	CreateUserInfoSuccess    = errorsUtil.New(2001, "create user info success")
	CreateUserInfoFail       = errorsUtil.New(2002, "create user info fail")
	UpdateNickNameSuccess    = errorsUtil.New(2003, "UpdateNickName success")
	UpdateNickNameFail       = errorsUtil.New(2004, "UpdateNickName fail")
	UpdatePortraitURLSuccess = errorsUtil.New(2005, "UpdatePortraitURL success")
	UpdatePortraitURLFail    = errorsUtil.New(2006, "UpdatePortraitURL fail")
	UpdateBirthdaySuccess    = errorsUtil.New(2007, "UpdateBirthday success")
	UpdateBirthdayFail       = errorsUtil.New(2008, "UpdateBirthday fail")
	UpdatePhoneSuccess       = errorsUtil.New(2009, "UpdatePhone success")
	UpdatePhoneFail          = errorsUtil.New(2010, "UpdatePhone fail")
	UpdateGenderSuccess      = errorsUtil.New(2011, "UpdateGender success")
	UpdateGenderFail         = errorsUtil.New(2012, "UpdateGender fail")
	QueryUserInfoSuccess     = errorsUtil.New(2013, "QueryUserInfoFail success")
	QueryUserInfoFail        = errorsUtil.New(2014, "QueryUserInfoFail fail")
)
