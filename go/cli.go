package na01

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Cli() {
	args := os.Args
	if len(args) < 2 {
		help()
	}

	var action func(string, string) error

	switch args[1] {
	case "present":
		action = Present
	case "cleanup":
		action = CleanUp
	case "www":
		www()
	default:
		help()
	}
	if action == nil {
		return
	}
	if len(args) != 4 {
		help()
	}
	err := action(args[2], args[3])
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}

func help() {
	me, _ := os.Executable()
	me = filepath.Base(me)
	fmt.Printf(`Usage:
  %[1]v present fqdn text - Create TXT DNS record
  %[1]v cleanup fqdn text - Remove TXT DNS record
  %[1]v www               - Start Web server
`, me)
	os.Exit(1)
}

func www() {
	http.HandleFunc("/alive", alive)
	http.HandleFunc("/quit", quit)
	http.HandleFunc("/present", handler)
	http.HandleFunc("/cleanup", handler)
	http.HandleFunc("/", fallback)
	err := http.ListenAndServe("", nil)
	if err != nil {
		panic(err)
	}
}

func alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func quit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bye!")
	c := time.After(1 * time.Second)
	go func() {
		<-c
		os.Exit(0)
	}()
}

func fallback(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/ukoloff/netangels.dns-01", http.StatusMovedPermanently)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.URL)
}
