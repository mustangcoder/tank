package rest

import (
	"github.com/eyebluecn/tank/code/constant"
	"github.com/eyebluecn/tank/code/tool/result"
	"net/http"
)

type ServerStatusController struct {
	BaseController
}

func (this *ServerStatusController) Init() {
	this.BaseController.Init()
}

func (this *ServerStatusController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {

	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))
	routeMap["/api/server/healthCheck"] = this.Wrap(this.HealthCheck, constant.USER_ROLE_GUEST)

	return routeMap
}

func (this *ServerStatusController) HealthCheck(writer http.ResponseWriter, request *http.Request) *result.WebResult {
	this.Logger.Info("HealthCheck is ok!")
	return this.Success("OK")
}
