# fullers-go

[![GoDoc](https://godoc.org/github.com/mattevans/fullers-go?status.svg)](https://godoc.org/github.com/mattevans/fullers-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattevans/fullers-go)](https://goreportcard.com/report/github.com/mattevans/fullers-go)

Scraps timetables and alerts of [Fullers](https://www.fullers.co.nz/) destinations from their site.

Installation
-----------------

`go get -u github.com/mattevans/fullers-go`

Example
-------------

Find supported destinations.

```go
client := fullers.NewClient()
locations, err := client.Location.GetLocations()
if err != nil {
  return err
}
```

Which will give you...

```json
[
    {
        "name": "Auckland City",
        "code": "AUCK"
    },
    {
        "name": "Bayswater",
        "code": "BAYS"
    },
    {
        "name": "Beach Haven",
        "code": "BEAC"
    },
    ...
]
```

Find a specific timetable.

```go
client := fullers.NewClient()
timetables, err := client.Timetable.GetTimetable(&fullers.TimetableRequest{
    Origin:      "BAYS",
    Destination: "AUCK",
})
```

Which will give you...

```json
[
    {
        "origin": "From: Bayswater",
        "destination": "To: Auckland City",
        "duration": "12 minutes",
        "column_titles": [
            "Mon to Fri",
            "Sat",
            "Sun & Pub Hols"
        ],
        "column_data": [
        [
            "6:40 AM",
            "7:10 AM",
            "7:40 AM",
            "8:10 AM",
        ...
```

Contributing
-----------------
If you've found a bug or would like to contribute, please create an issue here on GitHub, or better yet fork the project and submit a pull request!