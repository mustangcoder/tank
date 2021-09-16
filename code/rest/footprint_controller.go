package rest

import (
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/dao"
	"github.com/eyebluecn/tank/code/service"
	"net/http"
)

type FootprintController struct {
	BaseController
	footprintDao     *dao.FootprintDao
	footprintService *service.FootprintService
}

func (this *FootprintController) Init() {
	this.BaseController.Init()

	b := core.CONTEXT.GetBean(this.footprintDao)
	if b, ok := b.(*dao.FootprintDao); ok {
		this.footprintDao = b
	}

	b = core.CONTEXT.GetBean(this.footprintService)
	if b, ok := b.(*service.FootprintService); ok {
		this.footprintService = b
	}

}

func (this *FootprintController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {

	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))

	return routeMap
}
