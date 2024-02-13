package main

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CachedAd struct {
	ad       Ad
	gender   map[string]bool
	country  map[string]bool
	platform map[string]bool
}

type Cache struct {
	genderIndex   map[string](map[*CachedAd]bool)
	countryIndex  map[string](map[*CachedAd]bool)
	platformIndex map[string](map[*CachedAd]bool)
	ageIndex      []([]*CachedAd) // 0-100
	ads           []*CachedAd
	lock          sync.RWMutex
}

// updateFromCollection updates the cache with active ads from the given MongoDB Collection.
// Returns true if the update is successful, false otherwise.
func (cache *Cache) updateFromCollection(coll *mongo.Collection) bool {
	now := time.Now()

	filter := bson.M{
		"startat": bson.M{"$lte": now},
		"endat":   bson.M{"$gte": now},
	}
	opts := options.Find().SetSort(bson.D{{Key: "startat", Value: 1}})

	cur, err := coll.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println("Error fetching active ads from DB:", err)
		return false
	}

	var results []Ad
	for cur.Next(context.Background()) {
		var ad Ad
		if err := cur.Decode(&ad); err != nil {
			log.Println("Error decoding ad:", err)
			continue
		}
		results = append(results, ad)
	}

	if err := cur.Err(); err != nil {
		log.Println("Cursor error:", err)
	}
	cur.Close(context.Background())

	cache.update(results)

	return true
}

// update updates the cache with the given ads.
// It populates the genderIndex, countryIndex, platformIndex, ageIndex, and ads fields of the Cache.
// The ads parameter is a slice of Ad structs representing the ads to be added to the cache.
// Each ad is processed to update the corresponding indexes and maps in the cache.
// The function is thread-safe and uses a lock to ensure concurrent access to the cache is synchronized.
func (cache *Cache) update(ads []Ad) {
	log.Println("Updating cache ...")

	cachedAds := make([]*CachedAd, 0, len(ads))
	genderIndex := make(map[string](map[*CachedAd]bool))
	countryIndex := make(map[string](map[*CachedAd]bool))
	platformIndex := make(map[string](map[*CachedAd]bool))
	ageIndex := make([]([]*CachedAd), 101)
	for _, ad := range ads {
		cachedAd := &CachedAd{
			ad:       ad,
			gender:   make(map[string]bool),
			country:  make(map[string]bool),
			platform: make(map[string]bool),
		}

		if ad.Conditions.Gender != nil {
			for _, g := range *ad.Conditions.Gender {
				if genderIndex[g] == nil {
					genderIndex[g] = make(map[*CachedAd]bool)
				}
				genderIndex[g][cachedAd] = true
				cachedAd.gender[g] = true
			}
		}

		if ad.Conditions.Country != nil {
			for _, c := range *ad.Conditions.Country {
				if countryIndex[c] == nil {
					countryIndex[c] = make(map[*CachedAd]bool)
				}
				countryIndex[c][cachedAd] = true
				cachedAd.country[c] = true
			}
		}

		if ad.Conditions.Platform != nil {
			for _, p := range *ad.Conditions.Platform {
				if platformIndex[p] == nil {
					platformIndex[p] = make(map[*CachedAd]bool)
				}
				platformIndex[p][cachedAd] = true
				cachedAd.platform[p] = true
			}
		}

		if ad.Conditions.AgeStart != nil && ad.Conditions.AgeEnd != nil {
			for age := *ad.Conditions.AgeStart; age <= *ad.Conditions.AgeEnd; age++ {
				if ageIndex[age] == nil {
					ageIndex[age] = make([]*CachedAd, 0)
				}
				ageIndex[age] = append(ageIndex[age], cachedAd)
			}
		}

		cachedAds = append(cachedAds, cachedAd)
	}

	cache.lock.Lock()
	cache.genderIndex = genderIndex
	cache.countryIndex = countryIndex
	cache.platformIndex = platformIndex
	cache.ageIndex = ageIndex
	cache.ads = cachedAds
	cache.lock.Unlock()

	log.Println("Cache updated", len(cache.ads), "ads")
}

// Continuously updates the cache from the collection of ads at the specified time interval.
// It takes a TTL duration as a parameter determines how often the cache should be updated.
// The updater function runs indefinitely.
func (cache *Cache) updater(ttl time.Duration) {
	for {
		cache.updateFromCollection(ads)
		time.Sleep(ttl)
	}
}

// filter applies the given query to the cache and returns a list of matching ads.
func (cache *Cache) filter(query AdQuery) []Ad {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	matchingAds := make([]*CachedAd, 0, len(cache.ads))
	if query.Age > 0 {
		if cache.ageIndex[query.Age] != nil {
			matchingAds = append(matchingAds, cache.ageIndex[query.Age]...)
		}
	} else {
		matchingAds = append(matchingAds, cache.ads...)
	}

	if query.Country != "" {
		if cache.countryIndex[query.Country] != nil {
			for i, ad := range matchingAds {
				if !cache.countryIndex[query.Country][ad] {
					matchingAds[i] = nil
				}
			}
		} else {
			matchingAds = make([]*CachedAd, 0)
		}
	}

	if query.Platform != "" {
		if cache.platformIndex[query.Platform] != nil {
			for i, ad := range matchingAds {
				if !cache.platformIndex[query.Platform][ad] {
					matchingAds[i] = nil
				}
			}
		} else {
			matchingAds = make([]*CachedAd, 0)
		}
	}

	if query.Gender != "" {
		if cache.genderIndex[query.Gender] != nil {
			for i, ad := range matchingAds {
				if !cache.genderIndex[query.Gender][ad] {
					matchingAds[i] = nil
				}
			}
		} else {
			matchingAds = make([]*CachedAd, 0)
		}
	}

	var skiped int64 = 0
	results := []Ad{}
	for _, ad := range matchingAds {
		if ad == nil {
			continue
		}
		if skiped < query.Offset {
			skiped++
			continue
		}
		results = append(results, ad.ad)
		if int64(len(results)) >= query.Limit {
			break
		}
	}

	return results
}
