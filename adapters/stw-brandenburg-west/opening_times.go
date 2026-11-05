package stw_brandenburg_west

import (
	"context"
	"fmt"
	"log"
	"mensa-queue/internal/repository"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

// OutletModel returns every outlet of the provider in one response. The
// location query parameter is ignored by the API, so a single request is
// enough for all mensas.
const OutletModel Model = "outlet"

type Outlet struct {
	Name         string         `json:"name"`
	MoZeit1      *string        `json:"moZeit1"`
	MoZeit2      *string        `json:"moZeit2"`
	DiZeit1      *string        `json:"diZeit1"`
	DiZeit2      *string        `json:"diZeit2"`
	MiZeit1      *string        `json:"miZeit1"`
	MiZeit2      *string        `json:"miZeit2"`
	DoZeit1      *string        `json:"doZeit1"`
	DoZeit2      *string        `json:"doZeit2"`
	FrZeit1      *string        `json:"frZeit1"`
	FrZeit2      *string        `json:"frZeit2"`
	SaZeit1      *string        `json:"saZeit1"`
	SaZeit2      *string        `json:"saZeit2"`
	SoZeit1      *string        `json:"soZeit1"`
	SoZeit2      *string        `json:"soZeit2"`
	LocationInfo OutletLocation `json:"locationInfo"`
}

type OutletLocation struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// DayOpening is one open weekday with its parsed times. Closed weekdays are
// simply not part of the list.
type DayOpening struct {
	Day         repository.EnumOpeningDayDay
	OpeningFrom string
	OpeningTo   string
	ServiceFrom *string
	ServiceTo   *string
}

// The provider free-text fields vary: en dashes, optional "Uhr", optional
// "Essenausgabe" prefix with or without a space, single digit hours and
// stray whitespace (" Geschlossen").
var timeRangePattern = regexp.MustCompile(`(\d{1,2}:\d{2})\s*[–\-]\s*(\d{1,2}:\d{2})`)

// ParseOutlets fetches every outlet of the provider with one request.
func ParseOutlets() ([]Outlet, error) {
	// The location parameter does not filter the response, any value works.
	location := int32(9602)
	response, err := ParseModel[Outlet](OutletModel, &location, repository.EnumLocaleLocaleDe)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, fmt.Errorf("outlet API error: success is false")
	}
	return response.Content, nil
}

// findOutlet picks the canteen outlet for a standort. Pickup points
// ("Abholung / Pick up: ...") and cafeterias share the standort ID with the
// mensa, so only outlets named "Mensa ..." are considered.
func findOutlet(outlets []Outlet, stwID int32) (*Outlet, error) {
	for i := range outlets {
		outlet := &outlets[i]
		if outlet.LocationInfo.ID == stwID && strings.HasPrefix(outlet.Name, "Mensa") {
			return outlet, nil
		}
	}
	return nil, fmt.Errorf("no outlet named 'Mensa' found for location %d", stwID)
}

// dayStrings returns Zeit1 (opening) and Zeit2 (serving/Essenausgabe) per
// weekday, monday through sunday.
func (o *Outlet) dayStrings() [7][2]*string {
	return [7][2]*string{
		{o.MoZeit1, o.MoZeit2},
		{o.DiZeit1, o.DiZeit2},
		{o.MiZeit1, o.MiZeit2},
		{o.DoZeit1, o.DoZeit2},
		{o.FrZeit1, o.FrZeit2},
		{o.SaZeit1, o.SaZeit2},
		{o.SoZeit1, o.SoZeit2},
	}
}

func ptrString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func isClosed(raw string) bool {
	return strings.Contains(strings.ToLower(raw), "geschlossen")
}

// parseTimeRange extracts and normalizes "08:00 – 18:00 Uhr" into
// ("08:00", "18:00").
func parseTimeRange(raw string) (from string, to string, ok bool) {
	match := timeRangePattern.FindStringSubmatch(raw)
	if match == nil {
		return "", "", false
	}
	from = normalizeTime(match[1])
	to = normalizeTime(match[2])
	if from == "" || to == "" || from >= to {
		return "", "", false
	}
	return from, to, true
}

// normalizeTime pads a single digit hour, e.g. "8:05" -> "08:05".
func normalizeTime(raw string) string {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return ""
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return ""
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return ""
	}
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

var weekDays = [7]repository.EnumOpeningDayDay{
	repository.EnumOpeningDayDayMonday,
	repository.EnumOpeningDayDayTuesday,
	repository.EnumOpeningDayDayWednesday,
	repository.EnumOpeningDayDayThursday,
	repository.EnumOpeningDayDayFriday,
	repository.EnumOpeningDayDaySaturday,
	repository.EnumOpeningDayDaySunday,
}

// parseWeek turns the free-text Zeit1/Zeit2 fields into structured weekdays.
// Closed weekdays are omitted. An unparsable opening time aborts the parse so
// that existing data is never wiped by a provider format change; an
// unparsable serving time only drops the (optional) service range.
func (o *Outlet) parseWeek() ([]DayOpening, error) {
	var days []DayOpening

	for i, pair := range o.dayStrings() {
		rawOpening := ptrString(pair[0])
		rawService := ptrString(pair[1])
		dayName := string(weekDays[i])

		if rawOpening == "" {
			if rawService != "" {
				return nil, fmt.Errorf("%s: serving time %q without opening time", dayName, rawService)
			}
			// No data for this day: closed
			continue
		}
		if isClosed(rawOpening) {
			continue
		}

		from, to, ok := parseTimeRange(rawOpening)
		if !ok {
			return nil, fmt.Errorf("%s: unparsable opening time %q", dayName, rawOpening)
		}

		day := DayOpening{
			Day:         weekDays[i],
			OpeningFrom: from,
			OpeningTo:   to,
		}

		if rawService != "" && !isClosed(rawService) {
			serviceFrom, serviceTo, ok := parseTimeRange(rawService)
			if !ok {
				log.Printf("opening times: %s: dropping unparsable serving time %q", dayName, rawService)
			} else {
				day.ServiceFrom = &serviceFrom
				day.ServiceTo = &serviceTo
			}
		}

		days = append(days, day)
	}

	return days, nil
}

// SyncOpeningTimes parses the weekly opening times of every registered mensa
// and stores them as shared opening_day rows referenced from mensa_rels.
// Identical day patterns are stored once and shared between mensas.
func (a *StwBrandenburgWestAdapter) SyncOpeningTimes(ctx context.Context, conn *pgx.Conn) error {
	outlets, err := ParseOutlets()
	if err != nil {
		return fmt.Errorf("fetching outlets: %w", err)
	}

	repo := repository.New(conn)

	for _, mensa := range a.Mensas {
		if !mensa.IsRegistered() || mensa.MensaHubID == nil {
			continue
		}

		outlet, err := findOutlet(outlets, mensa.StwID)
		if err != nil {
			log.Printf("opening times: mensa %s: %v", mensa.Name, err)
			continue
		}

		week, err := outlet.parseWeek()
		if err != nil {
			// Keep the previously stored times instead of storing garbage
			log.Printf("opening times: mensa %s: %v", mensa.Name, err)
			continue
		}

		var dayIDs []int32
		for _, day := range week {
			id, err := repo.FindOrCreateOpeningDay(ctx, repository.FindOrCreateOpeningDayParams{
				Day:         day.Day,
				OpeningFrom: day.OpeningFrom,
				OpeningTo:   day.OpeningTo,
				ServiceFrom: day.ServiceFrom,
				ServiceTo:   day.ServiceTo,
			})
			if err != nil {
				return fmt.Errorf("storing opening day for %s (%s): %w", mensa.Name, day.Day, err)
			}
			dayIDs = append(dayIDs, id)
		}

		currentPointers, err := repo.GetMensaOpeningDayIDs(ctx, *mensa.MensaHubID)
		if err != nil {
			return fmt.Errorf("reading opening days for %s: %w", mensa.Name, err)
		}
		var currentIDs []int32
		for _, id := range currentPointers {
			if id != nil {
				currentIDs = append(currentIDs, *id)
			}
		}

		// Only rewrite the relations when the week actually changed
		if slices.Equal(currentIDs, dayIDs) {
			continue
		}

		if err := repo.DeleteMensaOpeningDayRels(ctx, *mensa.MensaHubID); err != nil {
			return fmt.Errorf("clearing opening days for %s: %w", mensa.Name, err)
		}
		for order, dayID := range dayIDs {
			order := int32(order)
			dayID := dayID
			if err := repo.InsertMensaOpeningDayRel(ctx, repository.InsertMensaOpeningDayRelParams{
				ParentID:     *mensa.MensaHubID,
				Order:        &order,
				OpeningDayID: &dayID,
			}); err != nil {
				return fmt.Errorf("linking opening day for %s: %w", mensa.Name, err)
			}
		}

		fmt.Printf("Updated opening times for %s (%d days)\n", mensa.Name, len(dayIDs))
	}

	return nil
}
