package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
	"sync"
)

var DB *gorm.DB
var DbMutex sync.Mutex

func init() {
	e := godotenv.Load()

	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_NAME")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)

	conn, err := gorm.Open("postgres", dbUri)

	if err != nil {
		fmt.Print(err)
		panic(err)
	}

	fmt.Println("Database connection established successfully")

	DB = conn
	DB.Debug().AutoMigrate(&Tag{})
	DB.Debug().AutoMigrate(&Post{})
	DB.Debug().AutoMigrate(&Project{})
	DB.Debug().AutoMigrate(&User{})
	DB.Debug().AutoMigrate(&Team{})
	DB.Debug().AutoMigrate(&Invitation{})

	fmt.Println("Database migration completed")
}

func GetDB() *gorm.DB {
	return DB
}
