package geoLocation

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	geoip2 "github.com/oschwald/geoip2-golang"
)

type Reader struct {
	*geoip2.Reader
	sync.RWMutex
	updateDuration time.Duration
	path           string
}

func Open(path string, updateDuration time.Duration, ctx context.Context) (*Reader, error) {
	db, err := geoip2.Open(path)
	if err != nil {
		return nil, err
	}
	reader := new(Reader)
	reader.Reader = db
	reader.updateDuration = updateDuration
	reader.path = path
	if reader.updateDuration > time.Duration(0) {
		go reader.update(ctx)
	}
	return reader, nil
}

func (reader *Reader) update(ctx context.Context) {
	for {
		select {
		case <-time.After(reader.updateDuration):
		case <-ctx.Done():
			return
		}
		db, err := geoip2.Open(reader.path)
		if err != nil {
			fmt.Println("geoip2 read failed: ", err)
		} else {
			reader.Lock()
			reader.Reader = db
			reader.Unlock()
		}
	}
}

func (reader *Reader) CountryIsoCode(ip net.IP) (string, bool) {
	reader.RLock()
	defer reader.RUnlock()
	location, err := reader.Country(ip)
	if err != nil {
		return "", false
	}
	return location.Country.IsoCode, true
}

func (reader *Reader) CityName(ip net.IP) (string, bool) {
	reader.RLock()
	defer reader.RUnlock()
	location, err := reader.City(ip)
	if err != nil {
		return "", false
	}
	if city, found := location.City.Names["en"]; found {
		return city, true
	}
	return "", false
}
