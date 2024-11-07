package models

import (
	"log"

	"github.com/aungmyozaw92/go-restapi-sta/config"
)

func MigrateTable() {
	db := config.GetDB()

	err := db.AutoMigrate(
		&User{},
		&Role{},
		&Module{},
		&RoleModule{},
		// &Unit{},
		// &Category{},
		// &Image{},
		// &Product{},
	)
	if err != nil {
		log.Fatal(err)
	}
}