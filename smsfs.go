package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/user"

	"github.com/altid/libs/config"
	"github.com/altid/libs/config/types"
	"github.com/altid/libs/fs"
)

var (
	mtpt  = flag.String("p", "/tmp/altid", "Path for filesystem")
	srv   = flag.String("s", "sms", "Name of service")
	debug = flag.Bool("d", false, "enable debug logging")
	setup = flag.Bool("conf", false, "Run configuration setup")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	u, _ := user.Current()
	conf := &struct {
		User          string
		ListenAddress types.ListenAddress
		Logdir        types.Logdir
	}{u.Name, "none", "none"}

	if *setup {
		if e := config.Create(conf, *srv, "", *debug); e != nil {
			log.Fatal(e)
		}

		os.Exit(0)
	}

	if e := config.Marshal(conf, *srv, "", *debug); e != nil {
		log.Fatal(e)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := &server{cancel, getRunner()}

	ctrl, err := fs.CreateCtlFile(ctx, s, string(conf.Logdir), *mtpt, *srv, "feed", *debug)
	if err != nil {
		log.Fatal(err)
	}

	defer ctrl.Cleanup()
	ctrl.SetCommands(Commands...)
	ctrl.CreateBuffer("server", "document")

	if e := s.setup(ctrl, conf.User); e != nil {
		log.Fatal(e)
	}

	go s.listen()

	if e := ctrl.Listen(); e != nil {
		log.Fatal(e)
	}
}
