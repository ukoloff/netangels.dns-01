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
	me, err := buildCLI()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(me)

	t.Run("www", func(t *testing.T) {
		cmd := exec.Command(me, "www")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			t.Fatal(err)
		}
		defer na01.FireQuit()
		time.Sleep(108 * time.Millisecond)

		FireWebTest()
	})

	t.Run("Lego/exec", func(t *testing.T) {
		t.SkipNow()
		err := lego("exec", map[string]string{
			"EXEC_PATH":                me,
			"EXEC_POLLING_INTERVAL":    "10",
			"EXEC_PROPAGATION_TIMEOUT": "300",
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func exePath() string {
	me, _ := os.Executable()
	f, _ := os.CreateTemp("", "na01-*"+filepath.Ext(me))
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

func buildCLI() (string, error) {
	me := exePath()
	cmd := exec.Command("go", "build", "-C", "main", "-ldflags", "-s -w", "-o", me)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return me, cmd.Run()
}
