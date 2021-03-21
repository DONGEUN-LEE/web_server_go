package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"github.com/joho/godotenv"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/labstack/echo"
)

var db *gorm.DB //database

// Product is ...
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func initDb() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")


	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Product{}) //Database migration
}

func initWeb() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("/again", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World Again!")
	})
	e.GET("/plan", func(c echo.Context) error {
		var prods []Product
		db.Find(&prods)
		doc, _ := json.Marshal(prods);
		return c.String(http.StatusOK, string(doc))
	})
	e.Logger.Fatal(e.Start(":1213"))
}

func main() {
	initDb()
	initWeb()
}
