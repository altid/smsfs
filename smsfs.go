package main

import (
	"flag"
	"log"
	"os"
	"os/user"

	"github.com/altid/libs/config"
	"github.com/altid/libs/config/types"
	"github.com/altid/libs/fs"
)

var (
	mtpt    = flag.String("p", "/tmp/altid", "Path for filesystem")
	srv     = flag.String("s", "sms", "Name of service")
	cfgfile = flag.String("c", "", "Directory of configuration file")
	debug   = flag.Bool("d", false, "enable debug logging")
	setup   = flag.Bool("conf", false, "Run configuration setup")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	u, _ := user.Current()
	conf := &struct {
		User   string              `altid:"user,prompt:Your name to show in messages"`
		Listen types.ListenAddress `altid:"listen_address,no_prompt"`
		Logdir types.Logdir        `altid:"logdir,no_prompt"`
	}{u.Name, "none", "none"}

	if *setup {
		if e := config.Create(conf, *srv, *cfgfile, *debug); e != nil {
			log.Fatal(e)
		}

		os.Exit(0)
	}

	if e := config.Marshal(conf, *srv, *cfgfile, *debug); e != nil {
		log.Fatal(e)
	}

	s := &server{getRunner()}

	ctrl, err := fs.New(s, string(conf.Logdir), *mtpt, *srv, "feed", *debug)
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
