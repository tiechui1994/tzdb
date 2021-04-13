# tzdb (Time Zone Database)

tzdb is a library to determine whether the time is in daylight saving time.

### IANA Time Zone Database

The data returned and used by TZInfo is sourced from the [IANA Time Zone Database](https://www.iana.org/time-zones). 
The [Theory and pragmatics of the tz code and data](https://data.iana.org/time-zones/theory.html) document gives 
details of how the data is organized and managed.


### Installation

```bash
go get -u github.com/tiechui1994/tzdb
```

### Example Usage 

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/tiechui1994/tzdb"
)

func main() {
    loc, _ := time.LoadLocation("America/Chicago")
    timeFormat := "2006-01-02 15:04:05"
    
    testz, _ := time.ParseInLocation(timeFormat, "2021-03-14 01:59:00", loc)
    fmt.Println(testz, testz.UTC(), tzdb.IsDST(testz))
    
    testz = testz.Add(time.Minute)
    fmt.Println(testz, testz.UTC(), tzdb.IsDST(testz))
    
    testz = testz.Add(time.Minute)
    fmt.Println(testz, testz.UTC(), tzdb.IsDST(testz))
}
```

### Update location zoneinfo

```bash
curl https://raw.githubusercontent.com/tiechui1994/tzdb/main/scripts/update.sh | sudo sh -
```