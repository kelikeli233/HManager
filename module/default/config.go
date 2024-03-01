package config

import (
	inters "ehmanager/module/datatypes"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

func getSysNameDefaultData() map[string]string {
	SysNameCache := map[string]string{
		"monitoring":       "容器平台",
		"cloudtogo-system": "容器平台",
		"kube-system":      "容器平台",
		"cps-dev":          "容器平台",
		"cps-dev1":         "容器平台",
		"NA":               "NA",
	}
	return SysNameCache
}

func LoadSysNameDefaultConfig() (map[string]string, error) {
	SysNameConfigFile := "sysname.cfg"
	config, err := readSysnameConfig(SysNameConfigFile)
	if err != nil {
		//读失败 写默认
		log.Println("写入默认系统名称配置:", err)
		config = getSysNameDefaultData()
		err := writeSysnameConfig(SysNameConfigFile, config)
		if err != nil {
			log.Println("写入默认系统名称配置失败:", err)
		}
		return config, nil
	}

	return config, nil
}

func readSysnameConfig(filePath string) (map[string]string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config map[string]string
	err = json.Unmarshal(content, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func writeSysnameConfig(filePath string, config map[string]string) error {
	content, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	//读写权限
	err = ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getStartDefaultConfig() inters.Config {
	return inters.Config{
		HTTP: inters.HTTPConfig{
			Address: "0.0.0.0:9568",
		},

		Database: inters.DBConfig{
			Master:                 inters.DBNodeConfig{DSN: "root:123456@tcp(127.0.0.1:3306)/alertsnitch"},
			Replica:                inters.DBNodeConfig{DSN: "root:123456@tcp(127.0.0.2:3306)/alertsnitch"},
			MaxIdleConns:           5,
			MaxOpenConns:           10,
			MaxConnLifetimeSeconds: 600,
			Backend:                "mysql",
		},

		Other: inters.OtherConfig{
			Debug:   true,
			Dryrun:  false,
			Version: false,
		},
	}
}

func LoadStartConfig() (inters.Config, error) {
	var config inters.Config

	configFileName := "config.yaml"
	//打开配置文件
	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		//配置不存在
		if os.IsNotExist(err) {
			config = getStartDefaultConfig()
			err := writeConfig(configFileName, config)
			if err != nil {
				return config, err
			}
			return config, nil
		}
	}

	err = yaml.Unmarshal(data, &config)
	//Ldap配置更新

	if err != nil {
		return config, err
	}
	return config, nil
}

// 配置写入文件
func writeConfig(filename string, config inters.Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
