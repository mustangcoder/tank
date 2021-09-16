package rest

import (
	"fmt"
	"github.com/eyebluecn/tank/code/constant"
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/dao"
	"github.com/eyebluecn/tank/code/service"
	"github.com/eyebluecn/tank/code/tool/i18n"
	"github.com/eyebluecn/tank/code/tool/result"
	"github.com/eyebluecn/tank/code/tool/util"
	"github.com/json-iterator/go"
	"go/types"
	"net/http"
)

type BaseController struct {
	core.BaseBean
	UserService *service.UserService
	UserDao     *dao.UserDao
	sessionDao  *dao.SessionDao
}

func (this *BaseController) Init() {

	this.BaseBean.Init()

	b := core.CONTEXT.GetBean(this.UserService)
	if b, ok := b.(*service.UserService); ok {
		this.UserService = b
	}

	b = core.CONTEXT.GetBean(this.UserDao)
	if b, ok := b.(*dao.UserDao); ok {
		this.UserDao = b
	}

	b = core.CONTEXT.GetBean(this.sessionDao)
	if b, ok := b.(*dao.SessionDao); ok {
		this.sessionDao = b
	}

}

func (this *BaseController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {

	return make(map[string]func(writer http.ResponseWriter, request *http.Request))
}

//handle some special routes, eg. params in the url.
func (this *BaseController) HandleRoutes(writer http.ResponseWriter, request *http.Request) (func(writer http.ResponseWriter, request *http.Request), bool) {
	return nil, false
}

//wrap the handle method.
func (this *BaseController) Wrap(f func(writer http.ResponseWriter, request *http.Request) *result.WebResult, qualifiedRole string) func(w http.ResponseWriter, r *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		var webResult *result.WebResult = nil

		//if the api not annotated with GUEST. login is required.
		if qualifiedRole != constant.USER_ROLE_GUEST {
			user := this.UserService.CheckUser(request)

			if user.Status == constant.USER_STATUS_DISABLED {
				//check user's status
				webResult = result.CustomWebResultI18n(request, result.USER_DISABLED, i18n.UserDisabled)
			} else {
				if qualifiedRole == constant.USER_ROLE_ADMINISTRATOR && user.Role != constant.USER_ROLE_ADMINISTRATOR {
					webResult = result.ConstWebResult(result.UNAUTHORIZED)
				} else {
					webResult = f(writer, request)
				}
			}

		} else {
			webResult = f(writer, request)
		}

		//if webResult not nil. response a json. if webResult is nil, return empty body or binary content.
		if webResult != nil {

			writer.Header().Set("Content-Type", "application/json;charset=UTF-8")

			b, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(webResult)

			this.PanicError(err)

			writer.WriteHeader(result.FetchHttpStatus(webResult.Code))

			_, err = fmt.Fprintf(writer, string(b))
			this.PanicError(err)
		}

	}
}

//response a success result. 1.string 2. WebResult 3.nil pointer 4.any type
func (this *BaseController) Success(data interface{}) *result.WebResult {
	var webResult *result.WebResult = nil
	if value, ok := data.(string); ok {
		//a simple message
		webResult = &result.WebResult{Code: result.OK.Code, Msg: value}
	} else if value, ok := data.(*result.WebResult); ok {
		//a webResult
		webResult = value
	} else if _, ok := data.(types.Nil); ok {
		//nil pointer means OK.
		webResult = result.ConstWebResult(result.OK)
	} else {
		//other type.
		webResult = &result.WebResult{Code: result.OK.Code, Data: data}
	}
	return webResult
}

//allow cors.
func (this *BaseController) AllowCORS(writer http.ResponseWriter) {
	util.AllowCORS(writer)
}
