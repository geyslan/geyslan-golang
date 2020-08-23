package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Println(`Usage:
DBHOST=host DBUSER=user DBPASS=pass DBDATABASE=db oowlish logfile`)
}

func main() {

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	go startService()
	watchLog(os.Args[1], os.Stdout)
}
