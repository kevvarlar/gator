package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/kevvarlar/gator/internal/config"
	"github.com/kevvarlar/gator/internal/database"
	_ "github.com/lib/pq"
)

func registerAll(gatorCommands *commands) {
	gatorCommands.register("login", middlewareLoggedIn(handlerLogin))
	gatorCommands.register("register", handlerRegister)
	gatorCommands.register("reset", handlerReset)
	gatorCommands.register("users", handlerUsers)
	gatorCommands.register("agg", handlerAgg)
	gatorCommands.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	gatorCommands.register("feeds", handlerFeeds)
	gatorCommands.register("follow", middlewareLoggedIn(handlerFollow))
	gatorCommands.register("following", handlerFollowing)
	gatorCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	gatorCommands.register("browse", handlerBrowse)
}

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
	registerAll(&gatorCommands)
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