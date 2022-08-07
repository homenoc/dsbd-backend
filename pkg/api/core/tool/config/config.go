package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const ToMainSlackNotify = "main"
const ToPaymentSlackNotify = "payment"
const ToPaymentLogSlackNotify = "payment_log"

type Config struct {
	Controller Controller `json:"controller"`
	Web        Web        `json:"web"`
	DB         DB         `json:"db"`
	Stripe     Stripe     `json:"stripe"`
	Mail       MailServer `json:"mail"`
	JPNIC      JPNIC      `json:"jpnic"`
	Radius     Radius     `json:"radius"`
	Slack      Slack      `json:"slack"`
	Log        Log        `json:"log"`
	Template   Template   `json:"template"`
}

type Controller struct {
	User     User   `json:"user"`
	Admin    Admin  `json:"admin"`
	Auth     Auth   `json:"auth"`
	TimeZone string `json:"timezone"`
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

type Web struct {
	URL string `json:"url"`
}

type Stripe struct {
	WebhookSecretKey        string `json:"webhook_secret_key"`
	SecretKey               string `json:"secret_key"`
	MembershipConfiguration string `json:"membership_configuration"`
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

type MailServer struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	From     string `json:"from"`
	CC       string `json:"cc"`
	Contract string `json:"contract"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
}

type JPNIC struct {
	URL            string `json:"url"`
	V4CertFilePath string `json:"v4_cert_file_path"`
	V4KeyFilePath  string `json:"v4_key_file_path"`
	V6CertFilePath string `json:"v6_cert_file_path"`
	V6KeyFilePath  string `json:"v6_key_file_path"`
	CAFilePath     string `json:"ca_file_path"`
}

type Radius struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Slack struct {
	Token    string        `json:"token"`
	Channels SlackChannels `json:"channels"`
}

type SlackChannels struct {
	Main       string `json:"main"`
	Log        string `json:"log"`
	Payment    string `json:"payment"`
	PaymentLog string `json:"payment_log"`
}

type Template struct {
	Membership []Membership `json:"membership"`
	NTT        []string     `json:"NTT"`
	V4         []string     `json:"v4"`
	V6         []string     `json:"v6"`
	V4Route    []string     `json:"v4_route"`
	V6Route    []string     `json:"v6_route"`
	Mail       []Mail       `json:"mail"`
}

type Membership struct {
	Title   string `json:"title"`
	Plan    string `json:"plan"`
	PriceID string `json:"price_id"`
	Fee     string `json:"fee"`
}

type Mail struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Log struct {
	Path string `json:"path"`
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
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}
	Conf = data
	return nil
}
