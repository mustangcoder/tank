package rest

import (
	"github.com/eyebluecn/tank/code/constant"
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/dao"
	"github.com/eyebluecn/tank/code/model"
	"github.com/eyebluecn/tank/code/service"
	"github.com/eyebluecn/tank/code/tool/builder"
	"github.com/eyebluecn/tank/code/tool/i18n"
	"github.com/eyebluecn/tank/code/tool/result"
	"net/http"
	"strconv"
	"strings"
)

type MatterController struct {
	BaseController
	matterDao         *dao.MatterDao
	matterService     *service.MatterService
	preferenceService *service.PreferenceService
	downloadTokenDao  *dao.DownloadTokenDao
	imageCacheDao     *dao.ImageCacheDao
	shareDao          *dao.ShareDao
	shareService      *service.ShareService
	bridgeDao         *dao.BridgeDao
	imageCacheService *service.ImageCacheService
	userService       *service.UserService
}

func (this *MatterController) Init() {
	this.BaseController.Init()

	b := core.CONTEXT.GetBean(this.userService)
	if b, ok := b.(*service.UserService); ok {
		this.userService = b
	}

	b = core.CONTEXT.GetBean(this.matterDao)
	if b, ok := b.(*dao.MatterDao); ok {
		this.matterDao = b
	}

	b = core.CONTEXT.GetBean(this.matterService)
	if b, ok := b.(*service.MatterService); ok {
		this.matterService = b
	}

	b = core.CONTEXT.GetBean(this.downloadTokenDao)
	if b, ok := b.(*dao.DownloadTokenDao); ok {
		this.downloadTokenDao = b
	}

	b = core.CONTEXT.GetBean(this.imageCacheDao)
	if b, ok := b.(*dao.ImageCacheDao); ok {
		this.imageCacheDao = b
	}

	b = core.CONTEXT.GetBean(this.shareDao)
	if b, ok := b.(*dao.ShareDao); ok {
		this.shareDao = b
	}

	b = core.CONTEXT.GetBean(this.shareService)
	if b, ok := b.(*service.ShareService); ok {
		this.shareService = b
	}

	b = core.CONTEXT.GetBean(this.bridgeDao)
	if b, ok := b.(*dao.BridgeDao); ok {
		this.bridgeDao = b
	}

	b = core.CONTEXT.GetBean(this.imageCacheService)
	if b, ok := b.(*service.ImageCacheService); ok {
		this.imageCacheService = b
	}
}

func (this *MatterController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {

	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))

	routeMap["/api/matter/create/directory"] = this.Wrap(this.CreateDirectory, constant.USER_ROLE_USER)
	routeMap["/api/matter/upload"] = this.Wrap(this.Upload, constant.USER_ROLE_USER)
	routeMap["/api/matter/crawl"] = this.Wrap(this.Crawl, constant.USER_ROLE_USER)
	routeMap["/api/matter/soft/delete"] = this.Wrap(this.SoftDelete, constant.USER_ROLE_USER)
	routeMap["/api/matter/soft/delete/batch"] = this.Wrap(this.SoftDeleteBatch, constant.USER_ROLE_USER)
	routeMap["/api/matter/recovery"] = this.Wrap(this.Recovery, constant.USER_ROLE_USER)
	routeMap["/api/matter/recovery/batch"] = this.Wrap(this.RecoveryBatch, constant.USER_ROLE_USER)
	routeMap["/api/matter/delete"] = this.Wrap(this.Delete, constant.USER_ROLE_USER)
	routeMap["/api/matter/delete/batch"] = this.Wrap(this.DeleteBatch, constant.USER_ROLE_USER)
	routeMap["/api/matter/clean/expired/deleted/matters"] = this.Wrap(this.CleanExpiredDeletedMatters, constant.USER_ROLE_ADMINISTRATOR)
	routeMap["/api/matter/rename"] = this.Wrap(this.Rename, constant.USER_ROLE_USER)
	routeMap["/api/matter/change/privacy"] = this.Wrap(this.ChangePrivacy, constant.USER_ROLE_USER)
	routeMap["/api/matter/move"] = this.Wrap(this.Move, constant.USER_ROLE_USER)
	routeMap["/api/matter/detail"] = this.Wrap(this.Detail, constant.USER_ROLE_USER)
	routeMap["/api/matter/page"] = this.Wrap(this.Page, constant.USER_ROLE_GUEST)

	//mirror local files.
	routeMap["/api/matter/mirror"] = this.Wrap(this.Mirror, constant.USER_ROLE_USER)
	routeMap["/api/matter/zip"] = this.Wrap(this.Zip, constant.USER_ROLE_USER)

	return routeMap
}

func (this *MatterController) Detail(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuid := request.FormValue("uuid")
	if uuid == "" {
		panic(result.BadRequest("uuid cannot be null"))
	}

	matter := this.matterService.Detail(request, uuid)

	user := this.userService.CheckUser(request)
	if matter.UserUuid != user.Uuid {
		panic(result.UNAUTHORIZED)
	}

	return this.Success(matter)

}

func (this *MatterController) Page(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	pageStr := request.FormValue("page")
	pageSizeStr := request.FormValue("pageSize")
	orderCreateTime := request.FormValue("orderCreateTime")
	orderUpdateTime := request.FormValue("orderUpdateTime")
	orderDeleteTime := request.FormValue("orderDeleteTime")
	orderSort := request.FormValue("orderSort")
	orderTimes := request.FormValue("orderTimes")

	puuid := request.FormValue("puuid")
	name := request.FormValue("name")
	dir := request.FormValue("dir")
	deleted := request.FormValue("deleted")
	orderDir := request.FormValue("orderDir")
	orderSize := request.FormValue("orderSize")
	orderName := request.FormValue("orderName")
	extensionsStr := request.FormValue("extensions")

	var userUuid string

	//auth by shareUuid.
	shareUuid := request.FormValue("shareUuid")
	shareCode := request.FormValue("shareCode")
	shareRootUuid := request.FormValue("shareRootUuid")
	if shareUuid != "" {

		if puuid == "" {
			panic(result.BadRequest("puuid cannot be null"))
		}

		dirMatter := this.matterDao.CheckByUuid(puuid)
		if !dirMatter.Dir {
			panic(result.BadRequest("puuid is not a directory"))
		}

		user := this.userService.FindUser(request)

		this.shareService.ValidateMatter(request, shareUuid, shareCode, user, shareRootUuid, dirMatter)
		puuid = dirMatter.Uuid

	} else {
		//if cannot auth by share. Then login is required.
		user := this.userService.CheckUser(request)
		userUuid = user.Uuid

	}

	var page int
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	pageSize := 200
	if pageSizeStr != "" {
		tmp, err := strconv.Atoi(pageSizeStr)
		if err == nil {
			pageSize = tmp
		}
	}

	var extensions []string
	if extensionsStr != "" {
		extensions = strings.Split(extensionsStr, ",")
	}

	sortArray := []builder.OrderPair{
		{
			Key:   "dir",
			Value: orderDir,
		},
		{
			Key:   "create_time",
			Value: orderCreateTime,
		},
		{
			Key:   "update_time",
			Value: orderUpdateTime,
		},
		{
			Key:   "delete_time",
			Value: orderDeleteTime,
		},
		{
			Key:   "sort",
			Value: orderSort,
		},
		{
			Key:   "size",
			Value: orderSize,
		},
		{
			Key:   "name",
			Value: orderName,
		},
		{
			Key:   "times",
			Value: orderTimes,
		},
	}

	pager := this.matterDao.Page(page, pageSize, puuid, userUuid, name, dir, deleted, extensions, sortArray)

	return this.Success(pager)
}

func (this *MatterController) CreateDirectory(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	puuid := request.FormValue("puuid")
	name := request.FormValue("name")

	user := this.userService.CheckUser(request)

	var dirMatter = this.matterDao.CheckWithRootByUuid(puuid, user)

	matter := this.matterService.AtomicCreateDirectory(request, dirMatter, name, user)
	return this.Success(matter)
}

func (this *MatterController) Upload(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	puuid := request.FormValue("puuid")
	privacyStr := request.FormValue("privacy")
	file, handler, err := request.FormFile("file")
	this.PanicError(err)
	defer func() {
		err := file.Close()
		this.PanicError(err)
	}()

	user := this.userService.CheckUser(request)

	privacy := privacyStr == constant.TRUE

	err = request.ParseMultipartForm(32 << 20)
	this.PanicError(err)

	//for IE browser. filename may contains filepath.
	fileName := handler.Filename
	pos := strings.LastIndex(fileName, "\\")
	if pos != -1 {
		fileName = fileName[pos+1:]
	}
	pos = strings.LastIndex(fileName, "/")
	if pos != -1 {
		fileName = fileName[pos+1:]
	}

	dirMatter := this.matterDao.CheckWithRootByUuid(puuid, user)

	//support upload simultaneously
	matter := this.matterService.Upload(request, file, user, dirMatter, fileName, privacy)

	return this.Success(matter)
}

//crawl a file by url.
func (this *MatterController) Crawl(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	url := request.FormValue("url")
	destPath := request.FormValue("destPath")
	filename := request.FormValue("filename")

	user := this.userService.CheckUser(request)

	dirMatter := this.matterService.CreateDirectories(request, user, destPath)

	if url == "" || (!strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://")) {
		panic(" url must start with  http:// or https://")
	}

	if filename == "" {
		panic("filename cannot be null")
	}

	matter := this.matterService.AtomicCrawl(request, url, filename, user, dirMatter, true)

	return this.Success(matter)
}

//soft delete.
func (this *MatterController) SoftDelete(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuid := request.FormValue("uuid")
	if uuid == "" {
		panic(result.BadRequest("uuid cannot be null"))
	}

	matter := this.matterDao.CheckByUuid(uuid)

	user := this.userService.CheckUser(request)
	if matter.UserUuid != user.Uuid {
		panic(result.UNAUTHORIZED)
	}

	this.matterService.AtomicSoftDelete(request, matter, user)

	return this.Success("OK")
}

func (this *MatterController) SoftDeleteBatch(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuids := request.FormValue("uuids")
	if uuids == "" {
		panic(result.BadRequest("uuids cannot be null"))
	}
	user := this.userService.CheckUser(request)

	uuidArray := strings.Split(uuids, ",")

	matters := make([]*model.Matter, 0)
	for _, uuid := range uuidArray {

		matter := this.matterDao.FindByUuid(uuid)

		if matter == nil {
			this.Logger.Warn("%s not exist anymore", uuid)
			continue
		}

		if matter.UserUuid != user.Uuid {
			panic(result.UNAUTHORIZED)
		}

		matters = append(matters, matter)
	}

	for _, matter := range matters {
		this.matterService.AtomicSoftDelete(request, matter, user)
	}

	return this.Success("OK")
}

//recovery delete.
func (this *MatterController) Recovery(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuid := request.FormValue("uuid")
	if uuid == "" {
		panic(result.BadRequest("uuid cannot be null"))
	}

	matter := this.matterDao.CheckByUuid(uuid)

	user := this.userService.CheckUser(request)
	if matter.UserUuid != user.Uuid {
		panic(result.UNAUTHORIZED)
	}

	this.matterService.AtomicRecovery(request, matter, user)

	return this.Success("OK")
}

//recovery batch.
func (this *MatterController) RecoveryBatch(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuids := request.FormValue("uuids")
	if uuids == "" {
		panic(result.BadRequest("uuids cannot be null"))
	}

	uuidArray := strings.Split(uuids, ",")

	for _, uuid := range uuidArray {

		matter := this.matterDao.FindByUuid(uuid)

		if matter == nil {
			this.Logger.Warn("%s not exist anymore", uuid)
			continue
		}

		user := this.userService.CheckUser(request)
		if matter.UserUuid != user.Uuid {
			panic(result.UNAUTHORIZED)
		}

		this.matterService.AtomicRecovery(request, matter, user)

	}

	return this.Success("OK")
}

//complete delete.
func (this *MatterController) Delete(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuid := request.FormValue("uuid")
	if uuid == "" {
		panic(result.BadRequest("uuid cannot be null"))
	}

	matter := this.matterDao.CheckByUuid(uuid)

	user := this.userService.CheckUser(request)
	if matter.UserUuid != user.Uuid {
		panic(result.UNAUTHORIZED)
	}

	this.matterService.AtomicDelete(request, matter, user)

	return this.Success("OK")
}

func (this *MatterController) DeleteBatch(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuids := request.FormValue("uuids")
	if uuids == "" {
		panic(result.BadRequest("uuids cannot be null"))
	}

	uuidArray := strings.Split(uuids, ",")
	user := this.userService.CheckUser(request)
	matters := make([]*model.Matter, 0)
	for _, uuid := range uuidArray {

		matter := this.matterDao.FindByUuid(uuid)

		if matter == nil {
			this.Logger.Warn("%s not exist anymore", uuid)
			continue
		}

		if matter.UserUuid != user.Uuid {
			panic(result.UNAUTHORIZED)
		}

		matters = append(matters, matter)
	}

	for _, matter := range matters {

		this.matterService.AtomicDelete(request, matter, user)
	}

	return this.Success("OK")
}

//manual clean expired deleted matters.
func (this *MatterController) CleanExpiredDeletedMatters(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	this.matterService.CleanExpiredDeletedMatters()

	return this.Success("OK")
}

func (this *MatterController) Rename(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuid := request.FormValue("uuid")
	name := request.FormValue("name")

	user := this.userService.CheckUser(request)

	matter := this.matterDao.CheckByUuid(uuid)

	if matter.UserUuid != user.Uuid {
		panic(result.UNAUTHORIZED)
	}

	this.matterService.AtomicRename(request, matter, name, false, user)

	return this.Success(matter)
}

func (this *MatterController) ChangePrivacy(writer http.ResponseWriter, request *http.Request) *result.WebResult {
	uuid := request.FormValue("uuid")
	privacyStr := request.FormValue("privacy")
	privacy := false
	if privacyStr == constant.TRUE {
		privacy = true
	}

	matter := this.matterDao.CheckByUuid(uuid)

	if matter.Deleted {
		panic(result.BadRequest("matter has been deleted. Cannot change privacy."))
	}

	if matter.Privacy == privacy {
		panic(result.BadRequest("not changed. Invalid operation."))
	}

	user := this.userService.CheckUser(request)
	if matter.UserUuid != user.Uuid {
		panic(result.UNAUTHORIZED)
	}

	matter.Privacy = privacy
	this.matterDao.Save(matter)

	return this.Success("OK")
}

func (this *MatterController) Move(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	srcUuidsStr := request.FormValue("srcUuids")
	destUuid := request.FormValue("destUuid")

	var srcUuids []string
	if srcUuidsStr == "" {
		panic(result.BadRequest("srcUuids cannot be null"))
	} else {
		srcUuids = strings.Split(srcUuidsStr, ",")
	}

	user := this.userService.CheckUser(request)

	var destMatter = this.matterDao.CheckWithRootByUuid(destUuid, user)
	if !destMatter.Dir {
		panic(result.BadRequest("destination is not a directory"))
	}

	if destMatter.UserUuid != user.Uuid {
		panic(result.UNAUTHORIZED)
	}

	if destMatter.Deleted {
		panic(result.BadRequest("dest matter has been deleted. Cannot move."))
	}

	var srcMatters []*model.Matter
	for _, uuid := range srcUuids {
		srcMatter := this.matterDao.CheckByUuid(uuid)

		if srcMatter.Puuid == destMatter.Uuid {
			panic(result.BadRequest("no move, invalid operation"))
		}

		if srcMatter.Deleted {
			panic(result.BadRequest("src matter has been deleted. Cannot move."))
		}

		//check whether there are files with the same name.
		count := this.matterDao.CountByUserUuidAndPuuidAndDirAndName(user.Uuid, destMatter.Uuid, srcMatter.Dir, srcMatter.Name)

		if count > 0 {
			panic(result.BadRequestI18n(request, i18n.MatterExist, srcMatter.Name))
		}

		if srcMatter.UserUuid != destMatter.UserUuid {
			panic("owner not the same")
		}

		srcMatters = append(srcMatters, srcMatter)
	}

	this.matterService.AtomicMoveBatch(request, srcMatters, destMatter, user)

	return this.Success(nil)
}

//mirror local files to EyeblueTank
func (this *MatterController) Mirror(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	srcPath := request.FormValue("srcPath")
	destPath := request.FormValue("destPath")
	overwriteStr := request.FormValue("overwrite")

	if srcPath == "" {
		panic(result.BadRequest("srcPath cannot be null"))
	}

	overwrite := false
	if overwriteStr == constant.TRUE {
		overwrite = true
	}

	user := this.userService.CheckUser(request)

	this.matterService.AtomicMirror(request, srcPath, destPath, overwrite, user)

	return this.Success(nil)

}

//download zip.
func (this *MatterController) Zip(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuids := request.FormValue("uuids")
	if uuids == "" {
		panic(result.BadRequest("uuids cannot be null"))
	}

	uuidArray := strings.Split(uuids, ",")

	matters := this.matterDao.FindByUuids(uuidArray, nil)

	if matters == nil || len(matters) == 0 {
		panic(result.BadRequest("matters cannot be nil."))
	}

	for _, matter := range matters {
		if matter.Deleted {
			panic(result.BadRequest("matter has been deleted. Cannot download batch."))
		}
	}

	user := this.userService.CheckUser(request)
	puuid := matters[0].Puuid

	for _, m := range matters {
		if m.UserUuid != user.Uuid {
			panic(result.UNAUTHORIZED)
		} else if m.Puuid != puuid {
			panic(result.BadRequest("puuid not same"))
		}
	}

	this.matterService.DownloadZip(writer, request, matters)

	return nil
}
