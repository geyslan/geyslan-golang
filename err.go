package main

import "log"

func eFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
