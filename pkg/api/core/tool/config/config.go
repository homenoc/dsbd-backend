package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Controller Controller `json:"controller"`
	DB         DB         `json:"db"`
	Mail       Mail       `json:"mail"`
	Radius     Radius     `json:"radius"`
	Slack      []Slack    `json:"slack"`
}

type Controller struct {
	User  User  `json:"user"`
	Admin Admin `json:"admin"`
	Auth  Auth  `json:"auth"`
}

type User struct {
	IP   string `json:"ip"`
	Url  string `json:"url"`
	Port int    `json:"port"`
}

type Admin struct {
	IP        string    `json:"ip"`
	Port      int       `json:"port"`
	AdminAuth AdminAuth `json:"auth"`
}

type AdminAuth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Auth struct {
	Token1 string `json:"token1"`
	Token2 string `json:"token2"`
	Token3 string `json:"token3"`
}

type DB struct {
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	DBName string `json:"dbName"`
}

type Mail struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	From     string `json:"from"`
	CC       string `json:"cc"`
	Contract string `json:"contract"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
}

type Radius struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Slack struct {
	ID         string `json:"id"`
	WebHookUrl string `json:"url"`
	Channel    string `json:"channel"`
	Name       string `json:"name"`
}

var Conf Config

func GetConfig(inputConfPath string) error {
	configPath := "./data.json"
	if inputConfPath != "" {
		configPath = inputConfPath
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	var data Config
	json.Unmarshal(file, &data)
	Conf = data
	return nil
}

//func GetConfig(inputConfPath string) (Config, error) {
//	configPath := "./data.json"
//	if inputConfPath != "" {
//		configPath = inputConfPath
//	}
//	file, err := ioutil.ReadFile(configPath)
//	if err != nil {
//		return Config{}, err
//	}
//	var data Config
//	json.Unmarshal(file, &data)
//	return data, nil
//}
