package dao

import (
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/model"
	"github.com/eyebluecn/tank/code/tool/builder"
	"github.com/eyebluecn/tank/code/tool/result"
	"github.com/eyebluecn/tank/code/tool/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type BridgeDao struct {
	BaseDao
}

//find by uuid. if not found return nil.
func (this *BridgeDao) FindByUuid(uuid string) *model.Bridge {

	var bridge = &model.Bridge{}
	db := core.CONTEXT.GetDB().Where("uuid = ?", uuid).First(bridge)
	if db.Error != nil {
		if db.Error.Error() == result.DB_ERROR_NOT_FOUND {
			return nil
		} else {
			panic(db.Error)
		}
	}
	return bridge

}

//find by uuid. if not found panic NotFound error
func (this *BridgeDao) CheckByUuid(uuid string) *model.Bridge {

	entity := this.FindByUuid(uuid)
	if entity == nil {
		panic(result.NotFound("not found record with uuid = %s", uuid))
	}

	return entity

}

//find by shareUuid and matterUuid. if not found panic NotFound error.
func (this *BridgeDao) CheckByShareUuidAndMatterUuid(shareUuid string, matterUuid string) *model.Bridge {

	var bridge = &model.Bridge{}
	db := core.CONTEXT.GetDB().Where("share_uuid = ? AND matter_uuid = ?", shareUuid, matterUuid).First(bridge)
	if db.Error != nil {
		if db.Error.Error() == result.DB_ERROR_NOT_FOUND {
			panic(result.NotFound("not found record with shareUuid = %s and matterUuid = %s", shareUuid, matterUuid))
		} else {
			panic(db.Error)
		}
	}

	return bridge
}

//get pager
func (this *BridgeDao) PlainPage(page int, pageSize int, shareUuid string, sortArray []builder.OrderPair) (int, []*model.Bridge) {

	var wp = &builder.WherePair{}

	if shareUuid != "" {
		wp = wp.And(&builder.WherePair{Query: "share_uuid = ?", Args: []interface{}{shareUuid}})
	}

	var conditionDB *gorm.DB
	conditionDB = core.CONTEXT.GetDB().Model(&model.Bridge{}).Where(wp.Query, wp.Args...)

	count := 0
	db := conditionDB.Count(&count)
	this.PanicError(db.Error)

	var bridges []*model.Bridge
	db = conditionDB.Order(this.GetSortString(sortArray)).Offset(page * pageSize).Limit(pageSize).Find(&bridges)
	this.PanicError(db.Error)

	return count, bridges
}

//get pager
func (this *BridgeDao) Page(page int, pageSize int, shareUuid string, sortArray []builder.OrderPair) *model.Pager {

	count, bridges := this.PlainPage(page, pageSize, shareUuid, sortArray)
	pager := model.NewPager(page, pageSize, count, bridges)

	return pager
}

func (this *BridgeDao) Create(bridge *model.Bridge) *model.Bridge {

	timeUUID, _ := uuid.NewV4()
	bridge.Uuid = string(timeUUID.String())
	bridge.CreateTime = time.Now()
	bridge.UpdateTime = time.Now()
	bridge.Sort = time.Now().UnixNano() / 1e6
	db := core.CONTEXT.GetDB().Create(bridge)
	this.PanicError(db.Error)

	return bridge
}

func (this *BridgeDao) Save(bridge *model.Bridge) *model.Bridge {

	bridge.UpdateTime = time.Now()
	db := core.CONTEXT.GetDB().Save(bridge)
	this.PanicError(db.Error)

	return bridge
}

func (this *BridgeDao) Delete(bridge *model.Bridge) {

	db := core.CONTEXT.GetDB().Delete(&bridge)
	this.PanicError(db.Error)
}

func (this *BridgeDao) DeleteByMatterUuid(matterUuid string) {

	var wp = &builder.WherePair{}

	wp = wp.And(&builder.WherePair{Query: "matter_uuid = ?", Args: []interface{}{matterUuid}})

	db := core.CONTEXT.GetDB().Where(wp.Query, wp.Args).Delete(model.Bridge{})
	this.PanicError(db.Error)
}

func (this *BridgeDao) DeleteByShareUuid(shareUuid string) {

	var wp = &builder.WherePair{}

	wp = wp.And(&builder.WherePair{Query: "share_uuid = ?", Args: []interface{}{shareUuid}})

	db := core.CONTEXT.GetDB().Where(wp.Query, wp.Args).Delete(model.Bridge{})
	this.PanicError(db.Error)
}

func (this *BridgeDao) FindByShareUuid(shareUuid string) []*model.Bridge {

	if shareUuid == "" {
		panic(result.BadRequest("shareUuid cannot be nil"))
	}

	var bridges []*model.Bridge

	db := core.CONTEXT.GetDB().
		Where("share_uuid = ?", shareUuid).
		Find(&bridges)
	this.PanicError(db.Error)

	return bridges
}

func (this *BridgeDao) Cleanup() {
	this.Logger.Info("[BridgeDao] cleanup: delete all bridge records.")
	db := core.CONTEXT.GetDB().Where("uuid is not null").Delete(model.Bridge{})
	this.PanicError(db.Error)
}
