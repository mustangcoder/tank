package service

import (
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/dao"
	"github.com/eyebluecn/tank/code/model"
)

//@Service
type BridgeService struct {
	core.BaseBean
	bridgeDao *dao.BridgeDao
	userDao   *dao.UserDao
}

func (this *BridgeService) Init() {
	this.BaseBean.Init()

	b := core.CONTEXT.GetBean(this.bridgeDao)
	if b, ok := b.(*dao.BridgeDao); ok {
		this.bridgeDao = b
	}

	b = core.CONTEXT.GetBean(this.userDao)
	if b, ok := b.(*dao.UserDao); ok {
		this.userDao = b
	}

}

func (this *BridgeService) Detail(uuid string) *model.Bridge {

	bridge := this.bridgeDao.CheckByUuid(uuid)

	return bridge
}
