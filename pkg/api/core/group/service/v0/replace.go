package v0

import "github.com/homenoc/dsbd-backend/pkg/api/core"

func replaceHiddenFalseForJPNICAdmin(jpnicAdmin *core.JPNICAdmin) {
	jpnicAdmin.Hidden = false
}

func replaceHiddenFalseForJPNICTech(jpnicTechs *[]core.JPNICTech) {
	for _, jpnicTech := range *jpnicTechs {
		jpnicTech.Hidden = false
	}
}
