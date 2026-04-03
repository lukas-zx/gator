package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/lukas-zx/gator/internal/config"
	"github.com/lukas-zx/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db *database.Queries
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

	db, err := sql.Open("postgres", s.cfg.DbURL)
	if err != nil {
		fmt.Printf("failed to open database connection: %v", err)
	}
	dbQueries := database.New(db)
	s.db = dbQueries

	commands := commands{
		cmds: map[string]func(*state, command) error{},
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListUsers)
	commands.register("agg", handlerAggregate)
	commands.register("addfeed", handlerAddFeed)
	commands.register("feeds", handlerListFeeds)

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("too few arguments, expected at least 1, got %d\n", len(args))
		os.Exit(1)
	}

	command := command{
		name: args[0],
		args: args[1:],
	}
	commands.run(&s, command)
}
