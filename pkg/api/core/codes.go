package core

import (
	"fmt"
	"strconv"
)

// ServiceCode renders the user-facing service identifier, e.g. "1-IP3B001".
func ServiceCode(groupID uint, serviceType string, serviceNumber uint) string {
	return strconv.Itoa(int(groupID)) + "-" + serviceType + fmt.Sprintf("%03d", serviceNumber)
}

// ConnectionCode appends the connection segment, e.g. "1-IP3B001-EIP001".
func ConnectionCode(serviceCode, connectionType string, connectionNumber uint) string {
	return serviceCode + "-" + connectionType + fmt.Sprintf("%03d", connectionNumber)
}
