package main

import (
	"context"
	"net/http"

	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Player struct {
	bun.BaseModel `bun:"table:players,alias:p"`
	Name          string `json:"name"`
	UUID          string `json:"uuid" bun:",pk"`
	Coins         uint32 `json:"coins"`
	Rank          uint8  `json:"rank"`
}

var db *bun.DB

var ctx = context.Background()

func getPlayers(c *gin.Context) {

	var users []Player
	err := db.NewSelect().Model(&users).Where("rank = ?", 1).Scan(ctx)

	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, users)
}

func addPlayer(c *gin.Context) {
	var newPlayer Player

	newPlayer.Coins = 0
	newPlayer.Rank = 0

	if err := c.BindJSON(&newPlayer); err != nil {
		return
	}

	_, err := db.NewInsert().Model(&newPlayer).Exec(ctx)

	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newPlayer)
}

func main() {
	dsn := "postgres://postgres:12345@localhost:5432/test_go?sslmode=disable"

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db = bun.NewDB(sqldb, pgdialect.New())

	db.NewCreateTable().
		Model((*Player)(nil)).
		IfNotExists().
		Exec(ctx)

	router := gin.Default()
	router.GET("/players", getPlayers)
	router.POST("/players", addPlayer)

	router.Run("localhost:8080")
}
