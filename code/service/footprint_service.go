package service

import (
	"encoding/json"
	"github.com/eyebluecn/tank/code/constant"
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/dao"
	"github.com/eyebluecn/tank/code/model"
	"github.com/eyebluecn/tank/code/tool/util"
	"net/http"
	"time"
)

//@Service
type FootprintService struct {
	core.BaseBean
	footprintDao *dao.FootprintDao
	userDao      *dao.UserDao
	userService  *UserService
}

func (this *FootprintService) Init() {
	this.BaseBean.Init()

	b := core.CONTEXT.GetBean(this.footprintDao)
	if b, ok := b.(*dao.FootprintDao); ok {
		this.footprintDao = b
	}

	b = core.CONTEXT.GetBean(this.userDao)
	if b, ok := b.(*dao.UserDao); ok {
		this.userDao = b
	}

	b = core.CONTEXT.GetBean(this.userService)
	if b, ok := b.(*UserService); ok {
		this.userService = b
	}

}

func (this *FootprintService) Detail(uuid string) *model.Footprint {

	footprint := this.footprintDao.CheckByUuid(uuid)

	return footprint
}

//log a request.
func (this *FootprintService) Trace(request *http.Request, duration time.Duration, success bool) {

	params := make(map[string][]string)

	//POST params
	values := request.PostForm
	for key, val := range values {
		params[key] = val
	}
	//GET params
	values1 := request.URL.Query()
	for key, val := range values1 {
		params[key] = val
	}

	//ignore password.
	for key, _ := range params {
		if key == constant.PASSWORD_KEY || key == "password" || key == "adminPassword" {
			params[key] = []string{"******"}
		}
	}

	paramsString := "{}"
	paramsData, err := json.Marshal(params)
	if err == nil {
		paramsString = string(paramsData)
	}

	footprint := &model.Footprint{
		Ip:      util.GetIpAddress(request),
		Host:    request.Host,
		Uri:     request.URL.Path,
		Params:  paramsString,
		Cost:    int64(duration / time.Millisecond),
		Success: success,
	}

	//if db not config just print content.
	if core.CONFIG.Installed() {
		user := this.userService.FindUser(request)
		userUuid := ""
		if user != nil {
			userUuid = user.Uuid
		}
		footprint.UserUuid = userUuid
		footprint = this.footprintDao.Create(footprint)
	}

	this.Logger.Info("Ip:%s Cost:%d Uri:%s Params:%s", footprint.Ip, int64(duration/time.Millisecond), footprint.Uri, paramsString)

}

func (this *FootprintService) Bootstrap() {

	this.Logger.Info("Immediately delete Footprint data of 8 days ago.")

	go core.RunWithRecovery(this.CleanOldData)
}

func (this *FootprintService) CleanOldData() {

	day8Ago := time.Now()
	day8Ago = day8Ago.AddDate(0, 0, -8)
	day8Ago = util.FirstSecondOfDay(day8Ago)

	this.Logger.Info("Delete footprint data before %s", util.ConvertTimeToDateTimeString(day8Ago))

	this.footprintDao.DeleteByCreateTimeBefore(day8Ago)
}
