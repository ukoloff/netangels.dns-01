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
		})
	}
}
