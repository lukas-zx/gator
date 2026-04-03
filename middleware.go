package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/lukas-zx/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(
			context.Background(),
			sql.NullString{String: s.cfg.CurrentUserName, Valid: true},
		)
		if err != nil {
			fmt.Printf("error getting user: %v\n", err)
			os.Exit(1)
		}

		return handler(s, cmd, user)
	}

}
