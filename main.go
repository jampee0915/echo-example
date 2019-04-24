package main

import (
	"database/sql"
	"gopkg.in/gorp.v2"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

var dbDriver = "mysql"

type Comment struct {
	Id      int64     `json:"id" db:"id,primarykey,autoincrement"`
	Name    string    `json:"name" db:"name,notnull,default:'',size:200"`
	Text    string    `json:"text" db:"text,notnull,size:399"`
	Updated time.Time `json:"updated" db:"updated,notnull"`
	Created time.Time `json:"created" db:"created,notnull"`
}

type Controller struct {
	dbmap *gorp.DbMap
}

func setupDB() (*gorp.DbMap, error) {
	db, err := sql.Open(dbDriver, "root@tcp(127.0.0.1:3306)/example")
	if err != nil {
		return nil, err
	}

	var diarect gorp.Dialect = gorp.MySQLDialect{"InnoDB", "UTF8"}

	// for testing
	if dbDriver == "sqlite3" {
		diarect = gorp.SqliteDialect{}
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: diarect}
	dbmap.AddTableWithName(Comment{}, "comments").SetKeys(true, "id")
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	return dbmap, nil
}

func (controller *Controller) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Debug = true
	e.Logger.SetOutput(os.Stderr)

	e.Use(middleware.Logger())

	return e
}

func main() {
	dbmap, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}

	controller := &Controller{dbmap: dbmap}

	e := setupEcho()

	e.GET("/api/health", controller.Health)
	e.Logger.Fatal(e.Start(":8989"))
}
