package config

type AppConfig struct {
	Mysql
	Redis
	Nacos
	Consul
	AliPay
	//RabbitMQ
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
	//Username  string
	//Password  string
}

type Consul struct {
	Host        string
	Port        int
	ServiceName string
	ServicePort int
	TTL         int
}

//	type RabbitMQ struct {
//		Host     string
//		Port     int
//		User     string
//		Password string
//		Vhost    string
//	}
type AliPay struct {
	PrivateKey string
	AppId      string
	NotifyURL  string
	ReturnURL  string
}
