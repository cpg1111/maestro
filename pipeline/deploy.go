/*
Copyright 2016 Christian Grabowski All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
