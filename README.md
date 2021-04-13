# tzdb (Time Zone Database)

- [Time Zone Database](https://www.iana.org/time-zones)

update time zone database:

```
curl https://raw.githubusercontent.com/tiechui1994/tzdb/main/scripts/update.sh | sudo sh -
```

- Usage example

```cgo
package main

import (
    "fmt"
    "time"
    
    "github.com/tiechui1994/tzdb"
)

func main() {
   isdst := tzdb.IsDST(time.Now())
   fmt.Println("isdst", isdst)
}
```