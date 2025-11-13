package na01_test

import (
	"na01"
	"testing"
	"time"
)

func TestWWW(t *testing.T) {
	go func() {
		if err := na01.Start(); err != nil {
			panic(err)
		}
	}()
	defer na01.Stop()

	time.Sleep(108 * time.Millisecond)

	err := na01.FireAlive()
	if err != nil {
		t.Fatal(err)
	}

	z, err := na01.NewZone("http-" + na01.RandomString(7) + ".ru")
	if err != nil {
		t.Fatal(err)
	}
	defer na01.DropZone(z.ID)

	fqdn := "acme-" + na01.RandomString(5) + "." + z.Name
	text := na01.RandomString(12)

	err = na01.FirePresent(fqdn, text)
	if err != nil {
		t.Fatal(err)
	}

	rrs, err := na01.FindRRs(fqdn)
	if err != nil {
		t.Fatal(err)
	}
	if len(rrs) != 1 {
		t.Fatal("new RR not found")
	}
	if rrs[0].Name != fqdn || rrs[0].Data["value"] != text {
		t.Fatal("incorrect RR found")
	}

	err = na01.FireCleanUp(fqdn, text)
	if err != nil {
		t.Fatal(err)
	}

	rrs, err = na01.FindRRs(fqdn)
	if err != nil {
		t.Fatal(err)
	}
	if len(rrs) != 0 {
		t.Fatal("new RR has not been deleted")
	}
}
