package stw_brandenburg_west

import (
	"mensa-queue/internal/repository"
	"strings"
	"testing"
)

func TestParseTimeRange(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		from string
		to   string
		ok   bool
	}{
		{"en dash with Uhr", "08:00 – 18:00 Uhr", "08:00", "18:00", true},
		{"ascii hyphen", "8:00 - 18:00", "08:00", "18:00", true},
		{"serving prefix", "Essenausgabe 11:00 – 18:00 Uhr", "11:00", "18:00", true},
		{"serving prefix colon", "Essenausgabe:11:00 – 14:30 Uhr", "11:00", "14:30", true},
		{"surrounding whitespace", " 08:00 – 16:30 Uhr", "08:00", "16:30", true},
		{"single digit hour", "9:15 – 13:00", "09:15", "13:00", true},
		{"no range", "Nach Absprache", "", "", false},
		{"closed text", "Geschlossen", "", "", false},
		{"reversed range", "18:00 – 08:00", "", "", false},
		{"invalid time", "99:00 – 18:00", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to, ok := parseTimeRange(tt.raw)
			if ok != tt.ok {
				t.Fatalf("parseTimeRange(%q) ok = %v, want %v", tt.raw, ok, tt.ok)
			}
			if from != tt.from || to != tt.to {
				t.Fatalf("parseTimeRange(%q) = (%q, %q), want (%q, %q)", tt.raw, from, to, tt.from, tt.to)
			}
		})
	}
}

func strPtr(s string) *string { return &s }

func TestParseWeek(t *testing.T) {
	// Real provider data (Mensa Golm): Mon-Thu and Friday with service
	// times, weekend closed with differently cased "Geschlossen"
	golm := Outlet{
		MoZeit1: strPtr("08:00 – 18:00 Uhr"),
		MoZeit2: strPtr("Essenausgabe 11:00 – 18:00 Uhr"),
		DiZeit1: strPtr("08:00 – 18:00 Uhr"),
		DiZeit2: strPtr("Essenausgabe 11:00 – 18:00 Uhr"),
		MiZeit1: strPtr("08:00 – 18:00 Uhr"),
		MiZeit2: strPtr("Essenausgabe 11:00 – 18:00 Uhr"),
		DoZeit1: strPtr("08:00 – 18:00 Uhr"),
		DoZeit2: strPtr("Essenausgabe 11:00 – 18:00 Uhr"),
		FrZeit1: strPtr("08:00 – 14:30 Uhr"),
		FrZeit2: strPtr("Essenausgabe 11:00 – 14:00 Uhr"),
		SaZeit1: strPtr("Geschlossen"),
		SoZeit1: strPtr("geschlossen"),
	}

	week, err := golm.parseWeek()
	if err != nil {
		t.Fatalf("parseWeek() error: %v", err)
	}
	if len(week) != 5 {
		t.Fatalf("expected 5 open days, got %d", len(week))
	}
	if week[0].Day != repository.EnumOpeningDayDayMonday ||
		week[0].OpeningFrom != "08:00" || week[0].OpeningTo != "18:00" {
		t.Errorf("monday parsed as %+v", week[0])
	}
	if week[0].ServiceFrom == nil || *week[0].ServiceFrom != "11:00" ||
		week[0].ServiceTo == nil || *week[0].ServiceTo != "18:00" {
		t.Errorf("monday service parsed as %+v", week[0])
	}
	friday := week[4]
	if friday.Day != repository.EnumOpeningDayDayFriday || friday.OpeningTo != "14:30" ||
		friday.ServiceTo == nil || *friday.ServiceTo != "14:00" {
		t.Errorf("friday parsed as %+v", friday)
	}

	// Leading whitespace on " Geschlossen" (Kiepenheuerallee) and a
	// missing Zeit2 must still be treated as closed
	whitespace := Outlet{SaZeit1: strPtr(" Geschlossen")}
	week, err = whitespace.parseWeek()
	if err != nil {
		t.Fatalf("parseWeek() error: %v", err)
	}
	if len(week) != 0 {
		t.Fatalf("expected all days closed, got %+v", week)
	}
}

func TestParseWeekKeepsDataOnUnparsableOpeningTime(t *testing.T) {
	broken := Outlet{
		MoZeit1: strPtr("Nach Absprache geöffnet"),
	}
	if _, err := broken.parseWeek(); err == nil {
		t.Fatal("expected an error for an unparsable opening time")
	} else if !strings.Contains(err.Error(), "monday") {
		t.Errorf("error should mention the weekday, got: %v", err)
	}

	// An unparsable serving time only drops the service range
	fallback := Outlet{
		MoZeit1: strPtr("08:00 – 18:00 Uhr"),
		MoZeit2: strPtr("nach Absprache"),
	}
	week, err := fallback.parseWeek()
	if err != nil {
		t.Fatalf("parseWeek() error: %v", err)
	}
	if len(week) != 1 || week[0].ServiceFrom != nil {
		t.Fatalf("expected opening time without service range, got %+v", week)
	}
}

func TestFindOutlet(t *testing.T) {
	outlets := []Outlet{
		{Name: "Abholung / Pick up: Wohnheim Stahnsdorfer Straße", LocationInfo: OutletLocation{ID: 9601}},
		{Name: "Mensa Griebnitzsee", LocationInfo: OutletLocation{ID: 9601}},
		{Name: "Mensa Golm", LocationInfo: OutletLocation{ID: 9602}},
		{Name: "Cafeteria Am Neuen Palais", LocationInfo: OutletLocation{ID: 9607}},
	}

	outlet, err := findOutlet(outlets, 9601)
	if err != nil {
		t.Fatalf("findOutlet() error: %v", err)
	}
	if outlet.Name != "Mensa Griebnitzsee" {
		t.Fatalf("picked %q, want the mensa outlet", outlet.Name)
	}

	if _, err := findOutlet(outlets, 9603); err == nil {
		t.Fatal("expected an error for a location without a mensa outlet")
	}
}
