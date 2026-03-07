package web

import "testing"

func TestMusicFilters_NoActiveFilters(t *testing.T) {
	filters := musicFilters(nil, nil)
	for _, f := range filters {
		if f.Active {
			t.Errorf("%s should be inactive", f.Name)
		}
		expected := "/music?instruments=" + f.Name
		if f.Path != expected {
			t.Errorf("%s: expected path %q, got %q", f.Name, expected, f.Path)
		}
	}
}

func TestMusicFilters_ActiveInstrumentDeselects(t *testing.T) {
	filters := musicFilters([]string{"trumpet"}, nil)
	for _, f := range filters {
		if f.Name == "trumpet" {
			if !f.Active {
				t.Error("trumpet should be active")
			}
			if f.Path != "/music" {
				t.Errorf("active trumpet path should be /music, got %q", f.Path)
			}
		} else {
			if f.Active {
				t.Errorf("%s should be inactive", f.Name)
			}
		}
	}
}

func TestMusicFilters_PreservesComposerState(t *testing.T) {
	// active instrument = trumpet, active composer = !Adam Czerepinski
	filters := musicFilters([]string{"trumpet"}, []string{"!Adam Czerepinski"})
	for _, f := range filters {
		if f.Name == "trumpet" {
			// clicking active trumpet deselects it but keeps composer
			expected := "/music?composers=%21Adam+Czerepinski"
			if f.Path != expected {
				t.Errorf("active trumpet path: expected %q, got %q", expected, f.Path)
			}
		} else if f.Name == "piano" {
			// inactive instrument keeps composer state
			expected := "/music?instruments=piano&composers=%21Adam+Czerepinski"
			if f.Path != expected {
				t.Errorf("inactive piano path: expected %q, got %q", expected, f.Path)
			}
		}
	}
}

func TestComposerFilters_NoActiveFilters(t *testing.T) {
	filters := composerFilters(nil, nil)
	if len(filters) != 2 {
		t.Fatalf("expected 2 composer filters, got %d", len(filters))
	}
	for _, f := range filters {
		if f.Active {
			t.Errorf("%s should be inactive", f.Name)
		}
	}
	if filters[0].Name != "originals" {
		t.Errorf("first filter should be originals, got %q", filters[0].Name)
	}
	if filters[1].Name != "covers" {
		t.Errorf("second filter should be covers, got %q", filters[1].Name)
	}
}

func TestComposerFilters_OriginalsActive(t *testing.T) {
	filters := composerFilters([]string{"Adam Czerepinski"}, nil)
	for _, f := range filters {
		if f.Name == "originals" {
			if !f.Active {
				t.Error("originals should be active")
			}
			// clicking deselects, no instrument state to preserve
			if f.Path != "/music" {
				t.Errorf("active originals path should be /music, got %q", f.Path)
			}
		} else if f.Name == "covers" {
			if f.Active {
				t.Error("covers should be inactive")
			}
		}
	}
}

func TestComposerFilters_CoversActive(t *testing.T) {
	filters := composerFilters([]string{"!Adam Czerepinski"}, nil)
	for _, f := range filters {
		if f.Name == "covers" {
			if !f.Active {
				t.Error("covers should be active")
			}
			if f.Path != "/music" {
				t.Errorf("active covers path should be /music, got %q", f.Path)
			}
		}
	}
}

func TestComposerFilters_PreservesInstrumentState(t *testing.T) {
	// active composer = Adam Czerepinski, active instrument = trumpet
	filters := composerFilters([]string{"Adam Czerepinski"}, []string{"trumpet"})
	for _, f := range filters {
		if f.Name == "originals" {
			// clicking active originals deselects it but keeps instrument
			expected := "/music?instruments=trumpet"
			if f.Path != expected {
				t.Errorf("active originals path: expected %q, got %q", expected, f.Path)
			}
		} else if f.Name == "covers" {
			// inactive covers keeps instrument state
			expected := "/music?instruments=trumpet&composers=%21Adam+Czerepinski"
			if f.Path != expected {
				t.Errorf("inactive covers path: expected %q, got %q", expected, f.Path)
			}
		}
	}
}
