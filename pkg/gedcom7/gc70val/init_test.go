package gc70val

import (
	"sort"
	"testing"
)

func TestInit(t *testing.T) {
	tlLen := 183
	calLen := 4
	typeLen := 6
	esLen := 12

	if len(baseline.tags) == 0 {
		t.Fatal("no tags discovered")
	}
	if len(baseline.tags) != tlLen {
		t.Errorf("tagList length mismatch. Wanted %d; got %d", tlLen, len(baseline.tags))
		keys := make([]string, 0)
		for k, v := range baseline.tags {
			if k != v.FullTag {
				t.Errorf("tagList key mismatch. Wanted %s; got %s", k, v.FullTag)
			}
			keys = append(keys, string(k))
		}
		sort.Strings(keys)
		t.Logf("tagList keys: %v", keys)
	}

	if len(baseline.enumSets) != esLen {
		var keys string
		for k := range baseline.enumSets {
			keys += k + " "
		}
		t.Errorf("enumSets length mismatch. Wanted %d; got %d\n%s", esLen, len(baseline.enumSets), keys)
	}

	if len(baseline.calendars) == 0 {
		t.Fatal("calendars not initialized")
	}
	if len(baseline.calendars) != calLen {
		t.Errorf("calendars length mismatch. Wanted %d; got %d", calLen, len(baseline.calendars))
	}
	for k, v := range baseline.calendars {
		if k != v.Cal {
			t.Errorf("calendars key mismatch. Wanted %s; got %s", k, v.Cal)
		}
		if len(v.Months) == 0 {
			t.Errorf("calendars month list not initialized for %s", k)
		}
	}

	if len(baseline.types) == 0 {
		t.Fatal("types not initialized")
	}
	if len(baseline.types) != typeLen {
		t.Errorf("types length mismatch. Wanted %d; got %d", typeLen, len(baseline.types))
	}

	bl := baseline
	for k, v := range bl.tags {
		if v.EnumSetName != "" && len(v.EnumSet.Values) == 0 {
			t.Errorf("enumSet not initialized for %s", k)
		}
	}
}
