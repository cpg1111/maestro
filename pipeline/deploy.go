package pipeline

import (
	"log"
)

func check(srv *Service) error {
	check, checkErr := srv.execCheck()
	if checkErr != nil {
		return checkErr
	}
	if check {
		log.Println(check)
	}
	return nil
}
