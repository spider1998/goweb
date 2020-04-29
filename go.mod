module goweb

go 1.12

replace github.com/go-ozzo/ozzo-routing v2.1.4+incompatible => github.com/caeret/ozzo-routing v2.1.5-0.20181126103820-32c086104c57+incompatible

require (
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535 // indirect
	github.com/go-ozzo/ozzo-routing v2.1.4+incompatible
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/core v0.6.0
	github.com/go-xorm/xorm v0.7.1
	github.com/gobuffalo/packr v1.30.1
	github.com/golang/gddo v0.0.0-20200324184333-3c2cc9a6329d // indirect
	github.com/mediocregopher/radix.v2 v0.0.0-20181115013041-b67df6e626f9
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/nsqio/go-nsq v1.0.8
	github.com/pkg/errors v0.8.1
	github.com/robfig/cron v1.2.0
	github.com/rubenv/sql-migrate v0.0.0-20200423171638-eef9d3b68125
	golang.org/x/crypto v0.0.0-20200423211502-4bdfaf469ed5 // indirect
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
	golang.org/x/sys v0.0.0-20200331124033-c3d80250170d // indirect
	google.golang.org/appengine v1.6.6 // indirect
	gopkg.in/yaml.v2 v2.2.8
)
