package dbs

import (
	"e-commerce/pkg/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAutoMigrate(t *testing.T) {
	config, err := utils.LoadConfig("../..")
	require.NoError(t, err)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	err = AutoMigrate(db)
	require.NoError(t, err)
}
