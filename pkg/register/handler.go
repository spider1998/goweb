package register

import (
	"fmt"

	"goweb/pkg/code"
	"goweb/pkg/components"
	"goweb/pkg/util"
)

func ping() (codes code.Code, err error) {
	fmt.Println("OK!")
	return code.StatusOk, nil
}
func version() (codes code.Code, err error) {
	fmt.Println("0.1")
	return
}

func truncateLog() (code code.Code, err error) {
	logConfig := components.Conf.Log
	logSize, err := util.GetDirSize(logConfig.RuntimePath)
	if err != nil {

	}
}
