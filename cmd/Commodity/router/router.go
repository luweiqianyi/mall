package router

import (
	"github.com/gin-gonic/gin"
	"mall/cmd/Commodity/config"
	"mall/cmd/Commodity/entity"
	myErrors "mall/cmd/Commodity/errors"
	"mall/cmd/Commodity/service"
	"mall/pkg/log"
	"mall/pkg/util/errors"
	"net/http"
	"strconv"
	"sync"
)

// 前端请求参数
const (
	identifyIDParam     = "identifyID"
	mainCategoryParam   = "mainCategory"
	subCategoryParam    = "subCategory"
	commodityNameParam  = "commodityName"
	commodityPriceParam = "commodityPrice"
	commodityPDParam    = "commodityPD"
	commodityExpParam   = "commodityExp"
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

	r.AddRoute(AddCommodityURL)
	r.AddRoute(DeleteCommodityURL)
	r.AddRoute(QueryCommodityURL)
	r.AddRoute(QueryCommoditiesURL)
	r.AddRoute(UpdateMainCategoryURL)
	r.AddRoute(UpdateSubCategoryURL)
	r.AddRoute(UpdateNameURL)
	r.AddRoute(UpdatePriceURL)

	eg.Run(config.GetGConfig().GetServiceHttpPort())
}

func (r *WebRouter) AddRoute(URLPath string) {
	switch URLPath {
	case AddCommodityURL:
		r.eg.POST(URLPath, AddCommodityHandler)
	case DeleteCommodityURL:
		r.eg.POST(URLPath, DeleteCommodityURLHandler)
	case QueryCommodityURL:
		r.eg.POST(URLPath, QueryCommodityURLHandler)
	case QueryCommoditiesURL:
		r.eg.GET(URLPath, QueryCommoditiesURLHandler)
	case UpdateMainCategoryURL:
		r.eg.POST(URLPath, UpdateMainCategoryURLHandler)
	case UpdateSubCategoryURL:
		r.eg.POST(URLPath, UpdateSubCategoryURLHandler)
	case UpdateNameURL:
		r.eg.POST(URLPath, UpdateNameURLHandler)
	case UpdatePriceURL:
		r.eg.POST(URLPath, UpdatePriceURLHandler)
	}
}

func BuildHttpResponse(context *gin.Context, httpStatusCode int, err errors.Code) {
	context.JSON(httpStatusCode, gin.H{
		"code": err.Code(),
		"msg":  err.Message(),
	})
}

func getParameter(param string, context *gin.Context) interface{} {
	switch param {
	case identifyIDParam:
		return context.PostForm(identifyIDParam)
	case mainCategoryParam:
		return context.PostForm(mainCategoryParam)
	case subCategoryParam:
		return context.PostForm(subCategoryParam)
	case commodityNameParam:
		return context.PostForm(commodityNameParam)
	case commodityPriceParam:
		return context.PostForm(commodityPriceParam)
	case commodityPDParam:
		return context.PostForm(commodityPDParam)
	case commodityExpParam:
		return context.PostForm(commodityExpParam)
	}
	return ""
}

// AddCommodityHandler 商品录入
func AddCommodityHandler(context *gin.Context) {
	price, err := strconv.ParseFloat(getParameter(commodityPriceParam, context).(string), 64)
	if err != nil {
		BuildHttpResponse(context, http.StatusOK, myErrors.ParameterError)
		return
	}
	info := entity.TbCommodity{
		IdentifyID:   getParameter(identifyIDParam, context).(string),
		MainCategory: getParameter(mainCategoryParam, context).(string),
		SubCategory:  getParameter(subCategoryParam, context).(string),
		CommodityProperty: entity.CommodityProperty{
			Name:  getParameter(commodityNameParam, context).(string),
			Price: price,
			PD:    getParameter(commodityPDParam, context).(string),
			Exp:   getParameter(commodityExpParam, context).(string),
		},
	}
	err = service.AddCommodity(info)
	if err != nil {
		log.PrintLog("AddCommodityHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.AddCommodityFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.AddCommoditySuccess)
	}
}

func DeleteCommodityURLHandler(context *gin.Context) {
	identifyID := getParameter(identifyIDParam, context).(string)
	err := service.DeleteCommodity(identifyID)
	if err != nil {
		log.PrintLog("DeleteCommodityURLHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.DeleteCommodityFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.DeleteCommoditySuccess)
	}
}

func QueryCommodityURLHandler(context *gin.Context) {
	identifyID := getParameter(identifyIDParam, context).(string)
	info, err := service.QueryCommodity(identifyID)
	if err != nil {
		log.PrintLog("QueryCommodityURLHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.QueryCommodityFail)
	} else {
		context.JSON(http.StatusOK, gin.H{
			"code": myErrors.QueryCommoditySuccess.Code(),
			"info": info,
		})
	}
}

func QueryCommoditiesURLHandler(context *gin.Context) {
	commodities, err := service.QueryCommodities()
	if err != nil {
		log.PrintLog("QueryCommodityURLHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.QueryCommoditiesFail)
	} else {
		context.JSON(http.StatusOK, gin.H{
			"code": myErrors.QueryCommoditiesSuccess.Code(),
			"info": commodities,
		})
	}
}

func UpdateMainCategoryURLHandler(context *gin.Context) {
	identifyID := getParameter(identifyIDParam, context).(string)
	mainCategory := getParameter(mainCategoryParam, context).(string)
	err := service.UpdateCommodityMainCategory(identifyID, mainCategory)
	if err != nil {
		log.PrintLog("UpdateMainCategoryURLHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateMainCategoryFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateMainCategorySuccess)
	}
}

func UpdateSubCategoryURLHandler(context *gin.Context) {
	identifyID := getParameter(identifyIDParam, context).(string)
	subCategory := getParameter(subCategoryParam, context).(string)
	err := service.UpdateCommoditySubCategory(identifyID, subCategory)
	if err != nil {
		log.PrintLog("UpdateSubCategoryURLHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateSubCategoryFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateSubCategorySuccess)
	}
}

func UpdateNameURLHandler(context *gin.Context) {
	identifyID := getParameter(identifyIDParam, context).(string)
	name := getParameter(commodityNameParam, context).(string)
	err := service.UpdateCommodityName(identifyID, name)
	if err != nil {
		log.PrintLog("UpdateNameURLHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateNameFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdateNameSuccess)
	}
}

func UpdatePriceURLHandler(context *gin.Context) {
	identifyID := getParameter(identifyIDParam, context).(string)
	price := getParameter(commodityPriceParam, context).(string)
	commodityPrice, _ := strconv.ParseFloat(price, 64)
	err := service.UpdateCommodityPrice(identifyID, commodityPrice)
	if err != nil {
		log.PrintLog("UpdatePriceURLHandler failed,err=%v", err)
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdatePriceFail)
	} else {
		BuildHttpResponse(context, http.StatusOK, myErrors.UpdatePriceSuccess)
	}
}
