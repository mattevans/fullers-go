package fullers

import (
	"net/url"
)

const (
	backendURL = "https://www.fullers.co.nz/timetables-and-fares"
)

// Client ...
type Client struct {
	BackendURL *url.URL
	Location   *LocationService
	Timetable  *TimetableService
}

type service struct {
	client *Client
}

// NewClient ...
func NewClient() *Client {
	baseURL, _ := url.Parse(backendURL)

	c := &Client{
		BackendURL: baseURL,
	}

	c.Location = &LocationService{client: c}
	c.Timetable = &TimetableService{client: c}

	c.Location.InitLocations()

	return c
}
