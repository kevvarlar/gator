package main

import (
	"github.com/kevvarlar/gator/internal/config"
	"github.com/kevvarlar/gator/internal/database"
)

type state struct{
	gatorConfig *config.Config
	gatorDB *database.Queries
}