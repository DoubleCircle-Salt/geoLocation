package geoLocation

import (
	"context"
	"net"
	"testing"
)

func TestParse(t * testing.T) {
	reader, err := Open("GeoIP2-City.mmdb", 0, context.Background())
	if err != nil {
		t.Error(err)
	}
	if ip := net.ParseIP("58.221.88.150"); ip != nil {
		country, ok := reader.CountryIsoCode(ip)
		if !ok || country != "CN" {
			t.Error("country read failed")
			return
		}
		city, ok := reader.CityName(ip)
		if !ok || city != "Nantong" {
			t.Error("city read failed")
		}
	} else {
		t.Error("parse ip failed")
	}
}