package config

type AppConfig struct {
	Mysql
	Redis
	Nacos
}
type Mysql struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}
type Redis struct {
	Host     string
	Port     int
	Password string
	Database int
}

type Nacos struct {
	Addr      string
	Port      int
	Namespace string
	DataID    string
	Group     string
}
