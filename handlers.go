package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing username")
	}

	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("username set to %s", cmd.args[0])

	return nil
}

