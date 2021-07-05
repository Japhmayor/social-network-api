package main

import (
	"github.com/japhmayor/social-media-api/api"
)

func main() {
	a := api.App{}
	a.Initialize("root", "bontusfavor1994?", "127.0.0.1", "3306", "social_network")
	a.Run(":3000")
}
