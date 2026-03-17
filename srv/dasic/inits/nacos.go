package inits

import (
	"fmt"
	"lianxi/srv/dasic/config"

	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

func NacosInit() {
	clientConfig := constant.ClientConfig{
		NamespaceId:         config.Gen.Nacos.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: config.Gen.Nacos.Addr,
			Port:   uint64(config.Gen.Nacos.Port),
		},
	}

	nacosClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig":  clientConfig,
		"serverConfigs": serverConfigs,
	})
	if err != nil {
		fmt.Printf("创建 Nacos 客户端失败: %v\n", err)
		return
	}
	configContent, err := nacosClient.GetConfig(vo.ConfigParam{
		DataId: config.Gen.Nacos.DataID,
		Group:  config.Gen.Nacos.Group,
	})
	if err != nil {
		fmt.Printf("从 Nacos 获取配置失败: %v\n", err)
		return
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(strings.NewReader(configContent))
	if err != nil {
		fmt.Printf("解析 Nacos 配置失败: %v\n", err)
		return
	}
	err = viper.Unmarshal(config.Gen)
	if err != nil {
		fmt.Printf("反序列化 Nacos 配置失败: %v\n", err)
		return
	}

	fmt.Println("Nacos 配置读取成功")
}
