package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

var db *gorm.DB //database

func migrate() {
	db.AutoMigrate(
		&Account{},
		&Playlist{},
		&Song{},
	)
}

func getDbInfoFromEnv() (dbDialect string, dbUri string) {
	dbUrl := os.Getenv("DATABASE_URL")
	reg := regexp.MustCompile("^(postgres|mysql|sqlite|mssql)://(.+?):(.+?)@(.+?):([0-9]+)/(.+)$")
	creds := make(map[string]string, 0)

	if reg.MatchString(dbUrl) {
		submatches := reg.FindStringSubmatch(dbUrl)
		fmt.Println("match found", submatches)

		dbDialect = submatches[1]
		creds["user"] = submatches[2]
		creds["pass"] = submatches[3]
		creds["host"] = submatches[4]
		// creds["port"] = submatches[5]
		creds["database"] = submatches[6]
	} else {
		e := godotenv.Load() //Load .env file
		if e != nil {
			fmt.Print(e)
		}

		dbDialect = os.Getenv("db_dialect")
		creds["user"] = os.Getenv("db_user")
		creds["pass"] = os.Getenv("db_pass")
		creds["host"] = os.Getenv("db_host")
		// creds["port"] = os.Getenv("db_port")
		creds["database"] = os.Getenv("db_name")
	}

	dbUri = fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=require",
		creds["host"], creds["user"], creds["database"], creds["pass"]) //Build connection string
	fmt.Println(dbDialect, dbUri)
	return
}

func init() {
	// dbLang, dbUri := getDbInfoFromEnv()
	conn, err := gorm.Open(getDbInfoFromEnv()) //(dbLang, dbUri)
	if err != nil {
		fmt.Print(err)
		panic("failed to connect to database")
	}

	db = conn
	migrate()
	// db.Debug().AutoMigrate(&Account{}, &Playlist{}, &Song{}) //Database migration
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
