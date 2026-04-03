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
	if err != nil {
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

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		fmt.Printf("error deleting users: %v", err)
		os.Exit(1)
	}

	fmt.Printf("successfully deleted users")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error getting users: %v", err)
		os.Exit(1)
	}

	for _, user := range users {
		userName := user.Name.String
		if userName == s.cfg.CurrentUserName {
			userName += " (current)"
		}
		fmt.Printf("* %s\n", userName)
	}

	return nil
}

func handlerAggregate(s *state, cmd command) error {
	//feed, err := fetchFeed(context.Background(), cmd.args[0])
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Printf("error fetching feed: %v", err)
		os.Exit(1)
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		fmt.Printf("too few arguments, expected at least 2, got %d\n", len(cmd.args))
		os.Exit(1)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      sql.NullString{String: cmd.args[0], Valid: true},
		Url:       sql.NullString{String: cmd.args[1], Valid: true},
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		fmt.Printf("error creating feed: %v", err)
		os.Exit(1)
	}

	if _, err = s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feed.ID, Valid: true},
	}); err != nil {
		fmt.Printf("error linking feed to user: %v", err)
		os.Exit(1)
	}

	fmt.Println(feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Printf("error getting feeds: %v", err)
		os.Exit(1)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID.UUID)
		if err != nil {
			fmt.Printf("error getting user for feed: %v", err)
			os.Exit(1)
		}
		fmt.Printf("name: %s\nurl: %s\ncreated by: %s\n\n", feed.Name.String, feed.Url.String, user.Name.String)
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeedByURL(context.Background(), sql.NullString{String: cmd.args[0], Valid: true})
	if err != nil {
		fmt.Printf("error getting feed from url: %v", err)
		os.Exit(1)
	}

	feedFollows, err := s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feed.ID, Valid: true},
	})
	if err != nil {
		fmt.Printf("error creating feed follows: %v", err)
		os.Exit(1)
	}

	fmt.Printf("feed %s is now followed by %s\n", feedFollows.FeedName.String, feedFollows.UserName.String)
	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeedByURL(context.Background(), sql.NullString{String: cmd.args[0], Valid: true})
	if err != nil {
		fmt.Printf("error getting feed from url: %v", err)
		os.Exit(1)
	}

	if err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feed.ID, Valid: true},
	}); err != nil {
		fmt.Printf("error unfollowing feed: %v", err)
		os.Exit(1)
	}

	fmt.Printf("feed %s is no longer followed by %s\n", feed.Name.String, user.Name.String)
	return nil
}

func handlerGetFollowingForUser(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), uuid.NullUUID{UUID: user.ID, Valid: true})
	if err != nil {
		fmt.Printf("error getting followed feeds: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%s is following:\n", user.Name.String)
	for _, feedFollow := range feedFollows {
		fmt.Println(feedFollow.FeedName.String)
	}
	return nil
}
