package core

type BaseBean struct {
	Logger Logger
}

func (this *BaseBean) Init() {
	this.Logger = LOGGER
}

func (this *BaseBean) Bootstrap() {

}

//clean up the application.
func (this *BaseBean) Cleanup() {

}

//shortcut for panic check.
func (this *BaseBean) PanicError(err error) {
	PanicError(err)
}
