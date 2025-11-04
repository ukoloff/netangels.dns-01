package na01_test

import (
	"na01"
	"reflect"
	"testing"
)

func TestAuth(t *testing.T) {
	token, err := na01.Auth()
	if err != nil {
		t.Error(err)
	}
	if len(token) == 0 {
		t.Error("Empty token")
	}
}

func TestZones(t *testing.T) {
	zones, err := na01.Zones()
	if err != nil {
		panic(err)
	}
	for _, z := range zones {
		t.Run("Fetch zone: "+z.Name, func(t *testing.T) {
			zone, err := na01.GetZone(z.ID)
			if err != nil {
				t.Errorf("Fetch failed for %v: %v", z.Name, err)
			}
			if !reflect.DeepEqual(zone, z) {
				t.Errorf("Zone mismatch: %v != %v", zone, z)
			}
			rrs, err := na01.ZoneRRs(z.ID)
			if err != nil {
				t.Fatal(err)
			}
			if zone.Count != len(rrs) {
				t.Errorf("Count mismatch: %v != %v", zone.Count, len(rrs))
			}
		})
	}
}

func TestNewZone(t *testing.T) {
	zname := "test-" + na01.RandomString(7) + ".ru"

	read := func() int {
		zones, err := na01.Zones()
		if err != nil {
			t.Errorf("Reading failed: %v", err)
		}
		return len(zones)
	}
	count := read()

	z, err := na01.NewZone(zname)

	if err != nil {
		t.Fatal("Creation failed", err)
	}
	defer func() {
		if z.ID == 0 {
			return
		}
		na01.DropZone(z.ID)
	}()

	if read() != count+1 {
		t.Fatal("Cannot increment Zone count")
	}

	z2, err := na01.DropZone(z.ID)
	if err != nil {
		t.Fatal("Destroy failed", err)
	}

	z.ID = 0

	if !reflect.DeepEqual(z, z2) {
		t.Errorf("Zone mismatch: %v != %v", z, z2)
	}
	if read() != count {
		t.Fatal("Cannot decrement Zone count")
	}
}
