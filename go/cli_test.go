package na01_test

import (
	"na01"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestCLI(t *testing.T) {
	me := exePath()
	cmd := exec.Command("go", "build", "-C", "main", "-ldflags", "-s -w", "-o", me)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(me)

	t.Run("cli", func(t *testing.T) {
		// t.SkipNow()
		err := lego("exec", map[string]string{
			"EXEC_PATH":                me,
			"EXEC_POLLING_INTERVAL":    "10",
			"EXEC_PROPAGATION_TIMEOUT": "300",
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("www", func(t *testing.T) {
		cmd := exec.Command(me, "www")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(108 * time.Millisecond)
		defer na01.FireQuit()
	})
}

func exePath() string {
	me, _ := os.Executable()
	f, _ := os.CreateTemp("", "na01-*"+filepath.Ext(me))
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}
