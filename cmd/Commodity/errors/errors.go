package errors

import (
	errorsUtil "mall/pkg/util/errors"
)

var (
	AddCommoditySuccess       = errorsUtil.New(3001, "AddCommodity success")
	AddCommodityFail          = errorsUtil.New(3002, "AddCommodity fail")
	ParameterError            = errorsUtil.New(3003, "parameter error")
	DeleteCommoditySuccess    = errorsUtil.New(3004, "DeleteCommodity success")
	DeleteCommodityFail       = errorsUtil.New(3005, "DeleteCommodity fail")
	UpdateMainCategorySuccess = errorsUtil.New(3006, "UpdateMainCategory success")
	UpdateMainCategoryFail    = errorsUtil.New(3007, "DeleteCommodity fail")
	UpdateSubCategorySuccess  = errorsUtil.New(3008, "UpdateSubCategory success")
	UpdateSubCategoryFail     = errorsUtil.New(3009, "UpdateSubCategory fail")
	UpdateNameSuccess         = errorsUtil.New(3010, "UpdateName success")
	UpdateNameFail            = errorsUtil.New(3011, "UpdateName fail")
	UpdatePriceSuccess        = errorsUtil.New(3012, "UpdatePrice success")
	UpdatePriceFail           = errorsUtil.New(3013, "UpdatePrice fail")
	QueryCommoditySuccess     = errorsUtil.New(3014, "QueryCommodity success")
	QueryCommodityFail        = errorsUtil.New(3015, "QueryCommodity fail")
	QueryCommoditiesSuccess   = errorsUtil.New(3016, "QueryCommodities success")
	QueryCommoditiesFail      = errorsUtil.New(3017, "QueryCommodities fail")
)
