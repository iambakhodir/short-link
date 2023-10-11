package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_linkHttpDelivery "github.com/iambakhodir/short-link/link/delivery/http"
	_linkHttpMiddleware "github.com/iambakhodir/short-link/link/delivery/http/middleware"
	_linkRepo "github.com/iambakhodir/short-link/link/repository/mysql"
	"github.com/iambakhodir/short-link/link/usecase"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"time"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Tashkent")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.Ping()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _linkHttpMiddleware.InitMiddleware()
	e.Use(middL.CORS)

	linkRepo := _linkRepo.NewMysqlLinkRepository(dbConn)
	timeOutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	lu := usecase.NewLinkUseCase(linkRepo, timeOutContext)

	linkTagRepo := _linkRepo.NewMysqlLinkTagRepository(dbConn)
	linkTagUcase := usecase.NewLinkTagUseCase(linkTagRepo, timeOutContext)

	tagsRepo := _linkRepo.NewMysqlTagsRepository(dbConn)
	tagsUcase := usecase.NewTagsUseCase(tagsRepo, timeOutContext)

	_linkHttpDelivery.NewLinkHandler(e, lu, tagsUcase, linkTagUcase)

	log.Fatal(e.Start(viper.GetString("server.address"))) //nolint
}
