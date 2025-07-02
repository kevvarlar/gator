package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/kevvarlar/gator/internal/config"
	"github.com/kevvarlar/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", gatorConfig.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	gatorState := state{
		gatorConfig: &gatorConfig,
		gatorDB: dbQueries,
	}
	gatorCommands := commands{
		commandNames: make(map[string]func(*state, command) error),
	}
	gatorCommands.register("login", handlerLogin)
	gatorCommands.register("register", handlerRegister)
	gatorCommands.register("reset", handlerReset)
	gatorCommands.register("users", handlerUsers)
	gatorCommands.register("agg", handlerAgg)
	gatorCommands.register("addfeed", handlerAddfeed)
	gatorCommands.register("feeds", handlerFeeds)
	arguments := os.Args
	if len(arguments) < 2 {
		log.Fatal("No command name provided")
	}
	gatorCommand := command{
		name: arguments[1],
		arguments: arguments[2:],
	}
	err = gatorCommands.run(&gatorState, gatorCommand)
	if err != nil {
		log.Fatal(err)
	}
}