package main

import (
	"fmt"
	"net/http"

	"github.com/mariaulfa33/go-auth/router"
)

func main() {
	var router router.Router
	fmt.Println("Listening on port 3000...")
	http.ListenAndServe(":3000", router)
}
