package fullers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snikch/api/fail"
)

// LocationService ...
type LocationService service

// Location ...
type Location struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var locations []Location

// BuildLocationRequestURL ...
func (s *LocationService) BuildLocationRequestURL() string {
	return s.client.BackendURL.String()
}

// InitLocations ...
func (s *LocationService) InitLocations() {
	var err error
	locations, err = s.GetLocations()
	if err != nil {
		log.Fatal("Unable to load fullers locations")
	}
}

// GetLocations ...
func (s *LocationService) GetLocations() ([]Location, error) {
	doc, err := goquery.NewDocument(s.BuildLocationRequestURL())
	if err != nil {
		return locations, err
	}

	doc.Find("#timetableDate .departure select option").Each(func(i int, s *goquery.Selection) {
		attr, _ := s.Attr("value")
		if attr != "" {
			loc := Location{
				Name: s.Text(),
				Code: attr,
			}
			locations = append(locations, loc)
		}
	})

	return locations, err
}

// ValidLocation ...
func (s *LocationService) ValidLocation(locs []string) error {
	if len(locations) <= 0 {
		return fail.NewValidationError(errors.New("Unable to parse list of fullers locations"))
	}
	// Map our locations we have in memory.
	locationWhitelist := map[string]bool{}
	for _, l := range locations {
		locationWhitelist[strings.ToUpper(l.Code)] = true
	}
	// Range input locations, check we have a match.
	for _, location := range locs {
		if !locationWhitelist[strings.ToUpper(location)] {
			return fail.NewValidationError(fmt.Errorf("Unknown fullers location: %s, please try again", location))
		}
	}
	return nil
}
