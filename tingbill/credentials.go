package main

import (
	"code.google.com/p/gopass"
	"fmt"
	"log"
)

func Username() string {
	if *username != "" {
		return *username
	} else {
		var usr string
		fmt.Printf("Username: ")
		fmt.Scanf("%s", &usr)
		if usr == "" {
			log.Fatal("Cannot proceed without username.")
		}
		return usr
	}
}

func Password() string {
	if *password != "" {
		return *password
	} else {
		pwd, err := gopass.GetPass("Password: ")
		if err != nil {
			log.Fatal(err)
		}
		if pwd == "" {
			log.Fatal("Cannot proceed without password.")
		}
		return pwd
	}
}
