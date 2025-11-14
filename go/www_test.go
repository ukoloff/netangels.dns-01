package na01_test

import (
	"errors"
	"na01"
	"strings"
	"testing"
	"time"
)

func TestWWW(t *testing.T) {
	go func() {
		if err := na01.StartWWW(); err != nil {
			panic(err)
		}
	}()
	defer na01.StopWWW()

	err := FireWebTest()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Lego/httpreq", func(t *testing.T) {
		t.SkipNow()
		err := lego("exec", map[string]string{
			// "EXEC_PATH":                me,
			"EXEC_POLLING_INTERVAL":    "10",
			"EXEC_PROPAGATION_TIMEOUT": "300",
		})
		if err != nil {
			t.Fatal(err)
		}
	})

}

func FireWebTest() error {
	time.Sleep(108 * time.Millisecond)

	err := na01.FireAlive()
	if err != nil {
		return err
	}

	z, err := na01.NewZone("http-" + na01.RandomString(7) + ".ru")
	if err != nil {
		return err
	}
	defer na01.DropZone(z.ID)

	fqdn := strings.ToLower("acme-" + na01.RandomString(5) + "." + z.Name)
	text := na01.RandomString(12)

	rr, err := na01.FirePresent(fqdn, text)
	if err != nil {
		return err
	}
	if rr.Name != fqdn || rr.Data["value"] != text {
		return errors.New("got incorrect RR")
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

	rrs, err = na01.FireCleanUp(fqdn, text)
	if err != nil {
		return err
	}
	if rrs[0].Name != fqdn || rrs[0].Data["value"] != text {
		return errors.New("incorrect RR deleted")
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
