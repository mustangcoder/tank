package core

type Config interface {
	Installed() bool
	ServerPort() int
	//get the mysql url. eg. tank:tank123@tcp(127.0.0.1:3306)/tank?charset=utf8&parseTime=True&loc=Local
	MysqlUrl() string
	//files storage location.
	MatterPath() string
	//when installed by user. Write configs to tank.json
	FinishInstall(mysqlPort int, mysqlHost string, mysqlSchema string, mysqlUsername string, mysqlPassword string, mysqlCharset string)
}
