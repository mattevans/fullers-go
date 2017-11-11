package fullers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// TimetableService ...
type TimetableService service

// TimetableRequest ...
type TimetableRequest struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

// Timetable ...
type Timetable struct {
	Origin       string              `json:"origin"`
	Destination  string              `json:"destination"`
	Duration     string              `json:"duration"`
	ColumnTitles []string            `json:"column_titles"`
	ColumnData   [][]string          `json:"column_data"`
	Footnotes    []string            `json:"footnotes"`
	ColumnCount  int                 `json:"column_count"`
	RowCount     int                 `json:"row_count"`
	Alerts       []map[string]string `json:"alerts"`
}

// BuildTimetableAlertRequestURL ...
func (s *TimetableService) BuildTimetableAlertRequestURL(request *TimetableRequest) string {
	return fmt.Sprintf("%s/alerts/?from=%s&to=%s", s.client.BackendURL, strings.ToUpper(request.Origin), strings.ToUpper(request.Destination))
}

// BuildTimetableRequestURL ...
func (s *TimetableService) BuildTimetableRequestURL(request *TimetableRequest) string {
	return fmt.Sprintf("%s/?from=%s&to=%s", s.client.BackendURL, strings.ToUpper(request.Origin), strings.ToUpper(request.Destination))
}

// GetTimetable ...
func (s *TimetableService) GetTimetable(request *TimetableRequest) ([]Timetable, error) {
	if request == nil {
		return nil, errors.New("The timetable request cannot be nil")
	}

	// Check we have valid origin/destinations passed.
	err := s.client.Location.ValidLocation([]string{
		request.Origin,
		request.Destination,
	})
	if err != nil {
		return nil, err
	}

	// Scrape any alerts first.
	adoc, err := goquery.NewDocument(s.BuildTimetableAlertRequestURL(request))
	if err != nil {
		return nil, err
	}
	alerts := []map[string]string{}
	adoc.Find("section.alert").Each(func(i int, s *goquery.Selection) {
		alert := map[string]string{
			"title":   s.Find(".alert-title").Text(),
			"content": s.Find("p").Text(),
		}
		alerts = append(alerts, alert)
	})

	// Scrape timetables, build into goquery document.
	doc, err := goquery.NewDocument(s.BuildTimetableRequestURL(request))
	if err != nil {
		return nil, err
	}

	timetables := []Timetable{}
	doc.Find(".timetable-content .timetable-inner").Each(func(i int, s *goquery.Selection) {

		// Build our Timetable object.
		timetable := Timetable{
			Origin:       strings.TrimSpace(s.Find("div.location").First().Text()),
			Destination:  strings.TrimSpace(s.Find("div.destination").Text()),
			Duration:     strings.TrimSpace(s.Find("div.time").Text()),
			Footnotes:    []string{},
			ColumnData:   [][]string{},
			ColumnTitles: []string{},
			ColumnCount:  0,
			Alerts:       alerts,
		}

		// Range each timetable found.
		s.Find(".timetable-column").Each(func(k int, is *goquery.Selection) {
			// Nudge our column count.
			timetable.ColumnCount = timetable.ColumnCount + 1
			timetable.ColumnTitles = append(timetable.ColumnTitles, is.Find("span.day").Text())

			// Find the departure times.
			times := is.Find("span").Parent().NextAll().Text()

			filteredTimes := []string{}
			splitTimes := strings.Split(times, "\n")
			for _, t := range splitTimes {
				time := strings.TrimSpace(t)
				if time != "" {
					filteredTimes = append(filteredTimes, time)
				}
			}
			ft := [][]string{
				filteredTimes,
			}
			timetable.ColumnData = append(timetable.ColumnData, ft...)
			timetable.RowCount = len(filteredTimes)
		})

		// Assign any footnotes we might have.
		s.Find(".footnote").Each(func(k int, is *goquery.Selection) {
			fn := strings.TrimSpace(is.Text())
			timetable.Footnotes = append(timetable.Footnotes, fn)
		})

		// Append the timetable to our slice.
		timetables = append(timetables, timetable)
	})

	// Return.
	return timetables, err
}
