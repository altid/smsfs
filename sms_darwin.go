package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alexdavid/sigma"
	"github.com/altid/libs/fs"
	"github.com/altid/libs/markup"
)

type iMessage struct {
	cl   sigma.Client
	ctrl *fs.Control
	user string
	last map[int]int
	sync.Mutex
}

func getRunner() *iMessage {
	return &iMessage{
		last: make(map[int]int),
	}
}

func (i *iMessage) handle(path string, c *markup.Lexer) error {
	msg, err := c.String()
	if err != nil {
		return err
	}

	chats, err := i.cl.Chats()
	if err != nil {
		return err
	}

	for _, chat := range chats {
		if chat.DisplayName == path {
			return i.cl.SendMessage(chat.ID, msg)
		}
	}

	return errors.New("unable to find client")
}

func (i *iMessage) run(*fs.Control, *fs.Command) error {
	return nil
}

func (i *iMessage) quit() {
	i.cl.Close()
}

func (i *iMessage) listen() error {
	ew, err := i.ctrl.ErrorWriter()
	if err != nil {
		return err
	}

	for {
		chats, err := i.cl.Chats()
		if err != nil {
			return err
		}

		time.Sleep(time.Second * 3)

		for _, chat := range chats {
			go func(i *iMessage, chat sigma.Chat) {
				msgs, err := i.cl.Messages(chat.ID, sigma.MessageFilter{
					AfterID: i.last[chat.ID],
				})

				if err != nil {
					fmt.Fprintf(ew, "%v\n", err)
					return
				}

				if len(msgs) < 1 {
					return
				}

				if !i.ctrl.HasBuffer(chat.DisplayName, "feed") {
					i.ctrl.CreateBuffer(chat.DisplayName, "feed")
					i.ctrl.Input(chat.DisplayName)
				}

				if e := i.buildMessage(chat, msgs); e != nil {
					fmt.Fprintf(ew, "%v\n", e)
					return
				}

			}(i, chat)
		}
	}
}

func (i *iMessage) setup(ctrl *fs.Control, user string) error {
	i.ctrl = ctrl
	i.user = user
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

			ctrl.CreateBuffer(chat.DisplayName, "feed")

			if e := i.buildMessage(chat, msgs); e != nil {
				fmt.Fprintf(ew, "%v\n", e)
				return
			}

			i.ctrl.Input(chat.DisplayName)
		}(chat)
	}

	return nil
}

func (i *iMessage) buildMessage(chat sigma.Chat, msgs []sigma.Message) error {
	mw, err := i.ctrl.MainWriter(chat.DisplayName, "feed")
	if err != nil {
		return err
	}

	defer mw.Close()

	for n := len(msgs); n > 0; n-- {
		i.Lock()
		i.last[chat.ID] = msgs[n-1].ID
		i.Unlock()

		if msgs[n-1].FromMe {
			fmt.Fprintf(mw, "%%[%s](grey) %s\n", i.user, msgs[n-1].Text)
		} else {
			fmt.Fprintf(mw, "%%[%s](blue) %s\n", chat.DisplayName, msgs[n-1].Text)
		}
	}

	return nil
}
