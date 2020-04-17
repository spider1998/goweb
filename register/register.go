package register

import "goweb/engine"

func Init() {
	engine.Register("ping", ping)
}
