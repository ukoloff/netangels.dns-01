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
		t.Errorf("Cannot list zones: %v", err)
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

func TestNewRR(t *testing.T) {
	tests := []struct {
		name string
		data any
	}{
		{name: "A", data: &na01.RRa{IP: "3.2.1.0"}},
		{name: "TXT", data: &na01.RRtxt{Value: "Preved, medved!"}},
		{name: "CNAME", data: &na01.RRcname{Domain: "ya.ru"}},
	}

	z, err := na01.NewZone("test-" + na01.RandomString(7) + ".com")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		na01.DropZone(z.ID)
	}()

	getCount := func() int {
		rrs, err := na01.ZoneRRs(z.ID)
		if err != nil {
			t.Fatal(err)
		}
		return len(rrs)
	}
	count := getCount()
	checkCount := func(delta int) {
		if count+delta == getCount() {
			return
		}
		t.Fatal("Failed to change records count")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := reflect.ValueOf(test.data).Elem()
			v.FieldByName("Type").SetString(test.name)
			newName := "rr-" + na01.RandomString(5) + "." + z.Name
			v.FieldByName("Name").SetString(newName)

			res, err := na01.NewRR(test.data)
			if err != nil {
				t.Fatal(err)
			}

			checkCount(+1)

			if res.Type != test.name {
				t.Fatal("Invalid record type created")
			}

			found, err := na01.FindRRs(newName)
			if err != nil {
				t.Fatal(err)
			}
			if len(found) != 1 {
				t.Fatal("Failed to find new RR")
			}
			if found[0].Name != res.Name || found[0].Type != test.name {
				t.Fatal("Invalid RR found")
			}
			res, err = na01.DropRR(res.ID)
			if err != nil {
				t.Fatal(err)
			}
			checkCount(0)
		})
	}
}
