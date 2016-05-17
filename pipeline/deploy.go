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
		log.Println("Found service, updating now")
		err := update(srv)
		if err != nil {
			return err
		}
	} else {
		log.Println("Creating new service")
		err := create(srv)
		if err != nil {
			return err
		}
	}
	return nil
}

func create(srv *Service) error {
	err := srv.execCreate()
	if err != nil {
		return err
	}
	return nil
}

func update(srv *Service) error {
	err := srv.execUpdate()
	if err != nil {
		return err
	}
	return nil
}
