package na01_test

import (
	"errors"
	"na01"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCLI(t *testing.T) {
	me, err := buildCLI()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(me)

	t.Run("exec", func(t *testing.T) {
		err := execTest(me)
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
		defer na01.FireQuit()
		time.Sleep(108 * time.Millisecond)

		FireWebTest()
	})

	t.Run("Lego::exec", func(t *testing.T) {
		t.SkipNow()
		err := lego("exec", []string{
			"EXEC_PATH=" + me,
			"EXEC_POLLING_INTERVAL=10",
			"EXEC_PROPAGATION_TIMEOUT=300",
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

func runMe(me string, args ...string) error {
	cmd := exec.Command(me, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func execTest(me string) error {
	z, err := na01.NewZone("exec-" + na01.RandomString(7) + ".net")
	if err != nil {
		return err
	}
	defer na01.DropZone(z.ID)

	fqdn := strings.ToLower("cli-" + na01.RandomString(5) + "." + z.Name)
	text := na01.RandomString(12)

	err = runMe(me, "present", fqdn, text)
	if err != nil {
		return err
	}

	rrs, err := na01.FindRRs(fqdn)
	if err != nil {
		return err
	}
	if len(rrs) != 1 {
		return errors.New("new RR not found")
	}
	if rrs[0].Name != fqdn || rrs[0].Data["value"] != text {
		return errors.New("incorrect RR found")
	}

	err = runMe(me, "cleanup", fqdn, text)
	if err != nil {
		return err
	}

	rrs, err = na01.FindRRs(fqdn)
	if err != nil {
		return err
	}
	if len(rrs) != 0 {
		return errors.New("new RR has not been deleted")
	}
	return nil
}
