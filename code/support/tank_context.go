package support

import (
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/dao"
	"github.com/eyebluecn/tank/code/rest"
	"github.com/eyebluecn/tank/code/service"
	"github.com/eyebluecn/tank/code/tool/cache"
	"github.com/jinzhu/gorm"
	"net/http"
	"reflect"
)

type TankContext struct {
	//db connection
	db *gorm.DB
	//session cache
	SessionCache *cache.Table
	//bean map.
	BeanMap map[string]core.Bean
	//controller map
	ControllerMap map[string]core.Controller
	//router
	Router *TankRouter
}

func (this *TankContext) Init() {

	//create session cache
	this.SessionCache = cache.NewTable()

	//init map
	this.BeanMap = make(map[string]core.Bean)
	this.ControllerMap = make(map[string]core.Controller)

	//register beans. This method will put Controllers to ControllerMap.
	this.registerBeans()

	//init every bean.
	this.initBeans()

	//create and init router.
	this.Router = NewRouter()

	//if the application is installed. Bean's Bootstrap method will be invoked.
	this.InstallOk()

}

func (this *TankContext) GetDB() *gorm.DB {
	return this.db
}

func (this *TankContext) GetSessionCache() *cache.Table {
	return this.SessionCache
}

func (this *TankContext) GetControllerMap() map[string]core.Controller {
	return this.ControllerMap
}

func (this *TankContext) Cleanup() {
	for _, bean := range this.BeanMap {
		bean.Cleanup()
	}
}

//can serve as http server.
func (this *TankContext) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	this.Router.ServeHTTP(writer, request)
}

func (this *TankContext) OpenDb() {

	var err error = nil
	this.db, err = gorm.Open("mysql", core.CONFIG.MysqlUrl())

	if err != nil {
		core.LOGGER.Panic("failed to connect mysql database")
	}

	//whether open the db sql log. (only true when debug)
	this.db.LogMode(false)
}

func (this *TankContext) CloseDb() {

	if this.db != nil {
		err := this.db.Close()
		if err != nil {
			core.LOGGER.Error("occur error when closing db %s", err.Error())
		}
	}
}

func (this *TankContext) registerBean(bean core.Bean) {

	typeOf := reflect.TypeOf(bean)
	typeName := typeOf.String()

	if element, ok := bean.(core.Bean); ok {

		if _, ok := this.BeanMap[typeName]; ok {
			core.LOGGER.Error("%s has been registerd, skip", typeName)
		} else {
			this.BeanMap[typeName] = element

			//if is controller type, put into ControllerMap
			if controller, ok1 := bean.(core.Controller); ok1 {
				this.ControllerMap[typeName] = controller
			}

		}

	} else {
		core.LOGGER.Panic("%s is not the Bean type", typeName)
	}

}

func (this *TankContext) registerBeans() {

	//alien
	this.registerBean(new(rest.AlienController))
	this.registerBean(new(service.AlienService))

	//bridge
	this.registerBean(new(dao.BridgeDao))
	this.registerBean(new(service.BridgeService))

	//dashboard
	this.registerBean(new(rest.DashboardController))
	this.registerBean(new(dao.DashboardDao))
	this.registerBean(new(service.DashboardService))

	//downloadToken
	this.registerBean(new(dao.DownloadTokenDao))

	//imageCache
	this.registerBean(new(rest.ImageCacheController))
	this.registerBean(new(dao.ImageCacheDao))
	this.registerBean(new(service.ImageCacheService))

	//install
	this.registerBean(new(rest.InstallController))

	//matter
	this.registerBean(new(rest.MatterController))
	this.registerBean(new(dao.MatterDao))
	this.registerBean(new(service.MatterService))

	//preference
	this.registerBean(new(rest.PreferenceController))
	this.registerBean(new(dao.PreferenceDao))
	this.registerBean(new(service.PreferenceService))

	//footprint
	this.registerBean(new(rest.FootprintController))
	this.registerBean(new(dao.FootprintDao))
	this.registerBean(new(service.FootprintService))

	//session
	this.registerBean(new(dao.SessionDao))
	this.registerBean(new(service.SessionService))

	//share
	this.registerBean(new(rest.ShareController))
	this.registerBean(new(dao.ShareDao))
	this.registerBean(new(service.ShareService))

	//uploadToken
	this.registerBean(new(dao.UploadTokenDao))

	//task
	this.registerBean(new(service.TaskService))

	//user
	this.registerBean(new(rest.UserController))
	this.registerBean(new(dao.UserDao))
	this.registerBean(new(service.UserService))

	//webdav
	this.registerBean(new(rest.DavController))
	this.registerBean(new(service.DavService))

	//healthCheck
	this.registerBean(new(rest.ServerStatusController))
}

func (this *TankContext) GetBean(bean core.Bean) core.Bean {

	typeOf := reflect.TypeOf(bean)
	typeName := typeOf.String()

	if val, ok := this.BeanMap[typeName]; ok {
		return val
	} else {
		core.LOGGER.Panic("%s not registered", typeName)
		return nil
	}
}

func (this *TankContext) initBeans() {

	for _, bean := range this.BeanMap {
		bean.Init()
	}
}

//if application installed. invoke this method.
func (this *TankContext) InstallOk() {

	if core.CONFIG.Installed() {
		this.OpenDb()

		for _, bean := range this.BeanMap {
			bean.Bootstrap()
		}
	}

}

func (this *TankContext) Destroy() {
	this.CloseDb()
}
