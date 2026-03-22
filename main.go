package main

import (
	"fmt"
	"os"

	"github.com/lukas-zx/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	_, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("invalid command: %s", cmd.name)
	}

	c.cmds[cmd.name](s, cmd)
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func main() {
	s := state{
		cfg: config.Read(),
	}

	commands := commands{
		cmds: map[string]func(*state, command) error{},
	} 
	commands.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Printf("too few arguments, expected at least 2, got %d\n", len(args))
		os.Exit(1)
	}

	command := command{
		name: args[0],
		args: args[1:],
	}
	commands.run(&s, command)
}
