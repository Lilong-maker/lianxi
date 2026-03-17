package inits

import (
	"fmt"
	"lianxi/srv/dasic/config"

	"lianxi/srv/handler/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	ViperInit()
	MysqlInit()
	RedisInit()
	NacosInit()

}

var err error

func MysqlInit() {
	mysqlConfig := config.Gen.Mysql

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database,
	)

	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}
	fmt.Println("数据库连接成功")
	// Auto-migrate multiple domain models to ensure schema is up-to-date
	err = config.DB.AutoMigrate(
		&model.Goods{},
		&model.Order{},
		&model.OrderItem{},
	)
	if err != nil {
		fmt.Printf("表迁移失败: %v\n", err)
		return
	}
	fmt.Println("表迁移成功")
}
func ViperInit() {
	viper.SetConfigFile("C:\\Users\\Lenovo\\Desktop\\lianxi\\config.yml")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config.Gen)
	if err != nil {
		return
	}
	fmt.Println("配置文件加载成功")
}
