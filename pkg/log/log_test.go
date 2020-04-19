package log

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	fmt.Println("------------------------------")
	Debugf("test debug", "add debug")
	i := 1
	for {
		Errorf("test error"+strconv.Itoa(i), "add error")

		i++
		time.Sleep(time.Second * 1)
	}

}
