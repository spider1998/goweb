package components

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/gobuffalo/packr"
	migrate "github.com/rubenv/sql-migrate"
)

func NewDB(dsn string, sources packr.Box, logger Logger) (*xorm.Engine, error) {
	e, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}
	e.SetMapper(core.GonicMapper{})

	e.SetLogger(makeDbLogger(logger, core.LOG_DEBUG))
	e.ShowSQL(true)
	e.ShowExecTime(true)

	n, err := migrateDB(dsn, sources)
	if err != nil {
		return nil, err
	}
	if n > 0 {
		logger.Infof("make db migration.", "n", n)
	}

	return e, nil
}

func migrateDB(dsn string, sources packr.Box) (int, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return 0, err
	}
	config.ParseTime = true
	dbname := config.DBName
	config.DBName = ""
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return 0, err
	}
	defer db.Close()
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		return 0, err
	}
	_, err = db.Exec("USE " + dbname)
	if err != nil {
		return 0, err
	}
	migrations := &migrate.PackrMigrationSource{
		Box: sources,
	}
	migrate.SetTable("schema_migration")
	return migrate.Exec(db, "mysql", migrations, migrate.Up)
}
