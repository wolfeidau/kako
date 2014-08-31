# kako

Kako is a time series store which accepts data over HTTPS and streams it into bigquery.

## type Event
``` go
type Event struct {
    UserID      string       // GUID of the user who owns this series
    Key         string       // name of the series, otherwise known as key
    Measurement *Measurement // measurement
    LogEntry    *LogEntry    // log event
    Timestamp   time.Time    // when measurement was taken
}
```

# usage

``` go
import "."
```

### func (\*Event) ToRowMap
``` go
func (e *Event) ToRowMap() map[string]bigquery.JsonValue
```

## type LogEntry
``` go
type LogEntry struct {
    Facility, Severity, Tag, Content string
}
```

### func (\*LogEntry) ToRowMap
``` go
func (l *LogEntry) ToRowMap() map[string]bigquery.JsonValue
```

## type Measurement
``` go
type Measurement struct {
    Min, Max, Mean, Count, Percentile95 int64
    Value                               float64
}
```
"time", "count", "min", "max", "mean", "std-dev", "95-percentile"

### func (\*Measurement) ToRowMap
``` go
func (m *Measurement) ToRowMap() map[string]bigquery.JsonValue
```