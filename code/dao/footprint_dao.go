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

type FootprintDao struct {
	BaseDao
}

//find by uuid. if not found return nil.
func (this *FootprintDao) FindByUuid(uuid string) *model.Footprint {
	var entity = &model.Footprint{}
	db := core.CONTEXT.GetDB().Where("uuid = ?", uuid).First(entity)
	if db.Error != nil {
		if db.Error.Error() == result.DB_ERROR_NOT_FOUND {
			return nil
		} else {
			panic(db.Error)
		}
	}
	return entity
}

//find by uuid. if not found panic NotFound error
func (this *FootprintDao) CheckByUuid(uuid string) *model.Footprint {
	entity := this.FindByUuid(uuid)
	if entity == nil {
		panic(result.NotFound("not found record with uuid = %s", uuid))
	}
	return entity
}

func (this *FootprintDao) Page(page int, pageSize int, userUuid string, sortArray []builder.OrderPair) *model.Pager {

	var wp = &builder.WherePair{}

	if userUuid != "" {
		wp = wp.And(&builder.WherePair{Query: "user_uuid = ?", Args: []interface{}{userUuid}})
	}

	var conditionDB *gorm.DB
	conditionDB = core.CONTEXT.GetDB().Model(&model.Footprint{}).Where(wp.Query, wp.Args...)

	count := 0
	db := conditionDB.Count(&count)
	this.PanicError(db.Error)

	var footprints []*model.Footprint
	db = conditionDB.Order(this.GetSortString(sortArray)).Offset(page * pageSize).Limit(pageSize).Find(&footprints)
	this.PanicError(db.Error)
	pager := model.NewPager(page, pageSize, count, footprints)

	return pager
}

func (this *FootprintDao) Create(footprint *model.Footprint) *model.Footprint {

	timeUUID, _ := uuid.NewV4()
	footprint.Uuid = string(timeUUID.String())
	footprint.CreateTime = time.Now()
	footprint.UpdateTime = time.Now()
	footprint.Sort = time.Now().UnixNano() / 1e6
	db := core.CONTEXT.GetDB().Create(footprint)
	this.PanicError(db.Error)

	return footprint
}

func (this *FootprintDao) Save(footprint *model.Footprint) *model.Footprint {

	footprint.UpdateTime = time.Now()
	db := core.CONTEXT.GetDB().Save(footprint)
	this.PanicError(db.Error)

	return footprint
}

func (this *FootprintDao) Delete(footprint *model.Footprint) {

	db := core.CONTEXT.GetDB().Delete(&footprint)
	this.PanicError(db.Error)
}

func (this *FootprintDao) CountBetweenTime(startTime time.Time, endTime time.Time) int64 {
	var count int64
	db := core.CONTEXT.GetDB().Model(&model.Footprint{}).Where("create_time >= ? AND create_time <= ?", startTime, endTime).Count(&count)
	this.PanicError(db.Error)
	return count
}

func (this *FootprintDao) UvBetweenTime(startTime time.Time, endTime time.Time) int64 {

	var wp = &builder.WherePair{Query: "create_time >= ? AND create_time <= ?", Args: []interface{}{startTime, endTime}}

	var count int64
	db := core.CONTEXT.GetDB().Model(&model.Footprint{}).Where(wp.Query, wp.Args...).Count(&count)
	if count == 0 {
		return 0
	}

	db = core.CONTEXT.GetDB().Model(&model.Footprint{}).Where("create_time >= ? AND create_time <= ?", startTime, endTime).Select("COUNT(DISTINCT(ip))")
	this.PanicError(db.Error)
	row := db.Row()
	err := row.Scan(&count)
	this.PanicError(err)
	return count
}

func (this *FootprintDao) AvgCostBetweenTime(startTime time.Time, endTime time.Time) int64 {

	var wp = &builder.WherePair{Query: "create_time >= ? AND create_time <= ?", Args: []interface{}{startTime, endTime}}

	var count int64
	db := core.CONTEXT.GetDB().Model(&model.Footprint{}).Where(wp.Query, wp.Args...).Count(&count)
	if count == 0 {
		return 0
	}

	var cost float64
	db = core.CONTEXT.GetDB().Model(&model.Footprint{}).Where(wp.Query, wp.Args...).Select("AVG(cost)")
	this.PanicError(db.Error)
	row := db.Row()
	err := row.Scan(&cost)
	this.PanicError(err)
	return int64(cost)
}

func (this *FootprintDao) DeleteByCreateTimeBefore(createTime time.Time) {
	db := core.CONTEXT.GetDB().Where("create_time < ?", createTime).Delete(model.Footprint{})
	this.PanicError(db.Error)
}

func (this *FootprintDao) DeleteByUserUuid(userUuid string) {

	db := core.CONTEXT.GetDB().Where("user_uuid = ?", userUuid).Delete(model.Footprint{})
	this.PanicError(db.Error)

}

//System cleanup.
func (this *FootprintDao) Cleanup() {
	this.Logger.Info("[FootprintDao][DownloadTokenDao] clean up. Delete all Footprint")
	db := core.CONTEXT.GetDB().Where("uuid is not null").Delete(model.Footprint{})
	this.PanicError(db.Error)
}
