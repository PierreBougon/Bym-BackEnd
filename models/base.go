package models

import (
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB //database

func init() {

	dbUrl := os.Getenv("DATABASE_URL")
	dbLang := "postgres"
	dbUri := ""

	if  len(dbUrl) > 3 {
		parts := strings.Split(dbUrl, "://")
		dbLang = parts[0]
		dbUri = parts[1]

	} else {
		e := godotenv.Load() //Load .env file
		if e != nil {
			fmt.Print(e)
		}

		username := os.Getenv("db_user")
		password := os.Getenv("db_pass")
		dbName := os.Getenv("db_name")
		dbHost := os.Getenv("db_host")

		dbUri = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
		fmt.Println(dbUri)
	}

	conn, err := gorm.Open(dbLang, dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Playlist{}, &Song{}) //Database migration
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
