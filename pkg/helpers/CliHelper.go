package helpers

import (
	"os"

	"github.com/yosa12978/waterbottle/pkg/controllers"
)

func InitCLI() error {
	var err error
	command := os.Args[1]
	if command == "migrate" {
		err = controllers.DoMigrate(os.Args[2:])
	} else if command == "help" {
		err = controllers.Help(os.Args[2:])
	}
	if err != nil {
		return err
	}
	return nil
}
