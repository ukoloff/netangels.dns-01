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

	t.Run("present+cleanup", func(t *testing.T) {
		err := FireWebTest()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Lego::httpreq", func(t *testing.T) {
		t.SkipNow()
		err := lego("httpreq",
			"HTTPREQ_ENDPOINT=http://localhost",
			"HTTPREQ_POLLING_INTERVAL=10",
			"HTTPREQ_PROPAGATION_TIMEOUT=300",
		)
		if err != nil {
			t.Fatal(err)
		}
	})
}

type txtRR struct {
	fqdn  string
	value string
}

func newTxt(z na01.Zone) (result txtRR) {
	result.fqdn = strings.ToLower("acme-" + na01.RandomString(5) + "." + z.Name)
	result.value = na01.RandomString(12)
	return
}

func (txt txtRR) match(rr na01.RR) bool {
	return rr.Type == "TXT" && txt.fqdn == rr.Name && txt.value == rr.Data["value"]
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

	x := newTxt(z)

	rr, err := na01.FirePresent(x.fqdn, x.value)
	if err != nil {
		return err
	}
	if !x.match(rr) {
		return errors.New("got incorrect RR")
	}

	rrs, err := na01.FindRRs(x.fqdn)
	if err != nil {
		return err
	}
	if len(rrs) != 1 {
		return errors.New("new RR not found")
	}
	if !x.match(rrs[0]) {
		return errors.New("incorrect RR found")
	}

	rrs, err = na01.FireCleanUp(x.fqdn, x.value)
	if err != nil {
		return err
	}
	if len(rrs) != 1 || !x.match(rrs[0]) {
		return errors.New("incorrect RR deleted")
	}

	rrs, err = na01.FindRRs(x.fqdn)
	if err != nil {
		return err
	}
	if len(rrs) != 0 {
		return errors.New("new RR has not been deleted")
	}
	return nil
}
