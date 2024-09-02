package main

import (
	"log"

	"github.com/access6080/jesorm"
)

func main() {
	type User struct {
		Name       string `type:"VARCHAR(255)" constraints:"PrimaryKey"`
		Id         int    `type:"Int(2)"`
		Registered bool   `constraints:"ForeignKey(Messages.id)"`
		Account    int    `type:"INTEGER" constraints:"ForeignKey(Account.id), NotNull, Unique"`
	}

	config := jesorm.Config{
		DriverName: "sqlite3",
		DSN:        "test:test.db",
	}

	//Initialize Orm and database
	model, err := jesorm.Init(config)
	if err != nil {
		log.Fatal(err)
	}

	// Register your models with the orm
	if err = model.Add(User{}); err != nil {
		log.Fatal(err)
	}

	// Watch models for changes
	if err = jesorm.AutoMigrate(*model); err != nil {
		panic(err)
	}
}
