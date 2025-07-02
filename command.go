package main

import "fmt"

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
	err := s.gatorConfig.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("The user has been set")
	return nil
}