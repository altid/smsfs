package main

import (
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/alexdavid/sigma"
	"github.com/altid/libs/fs"
	"github.com/altid/libs/markup"
)

type iMessage struct {
	cl   sigma.Client
	ctrl *fs.Control
}

func getRunner() *iMessage {
	return &iMessage{}
}

func (i *iMessage) Handle(path string, c *markup.Lexer) error {
	id, err := strconv.Atoi(path)
	if err != nil {
		return err
	}

	return i.cl.SendMessage(id, c.String())
}

func (i *iMessage) run(*fs.Control, *fs.Command) error {
	return nil
}

func (i *iMessage) quit() {
	i.cl.Close()
}

func (i *iMessage) listen() error {
	// make a map of lastID
	// update that map as we loop
	return nil
}

func (i *iMessage) setup(ctrl *fs.Control, user string) error {
	cl, err := sigma.NewClient()
	if err != nil {
		log.Println("It's likely that you need to allow unlimited file access to your terminal appliations in privacy settings")
		return nil
	}

	i.cl = cl

	chats, err := cl.Chats()
	if err != nil {
		return err
	}

	ew, err := ctrl.ErrorWriter()
	if err != nil {
		return err
	}

	defer ew.Close()

	// Populate chats
	for _, chat := range chats {
		go func(chat sigma.Chat) {
			msgs, err := cl.Messages(chat.ID, sigma.MessageFilter{Limit: 500})
			if err != nil {
				fmt.Fprintf(ew, "%v\n", err)
				return
			}

			if len(msgs) < 1 {
				return
			}

			if e := buildMessage(chat, ctrl, msgs, user); e != nil {
				fmt.Fprintf(ew, "%v\n", e)
				return
			}

			input, err := fs.NewInput(i, path.Join(*mtpt, *srv), chat.DisplayName, *debug)
			if err != nil {
				fmt.Fprintf(ew, "%v\n", err)
				return
			}

			input.Start()
		}(chat)
	}

	return nil
}

func buildMessage(chat sigma.Chat, ctrl *fs.Control, msgs []sigma.Message, user string) error {
	ctrl.CreateBuffer(chat.DisplayName, "feed")
	mw, err := ctrl.MainWriter(chat.DisplayName, "feed")
	if err != nil {
		return err
	}

	defer mw.Close()

	for i := len(msgs); i > 0; i-- {
		if msgs[i-1].FromMe {
			fmt.Fprintf(mw, "%%[%s](grey) %s\n", user, msgs[i-1].Text)
		} else {
			fmt.Fprintf(mw, "%%[%s](blue) %s\n", chat.DisplayName, msgs[i-1].Text)
		}
	}

	return nil
}
