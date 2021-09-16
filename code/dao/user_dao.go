package dao

import (
	"github.com/eyebluecn/tank/code/constant"
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/model"
	"github.com/eyebluecn/tank/code/tool/builder"
	"github.com/eyebluecn/tank/code/tool/result"
	"github.com/eyebluecn/tank/code/tool/uuid"
	"math"
	"time"
)

type UserDao struct {
	BaseDao
}

func (this *UserDao) Init() {
	this.BaseDao.Init()
}

func (this *UserDao) Create(user *model.User) *model.User {

	if user == nil {
		panic(result.BadRequest("user cannot be nil"))
	}

	timeUUID, _ := uuid.NewV4()
	user.Uuid = string(timeUUID.String())
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	user.LastTime = time.Now()
	user.Sort = time.Now().UnixNano() / 1e6

	db := core.CONTEXT.GetDB().Create(user)
	this.PanicError(db.Error)

	return user
}

//find by uuid. if not found return nil.
func (this *UserDao) FindByUuid(uuid string) *model.User {
	var entity = &model.User{}
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
func (this *UserDao) CheckByUuid(uuid string) *model.User {
	entity := this.FindByUuid(uuid)
	if entity == nil {
		panic(result.NotFound("not found record with uuid = %s", uuid))
	}
	return entity
}

func (this *UserDao) FindByUsername(username string) *model.User {

	var user = &model.User{}
	db := core.CONTEXT.GetDB().Where(&model.User{Username: username}).First(user)
	if db.Error != nil {
		if db.Error.Error() == result.DB_ERROR_NOT_FOUND {
			return nil
		} else {
			panic(db.Error)
		}
	}
	return user
}

func (this *UserDao) Page(page int, pageSize int, username string, status string, sortArray []builder.OrderPair) *model.Pager {

	count, users := this.PlainPage(page, pageSize, username, status, sortArray)

	pager := model.NewPager(page, pageSize, count, users)

	return pager
}

func (this *UserDao) PlainPage(page int, pageSize int, username string, status string, sortArray []builder.OrderPair) (int, []*model.User) {

	var wp = &builder.WherePair{}

	if username != "" {
		wp = wp.And(&builder.WherePair{Query: "username LIKE ?", Args: []interface{}{"%" + username + "%"}})
	}

	if status != "" {
		wp = wp.And(&builder.WherePair{Query: "status = ?", Args: []interface{}{status}})
	}

	count := 0
	db := core.CONTEXT.GetDB().Model(&model.User{}).Where(wp.Query, wp.Args...).Count(&count)
	this.PanicError(db.Error)

	var users []*model.User
	orderStr := this.GetSortString(sortArray)
	if orderStr == "" {
		db = core.CONTEXT.GetDB().Where(wp.Query, wp.Args...).Offset(page * pageSize).Limit(pageSize).Find(&users)
	} else {
		db = core.CONTEXT.GetDB().Where(wp.Query, wp.Args...).Order(orderStr).Offset(page * pageSize).Limit(pageSize).Find(&users)
	}

	this.PanicError(db.Error)

	return count, users
}

//handle user page by page.
func (this *UserDao) PageHandle(username string, status string, fun func(user *model.User)) {

	//delete share and bridges.
	pageSize := 1000
	sortArray := []builder.OrderPair{
		{
			Key:   "uuid",
			Value: constant.DIRECTION_ASC,
		},
	}
	count, _ := this.PlainPage(0, pageSize, username, status, sortArray)
	if count > 0 {
		var totalPages = int(math.Ceil(float64(count) / float64(pageSize)))
		var page int
		for page = 0; page < totalPages; page++ {
			_, users := this.PlainPage(0, pageSize, username, status, sortArray)
			for _, u := range users {
				fun(u)
			}
		}
	}
}

func (this *UserDao) CountByUsername(username string) int {
	var count int
	db := core.CONTEXT.GetDB().
		Model(&model.User{}).
		Where("username = ?", username).
		Count(&count)
	this.PanicError(db.Error)
	return count
}

func (this *UserDao) Save(user *model.User) *model.User {

	user.UpdateTime = time.Now()
	db := core.CONTEXT.GetDB().
		Save(user)
	this.PanicError(db.Error)
	return user
}

//find all 2.0 users.
func (this *UserDao) FindUsers20() []*model.User {
	var users []*model.User
	var wp = &builder.WherePair{}
	wp = wp.And(&builder.WherePair{Query: "username like ?", Args: []interface{}{"%_20"}})

	db := core.CONTEXT.GetDB().Model(&model.User{}).Where(wp.Query, wp.Args...).Find(&users)
	this.PanicError(db.Error)
	return users
}

func (this *UserDao) DeleteUsers20() {
	var wp = &builder.WherePair{}
	wp = wp.And(&builder.WherePair{Query: "username like ?", Args: []interface{}{"%_20"}})

	db := core.CONTEXT.GetDB().Where(wp.Query, wp.Args...).Delete(model.User{})
	this.PanicError(db.Error)
}

func (this *UserDao) Delete(user *model.User) {

	db := core.CONTEXT.GetDB().Delete(&user)
	this.PanicError(db.Error)
}

//System cleanup.
func (this *UserDao) Cleanup() {
	this.Logger.Info("[UserDao] clean up. Delete all User")
	db := core.CONTEXT.GetDB().Where("uuid is not null and role != ?", constant.USER_ROLE_ADMINISTRATOR).Delete(model.User{})
	this.PanicError(db.Error)
}
