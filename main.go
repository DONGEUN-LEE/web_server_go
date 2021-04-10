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
	data "web_server.com/m/v1/data"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB //database

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
	db.Debug().AutoMigrate(&data.User{}, &data.Plan{}) //Database migration
}

func initWeb() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("/again", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World Again!")
	})
	e.GET("/api/plan", func(c echo.Context) error {
		var plans []data.Plan
		db.Find(&plans)
		doc, _ := json.Marshal(plans)
		return c.String(http.StatusOK, string(doc))
	})
	e.POST("/api/login", func(c echo.Context) error {
		json_map := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&json_map)
		if err != nil {
				return err
		} else {
			var user data.User;
			email := json_map["email"].(string)
			password := []byte(json_map["password"].(string))
			db.Where("email = ?", email).First(&user)
			
			// hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
			// if err != nil {
			// 		return err
			// }

			loginRes := data.LoginRes{}
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), password)
			if err == nil {
				token, _ := createToken(email)
				loginRes.Token = token
				loginRes.Message = "Success"
			} else {
				loginRes.Message = "Fail"
			}

			doc, _ := json.Marshal(loginRes)
			return c.String(http.StatusOK, string(doc))
		}
	})
	e.Logger.Fatal(e.Start(":1213"))
}

func createToken(userid string) (string, error) {
  var err error
  //Creating Access Token
  atClaims := jwt.MapClaims{}
  atClaims["authorized"] = true
  atClaims["user_id"] = userid
  at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
  token, err := at.SignedString([]byte(os.Getenv("SECRET_KEY")))
  if err != nil {
     return "", err
  }
  return token, nil
}

func main() {
	initDb()
	initWeb()
}
