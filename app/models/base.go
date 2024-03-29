package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
	"regexp"
	"strings"
)

var db *gorm.DB //database

func loadDotEnv() {
	cwd, err := os.Getwd()

	envPath := "./.env"
	if strings.HasSuffix(cwd, "tests") {
		envPath = "../../../.env"
	}
	err = godotenv.Load(envPath)
	if err != nil {
		fmt.Print(err)
	}
}

func migrate() {
	db.AutoMigrate(
		&Account{},
		&Playlist{},
		&Song{},
		&Vote{},
		&Role{},
		&PlaylistAccessControl{})
	db.Model(&Account{}).RemoveIndex("token")
	db.Model(&Song{}).AddForeignKey("playlist_id", "playlists(id)", "CASCADE", "CASCADE")
	db.Model(&Playlist{}).AddForeignKey("user_id", "accounts(id)", "CASCADE", "CASCADE")
	db.Model(&Vote{}).AddForeignKey("user_id", "accounts(id)", "CASCADE", "CASCADE")
	db.Model(&Vote{}).AddForeignKey("song_id", "songs(id)", "CASCADE", "CASCADE")
}

func getDbInfoFromEnv() (dbDialect string, dbUri string) {
	dbUrl := os.Getenv("DATABASE_URL")
	reg := regexp.MustCompile("^(postgres|mysql|sqlite|mssql)://(.+?):(.*?)@(.+?):([0-9]+)/(.+)$")
	creds := make(map[string]string, 0)

	if reg.MatchString(dbUrl) {
		submatches := reg.FindStringSubmatch(dbUrl)

		dbDialect = submatches[1]
		creds["user"] = submatches[2]
		creds["pass"] = submatches[3]
		creds["host"] = submatches[4]
		// creds["port"] = submatches[5]
		creds["database"] = submatches[6]
	} else {
		loadDotEnv()
		dbDialect = os.Getenv("db_dialect")
		creds["user"] = os.Getenv("db_user")
		creds["pass"] = os.Getenv("db_pass")
		creds["host"] = os.Getenv("db_host")
		// creds["port"] = os.Getenv("db_port")
		creds["database"] = os.Getenv("db_name")
	}

	sslmode := "require"
	if creds["host"] == "localhost" || creds["host"] == "127.0.0.1" {
		sslmode = "disable"
	}
	dbUri = fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=%s",
		creds["host"], creds["user"], creds["database"], creds["pass"], sslmode) //Build connection string
	return
}

func init() {
	dbDialect, dbCnx := getDbInfoFromEnv()
	conn, err := gorm.Open(dbDialect, dbCnx)
	if err != nil {
		fmt.Print(err)
		panic("failed to connect to database")
	}

	db = conn
	migrate()
	db.Set("gorm:auto_preload", true)
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
