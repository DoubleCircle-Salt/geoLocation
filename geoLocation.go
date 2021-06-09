package geoLocation

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	geoip2 "github.com/oschwald/geoip2-golang"
)

const (
	geoUpdateDuration = 300 * time.Second
)

var (
	geoDatabase *geoip2.Reader
	geoLock     sync.RWMutex
)

// Geo 结构体
type Geo struct {
	Enabled bool   `json:"enabled"`
	GeoIP2  string `json:"geoip2"`
}

func geoUpdate(ctx context.Context) {
	for {
		select {
		case <-time.After(geoUpdateDuration):
		case <-ctx.Done():
			return
		}
		geoConfig := &GetDynamicConfig().Geo
		db, err := geoip2.Open(geoConfig.GeoIP2)
		if err != nil {
			fmt.Println("geoip2 read failed: ", err)
		} else {
			geoLock.Lock()
			geoDatabase = db
			geoLock.Unlock()
		}
	}
}

func GetGeo() *geoip2.Reader {
	geoLock.RLock()
	defer geoLock.RUnlock()
	return geoDatabase
}

func GetCountryCity(ip net.IP) (country, city string) {
	db := GetGeo()
	if db == nil {
		return
	}
	location, err := db.City(ip)
	if err != nil {
		return
	}
	country = location.Country.IsoCode
	if cityName, found := location.City.Names["en"]; found {
		city = cityName
	}
	return
}

func init() {
	geoConfig := &GetDynamicConfig().Geo
	db, err := geoip2.Open(geoConfig.GeoIP2)
	if err != nil {
		fmt.Println("geoip2 read failed: ", err)
	} else {
		geoDatabase = db
	}
	go geoUpdate(context.Background())
}
