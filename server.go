package main

import (
	"context"
	"errors"

	"github.com/altid/libs/fs"
)

type runner interface {
	open(*fs.Control, string) error
	close(*fs.Control, string) error

	block(*fs.Control, *fs.Command) error
	setup(*fs.Control, string) error
	listen() error

	restart(*fs.Control) error
	refresh(*fs.Control) error
	quit()
}
type server struct {
	cancel context.CancelFunc
	run    runner
}

func (s *server) Open(c *fs.Control, msg string) error {
	return s.run.open(c, msg)
}

func (s *server) Close(c *fs.Control, msg string) error {
	return s.run.close(c, msg)
}

func (s *server) Link(c *fs.Control, from, msg string) error {
	return errors.New("Link command not supported for sms")
}

func (s *server) Default(c *fs.Control, cmd *fs.Command) error {
	switch cmd.Name {
	case "block":
		return s.run.block(c, cmd)
	}
	return nil
}

func (s *server) Restart(c *fs.Control) error {
	return s.run.restart(c)
}

func (s *server) Refresh(c *fs.Control) error {
	return s.run.refresh(c)
}

func (s *server) Quit() {
	s.run.quit()
	s.cancel()
}

func (s *server) setup(ctrl *fs.Control, user string) error {
	return s.run.setup(ctrl, user)
}

func (s *server) listen() error {
	return s.run.listen()
}
