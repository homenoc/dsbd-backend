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
}

type Controller struct {
	User  User  `json:"user"`
	Admin Admin `json:"admin"`
}

type User struct {
	Url  string `json:"url"`
	Port int    `json:"port"`
}

type Admin struct {
	Port int `json:"port"`
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
