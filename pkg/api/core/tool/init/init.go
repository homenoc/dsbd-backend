package init

import (
	"encoding/json"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	dbNOCTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbConnectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/connection/v0"
	dbNTTTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ntt/v0"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
	"io/ioutil"
	"log"
)

type Template struct {
	Service    []core.ServiceTemplate    `json:"service"`
	Connection []core.ConnectionTemplate `json:"connection"`
	NTT        []core.NTTTemplate        `json:"ntt"`
	NOC        []core.NOC                `json:"noc"`
}

func RegisterTemplateConfig(inputTemplatePath string) error {
	configPath := "./template.json"
	if inputTemplatePath != "" {
		configPath = inputTemplatePath
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	var data Template
	json.Unmarshal(file, &data)
	for _, tmp := range data.Service {
		_, err = dbServiceTemplate.Create(&tmp)
		if err != nil {
			log.Println(err)
		}
	}

	for _, tmp := range data.Connection {
		_, err = dbConnectionTemplate.Create(&tmp)
		if err != nil {
			log.Println(err)
		}
	}

	for _, tmp := range data.NTT {
		_, err = dbNTTTemplate.Create(&tmp)
		if err != nil {
			log.Println(err)
		}
	}

	for _, tmp := range data.NOC {
		_, err = dbNOCTemplate.Create(&tmp)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
