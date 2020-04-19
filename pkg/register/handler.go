package register

import (
	"fmt"
	"goweb/pkg/code"
)

func ping() (codes code.Code, err error) {
	fmt.Println("OK!")
	return code.StatusOk, nil
}
func version() (codes code.Code, err error) {
	fmt.Println("0.1")
	return
}

func truncateLog() {

}
