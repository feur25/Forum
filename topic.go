package main

import (
	"fmt"
)

func createTopicInDB(name string, id int) error {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `topic` (`name`, `user_id`) VALUES ('%s','%d')", name, id))
	checkError(err)
	defer insert.Close()
	return nil
}
