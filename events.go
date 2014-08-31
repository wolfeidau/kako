package kako

import (
	"time"

	"code.google.com/p/google-api-go-client/bigquery/v2"
)

type Event struct {
	UserID      string       // GUID of the user who owns this series
	Key         string       // name of the series, otherwise known as key
	Measurement *Measurement // measurement
	LogEntry    *LogEntry    // log event
	Timestamp   time.Time    // when measurement was taken
}

func (e *Event) ToRowMap() map[string]bigquery.JsonValue {

	v := make(map[string]bigquery.JsonValue)

	v["Key"] = e.Key
	v["UserID"] = e.UserID
	v["Timestamp"] = e.Timestamp

	if e.Measurement != nil {
		v["Measurement"] = e.Measurement.ToRowMap()
	}

	if e.LogEntry != nil {
		v["LogEntry"] = e.LogEntry.ToRowMap()
	}

	return v
}

// "time", "count", "min", "max", "mean", "std-dev", "95-percentile"
type Measurement struct {
	Min, Max, Mean, Count, Percentile95 int64
	Value                               float64
}

func (m *Measurement) ToRowMap() map[string]bigquery.JsonValue {
	v := make(map[string]bigquery.JsonValue)

	v["Min"] = m.Min
	v["Max"] = m.Max
	v["Mean"] = m.Mean
	v["Value"] = m.Value
	v["Count"] = m.Count
	v["Percentile95"] = m.Percentile95

	return v
}

type LogEntry struct {
	Facility, Severity, Tag, Content string
}

func (l *LogEntry) ToRowMap() map[string]bigquery.JsonValue {
	v := make(map[string]bigquery.JsonValue)

	v["Tag"] = l.Tag
	v["Content"] = l.Content
	v["Facility"] = l.Facility
	v["Severity"] = l.Severity

	return v
}
