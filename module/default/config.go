package config

import (
	inters "ehmanager/module/datatypes"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

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
			Address: "0.0.0.0:80",
		},

		Database: inters.DBConfig{
			Master: inters.DBNodeConfig{User: "root",
				Password: "123456",
				Addr:     "127.0.0.1:3306"},
			Replica: inters.DBNodeConfig{User: "root",
				Password: "123456",
				Addr:     "127.0.0.2:3306"},
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
