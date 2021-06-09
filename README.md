# geoLocation
maxmind geo country/city 名称解析


## Installation
```golang
go get github.com/DoubleCircle-Salt/geoLocation
```

## Instruction
```golang
func Open(path string, updateDuration time.Duration, ctx context.Context) (*Reader, error)

```


```golang
package main

import (
    "context"
    "fmt"
    "net"
    "time"

    "github.com/DoubleCircle-Salt/geoLocation"
)

func main(){
    reader, err := geoLocation.Open("GeoIP2-City.mmdb", time.Hour, context.Background())
    if err != nil {
         fmt.Println(err)
         return
    }
    if ip := net.ParseIP("58.221.88.150"); ip != nil {
        country, ok := reader.CountryIsoCode(ip)
        if !ok || country != "CN" {
            fmt.Println("country read failed")
            return
        }
        city, ok := reader.CityName(ip)
        if !ok || city != "Nantong" {
             fmt.Println("city read failed")
        }
    } else {
         fmt.Println("parse ip failed")
    }
}
```

