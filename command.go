package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevvarlar/gator/internal/database"
)

type command struct {
	name string
	arguments []string
}

type commands struct {
	commandNames map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commandNames[cmd.name]
	if !ok {
		return fmt.Errorf("command does not exist")
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandNames[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	_, err := s.gatorDB.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("you cannot login to an account that does not exist: %w", err)
	}
	err = s.gatorConfig.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("The user has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the register handler expects a single argument, the username")
	}
	_, err := s.gatorDB.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		Name: cmd.arguments[0],
	})
	if err != nil {
		return err
	}
	err = handlerLogin(s, cmd)
	if err != nil {
		return err
	}
	fmt.Println("user was created")
	return nil
}

func handlerReset(s *state, c command) error {
	err := s.gatorDB.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users database: %w", err)
	}
	fmt.Println("Successfully reseted users database")
	return nil
}