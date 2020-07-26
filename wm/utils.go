package wm

import (
	"log"
	"os/user"
)

func homeDir() string {

	usr, err := user.Current()
	if err != nil {
		log.Panic(err)
	}

	return usr.HomeDir
}
