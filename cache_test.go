package main

import (
	"testing"
	"time"
)

func TestCacheFilter(t *testing.T) {
	cache := Cache{}
	setupCache(&cache)

	t.Run("Filter by Limit", func(t *testing.T) {
		query := AdQuery{Offset: 0, Limit: 3}
		results := cache.Filter(query)
		if len(results) != 3 {
			t.Errorf("Expected 3 ads, got %d", len(results))
		}
		if results[0].Title != "Ad 1" || results[1].Title != "Ad 2" || results[2].Title != "Ad 3" {
			t.Errorf("Expected Ad 1, Ad 2, Ad 3, got %s, %s, %s", results[0].Title, results[1].Title, results[2].Title)
		}
	})

	t.Run("Filter by Offset", func(t *testing.T) {
		query := AdQuery{Offset: 1, Limit: 3}
		results := cache.Filter(query)
		if len(results) != 3 {
			t.Errorf("Expected 3 ads, got %d", len(results))
		}
		if results[0].Title != "Ad 2" || results[1].Title != "Ad 3" || results[2].Title != "Ad 4" {
			t.Errorf("Expected Ad 2, Ad 3, Ad 4, got %s, %s, %s", results[0].Title, results[1].Title, results[2].Title)
		}
	})

	t.Run("Filter by out of range Offset", func(t *testing.T) {
		query := AdQuery{Offset: 5, Limit: 3}
		results := cache.Filter(query)
		if len(results) != 0 {
			t.Errorf("Expected 0 ad, got %d", len(results))
		}
	})

	t.Run("Filter by Age", func(t *testing.T) {
		query := AdQuery{Offset: 0, Limit: 3, Age: 25}
		results := cache.Filter(query)
		if len(results) != 2 {
			t.Errorf("Expected 2 ads, got %d", len(results))
		}
		if results[0].Title != "Ad 2" || results[1].Title != "Ad 4" {
			t.Errorf("Expected Ad 2, Ad 4, got %s, %s", results[0].Title, results[1].Title)
		}
	})

	t.Run("Filter by Country", func(t *testing.T) {
		query := AdQuery{Offset: 0, Limit: 3, Country: "TW"}
		results := cache.Filter(query)
		if len(results) != 2 {
			t.Errorf("Expected 2 ads, got %d", len(results))
		}
		if results[0].Title != "Ad 1" || results[1].Title != "Ad 3" {
			t.Errorf("Expected Ad 1, Ad 3, got %s, %s", results[0].Title, results[1].Title)
		}
	})

	t.Run("Filter by Platform", func(t *testing.T) {
		query := AdQuery{Offset: 0, Limit: 3, Platform: "android"}
		results := cache.Filter(query)
		if len(results) != 2 {
			t.Errorf("Expected 2 ads, got %d", len(results))
		}
		if results[0].Title != "Ad 3" || results[1].Title != "Ad 4" {
			t.Errorf("Expected Ad 3, Ad 4, got %s, %s", results[0].Title, results[1].Title)
		}
	})

	t.Run("Filter by Country and Gender", func(t *testing.T) {
		query := AdQuery{Offset: 0, Limit: 3, Country: "TW", Gender: "M"}
		results := cache.Filter(query)
		if len(results) != 1 {
			t.Errorf("Expected 1 ad, got %d", len(results))
		}
		if results[0].Title != "Ad 1" {
			t.Errorf("Expected Ad 1, got %s", results[0].Title)
		}
	})

	t.Run("Filter by Country and Platform", func(t *testing.T) {
		query := AdQuery{Offset: 0, Limit: 3, Country: "JP", Platform: "ios"}
		results := cache.Filter(query)
		if len(results) != 1 {
			t.Errorf("Expected 1 ad, got %d", len(results))
		}
		if results[0].Title != "Ad 4" {
			t.Errorf("Expected Ad 4, got %s", results[0].Title)
		}
	})

	t.Run("Filter by Age and Platform", func(t *testing.T) {
		query := AdQuery{Offset: 0, Limit: 3, Age: 30, Platform: "web"}
		results := cache.Filter(query)
		if len(results) != 1 {
			t.Errorf("Expected 1 ad, got %d", len(results))
		}
		if results[0].Title != "Ad 4" {
			t.Errorf("Expected Ad 4, got %s", results[0].Title)
		}
	})

	t.Run("Filter by Age, Country, and Platform", func(t *testing.T) {
		query := AdQuery{Offset: 0, Limit: 3, Age: 30, Country: "TW", Platform: "web"}
		results := cache.Filter(query)
		if len(results) != 0 {
			t.Errorf("Expected 0 ad, got %d", len(results))
		}
	})
}

func setupCache(cache *Cache) {
	int20 := 20
	int30 := 30
	int40 := 40
	ad1 := Ad{
		Title:   "Ad 1",
		StartAt: time.Now().Add(-time.Hour * 24),
		EndAt:   time.Now().Add(time.Hour * 24),
		Conditions: AdConditions{
			Gender:  &[]string{"M"},
			Country: &[]string{"TW", "JP"},
		},
	}
	ad2 := Ad{
		Title:   "Ad 2",
		StartAt: time.Now().Add(-time.Hour * 24),
		EndAt:   time.Now().Add(time.Hour * 24),
		Conditions: AdConditions{
			AgeStart: &int20,
			AgeEnd:   &int30,
			Gender:   &[]string{"F", "M"},
		},
	}
	ad3 := Ad{
		Title:   "Ad 3",
		StartAt: time.Now().Add(-time.Hour * 24),
		EndAt:   time.Now().Add(time.Hour * 24),
		Conditions: AdConditions{
			AgeStart: &int30,
			AgeEnd:   &int40,
			Country:  &[]string{"TW"},
			Platform: &[]string{"android", "ios"},
		},
	}
	ad4 := Ad{
		Title:   "Ad 4",
		StartAt: time.Now().Add(-time.Hour * 24),
		EndAt:   time.Now().Add(time.Hour * 24),
		Conditions: AdConditions{
			AgeStart: &int20,
			AgeEnd:   &int40,
			Country:  &[]string{"JP"},
			Platform: &[]string{"android", "ios", "web"},
		},
	}

	cache.Update([]Ad{ad1, ad2, ad3, ad4})
}
