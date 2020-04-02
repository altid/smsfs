package main

import (
	"errors"

	"github.com/altid/libs/fs"
	"github.com/altid/libs/markup"
)

type runner interface {
	setup(*fs.Control, string) error
	listen() error

	run(*fs.Control, *fs.Command) error
	handle(string, *markup.Lexer) error
	quit()
}
type server struct {
	cmd runner
}

func (s *server) Run(c *fs.Control, cmd *fs.Command) error {
	switch cmd.Name {
	case "block", "open", "close":
		return s.cmd.run(c, cmd)
	default:
		return errors.New("command not supported")
	}
}

func (s *server) Handle(path string, c *markup.Lexer) error {
	return s.cmd.handle(path, c)
}

func (s *server) Quit() {
	s.cmd.quit()
}

func (s *server) setup(ctrl *fs.Control, user string) error {
	return s.cmd.setup(ctrl, user)
}

func (s *server) listen() error {
	return s.cmd.listen()
}
