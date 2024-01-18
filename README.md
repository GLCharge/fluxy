# go-fluxy

Go-fluxy is a simple library for creating Flux-compatible queries for InfluxDB V2.

## Installation

```bash
go get github.com/GLCharge/go-fluxy
```

## Usage

```go

import (
    "github.com/GLCharge/go-fluxy"
)

func main()  {
    query := NewFluxQueryBuilder().	From("mybucket").Range(time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC), nil).Filter("r._measurement == \"mymeasurement\"").GroupBy("something")
	fmt.Println(query.Build())
	
	// Output:
	// from(bucket: "mybucket")
	// 	|> range(start: 2022-01-01T00:00:00Z)
	// 	|> filter(fn: (r) => r._measurement == "mymeasurement")
	// 	|> group(by: ["something"])
	
	// Use the query to query InfluxDB
}
```