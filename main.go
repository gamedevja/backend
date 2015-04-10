package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gamedevja/backend/git"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/testpush", testpush)
	fmt.Println("listening...")
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, heroku via github")
}

func testpush(w http.ResponseWriter, r *http.Request) {
	t := fmt.Sprintf("%d", time.Now().UnixNano())
	if err := git.Push("assets/testimage/"+t, t, "remote"); err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, "done")
	}
}
