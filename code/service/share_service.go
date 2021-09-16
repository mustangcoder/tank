package service

import (
	"github.com/eyebluecn/tank/code/constant"
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/dao"
	"github.com/eyebluecn/tank/code/model"
	"github.com/eyebluecn/tank/code/tool/builder"
	"github.com/eyebluecn/tank/code/tool/i18n"
	"github.com/eyebluecn/tank/code/tool/result"
	"math"
	"net/http"
	"strings"
	"time"
)

//@Service
type ShareService struct {
	core.BaseBean
	shareDao  *dao.ShareDao
	matterDao *dao.MatterDao
	bridgeDao *dao.BridgeDao
	userDao   *dao.UserDao
}

func (this *ShareService) Init() {
	this.BaseBean.Init()

	b := core.CONTEXT.GetBean(this.shareDao)
	if b, ok := b.(*dao.ShareDao); ok {
		this.shareDao = b
	}

	b = core.CONTEXT.GetBean(this.matterDao)
	if b, ok := b.(*dao.MatterDao); ok {
		this.matterDao = b
	}

	b = core.CONTEXT.GetBean(this.bridgeDao)
	if b, ok := b.(*dao.BridgeDao); ok {
		this.bridgeDao = b
	}

	b = core.CONTEXT.GetBean(this.userDao)
	if b, ok := b.(*dao.UserDao); ok {
		this.userDao = b
	}

}

func (this *ShareService) Detail(uuid string) *model.Share {

	share := this.shareDao.CheckByUuid(uuid)

	return share
}

// check whether shareUuid and shareCode matches. check whether user can access.
func (this *ShareService) CheckShare(request *http.Request, shareUuid string, code string, user *model.User) *model.Share {

	share := this.shareDao.CheckByUuid(shareUuid)
	//if self, not need shareCode
	if user == nil || user.Uuid != share.UserUuid {
		//if not login or not self's share, shareCode is required.
		if code == "" {
			panic(result.CustomWebResultI18n(request, result.NEED_SHARE_CODE, i18n.ShareCodeRequired))
		} else if share.Code != code {
			panic(result.CustomWebResultI18n(request, result.SHARE_CODE_ERROR, i18n.ShareCodeError))
		} else {
			if !share.ExpireInfinity {
				if share.ExpireTime.Before(time.Now()) {
					panic(result.BadRequest("share expired"))
				}
			}
		}
	}
	return share
}

//check whether a user can access a matter. shareRootUuid is matter's parent(or parent's parent and so on)
func (this *ShareService) ValidateMatter(request *http.Request, shareUuid string, code string, user *model.User, shareRootUuid string, matter *model.Matter) {

	if matter == nil {
		panic(result.Unauthorized("matter cannot be nil"))
	}

	//if self's matter, ok.
	if user != nil && matter.UserUuid == user.Uuid {
		return
	}

	if shareUuid == "" || code == "" || shareRootUuid == "" {
		panic(result.Unauthorized("shareUuid,code,shareRootUuid cannot be null"))
	}

	share := this.CheckShare(request, shareUuid, code, user)

	//if shareRootUuid is root. Bridge must has record.
	if shareRootUuid == constant.MATTER_ROOT {

		this.bridgeDao.CheckByShareUuidAndMatterUuid(share.Uuid, matter.Uuid)

	} else {
		//check whether shareRootMatter is being sharing
		shareRootMatter := this.matterDao.CheckByUuid(shareRootUuid)
		this.bridgeDao.CheckByShareUuidAndMatterUuid(share.Uuid, shareRootMatter.Uuid)

		// shareRootMatter is ancestor of matter.
		child := strings.HasPrefix(matter.Path, shareRootMatter.Path)
		if !child {
			panic(result.BadRequest("%s is not %s's children", matter.Uuid, shareRootUuid))
		}
	}

}

//delete user's shares and corresponding bridges.
func (this *ShareService) DeleteSharesByUser(request *http.Request, currentUser *model.User) {

	//delete share and bridges.
	pageSize := 100
	var sortArray []builder.OrderPair
	count, _ := this.shareDao.PlainPage(0, pageSize, currentUser.Uuid, sortArray)
	if count > 0 {
		var totalPages = int(math.Ceil(float64(count) / float64(pageSize)))

		var page int
		for page = 0; page < totalPages; page++ {
			_, shares := this.shareDao.PlainPage(0, pageSize, currentUser.Uuid, sortArray)
			for _, share := range shares {
				this.bridgeDao.DeleteByShareUuid(share.Uuid)

				//delete this share
				this.shareDao.Delete(share)
			}
		}

	}

}
