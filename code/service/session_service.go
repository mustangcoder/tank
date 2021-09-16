package service

import (
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/dao"
)

//@Service
type SessionService struct {
	core.BaseBean
	userDao    *dao.UserDao
	sessionDao *dao.SessionDao
}

func (this *SessionService) Init() {
	this.BaseBean.Init()

	b := core.CONTEXT.GetBean(this.userDao)
	if b, ok := b.(*dao.UserDao); ok {
		this.userDao = b
	}

	b = core.CONTEXT.GetBean(this.sessionDao)
	if b, ok := b.(*dao.SessionDao); ok {
		this.sessionDao = b
	}

}

//System cleanup.
func (this *SessionService) Cleanup() {

	this.Logger.Info("[SessionService] clean up. Delete all Session. total:%d", core.CONTEXT.GetSessionCache().Count())

	core.CONTEXT.GetSessionCache().Truncate()
}
