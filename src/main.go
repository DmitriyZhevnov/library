package main

import (
	"github.com/DmitriyZhevnov/library/src/app"
	_ "github.com/DmitriyZhevnov/library/src/app"
)

func main() {
	a := app.App{}
	a.Initialize("root", "950621", "library")
	a.Run(":8080")
}
