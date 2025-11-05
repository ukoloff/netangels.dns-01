package na01

import (
	"fmt"
	"os"
	"path/filepath"
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
		println("WWW...")
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
