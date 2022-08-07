package init

import (
	"encoding/json"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	dbNOCTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbConnectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/connection/v0"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
	"io/ioutil"
)

type Template struct {
	Service    []core.ServiceTemplate    `json:"service"`
	Connection []core.ConnectionTemplate `json:"connection"`
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
	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	for _, tmp := range data.Service {
		_, err = dbServiceTemplate.Create(&tmp)
		if err != nil {
			return err
		}
	}

	for _, tmp := range data.Connection {
		_, err = dbConnectionTemplate.Create(&tmp)
		if err != nil {
			return err
		}
	}

	for _, tmp := range data.NOC {
		_, err = dbNOCTemplate.Create(&tmp)
		if err != nil {
			return err
		}
	}

	return nil
}
