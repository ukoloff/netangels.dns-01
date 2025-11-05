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
	switch args[1] {
	case "present":
	case "cleanup":
	case "www":
	default:
		help()
	}
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
