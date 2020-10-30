package main

import (
	"fmt"

	"github.com/rushteam/sqlcli/suggest"
)

var (
	//Version ..
	Version = "0.0.0"
	//DbDefConfFilePath ..
	DbDefConfFilePath = []string{
		"/etc/my.cnf",
		"/etc/mysql/my.cnf",
		"/usr/local/etc/my.cnf",
		"~/.my.cnf",
	}
)

//Connected ..

func main() {
	fmt.Printf("sqlcli %s\n", Version)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer fmt.Println("Bye!")
	eng := &suggest.DbEngine{
		Type:   "mysql",
		DbName: "",
		User:   "user",
		Pass:   "user",
		Host:   "192.168.1.1",
		Port:   3306,
	}
	state := &suggest.StateCommon{
		Engine: eng,
	}
	state.Run()
}
