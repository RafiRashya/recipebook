package database

import (
	"go-rest-modul/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func init() {
	dsn := "host=localhost port=5432 user=postgres password=123 dbname=recipe_db sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Tidak dapat terhubung ke database")
		panic(err.Error())
	}
	DB = db
	log.Println("Berhasil terhubung ke database")

	err = db.AutoMigrate(&models.Category{}, &models.Recipe{}, &models.Ingredient{}, &models.RecipeIngredient{})
	if err != nil {
		log.Println("Gagal melakukan migrasi")
	}
}
