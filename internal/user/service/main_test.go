package service

import (
	"e-commerce/pkg/dbs"
	"e-commerce/pkg/utils"
	"log"
	"os"
	"testing"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("can not load config")
		return
	}
	dsn := dbs.BuildDSN(config)
	testDB, err = dbs.InitDB(dsn)
	if err != nil {
		log.Fatal("can not connect to database")
		return
	}
	os.Exit(m.Run())
}
