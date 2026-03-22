package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lukas-zx/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing username")
	}

	userName := sql.NullString{
		String: cmd.args[0],
		Valid:  true,
	}
	_, err := s.db.GetUser(context.Background(), userName)
	if (err != nil) {
		fmt.Printf("user does not exist: %s\n", cmd.args[0])
		os.Exit(1)
	}

	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("username set to %s", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing username")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      sql.NullString{String: cmd.args[0], Valid: true},
	})
	if err != nil {
		fmt.Printf("error creating user: %v", err)
		os.Exit(1)
	}

	s.cfg.CurrentUserName = user.Name.String
	s.cfg.SetUser(user.Name.String)

	fmt.Printf("user created: %v", user)
	return nil
}
