package models

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB //database

func migrate() {
	db.AutoMigrate(
		&Account{},
		&Playlist{},
		&Song{},
		)
}

func getDbInfoFromEnv() (dbLang string, dbUri string) {
	dbUrl := os.Getenv("DATABASE_URL")
	if  len(dbUrl) > 4 { // need a regex check to validate
		parts := strings.Split(dbUrl, "://")
		dbLang = parts[0]
		dbUri = parts[1]
	} else {
		e := godotenv.Load() //Load .env file
		if e != nil {
			fmt.Print(e)
		}

		dbLang = os.Getenv("db_dialect")
		username := os.Getenv("db_user")
		password := os.Getenv("db_pass")
		dbName := os.Getenv("db_name")
		dbHost := os.Getenv("db_host")

		dbUri = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s sslmode=require", dbHost, username, dbName, password) //Build connection string
	}
	fmt.Println(dbLang, dbUri)
	return
}

func init() {
	// dbLang, dbUri := getDbInfoFromEnv()
	conn, err := gorm.Open(getDbInfoFromEnv())//(dbLang, dbUri)
	if err != nil {
		fmt.Print(err)
		panic("failed to connect to database")
	}

	db = conn
	// migrate()
	db.Debug().AutoMigrate(&Account{}, &Playlist{}, &Song{}) //Database migration
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
