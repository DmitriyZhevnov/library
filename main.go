package main

import (
	"github.com/DmitriyZhevnov/library/src/app"
)

func main() {
	a := app.App{}
	a.Initialize("postgres", "root", "root", "5432", "fullstack-postgres", "fullstack_api")
	a.Run(":8080")
}
