package main

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CachedAd struct {
	ad       Ad
	gender   map[string]bool
	country  map[string]bool
	platform map[string]bool
}

type Cache struct {
	Ads  []CachedAd
	lock sync.RWMutex
}

var (
	cache = Cache{}
)

func cacheUpdater(ttl time.Duration) {
	for {
		log.Println("Updating cache")
		now := time.Now()

		filter := bson.M{
			"startat": bson.M{"$lte": now},
			"endat":   bson.M{"$gte": now},
		}
		opts := options.Find().SetSort(bson.D{{Key: "startat", Value: 1}})

		cur, err := ads.Find(context.Background(), filter, opts)
		if err != nil {
			log.Println("Error fetching active ads from DB:", err)
			continue
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

		cachedAds := make([]CachedAd, len(results))
		for i, ad := range results {
			genderMap := make(map[string]bool)
			if ad.Conditions.Gender != nil {
				for _, g := range *ad.Conditions.Gender {
					genderMap[g] = true
				}
			}
			countryMap := make(map[string]bool)
			if ad.Conditions.Country != nil {
				for _, c := range *ad.Conditions.Country {
					countryMap[c] = true
				}
			}
			platformMap := make(map[string]bool)
			if ad.Conditions.Platform != nil {
				for _, p := range *ad.Conditions.Platform {
					platformMap[p] = true
				}
			}
			cachedAds[i] = CachedAd{
				ad:       ad,
				gender:   genderMap,
				country:  countryMap,
				platform: platformMap,
			}
		}

		cache.lock.Lock()
		cache.Ads = cachedAds
		cache.lock.Unlock()

		log.Println("Cache updated", len(results), "ads")
		time.Sleep(ttl)
	}
}

func filterFromCache(query AdQuery) []Ad {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	var skiped int64 = 0
	results := []Ad{}
	for _, cached := range cache.Ads {
		if query.Age != 0 && cached.ad.Conditions.AgeStart != nil && cached.ad.Conditions.AgeEnd != nil {
			if query.Age < *cached.ad.Conditions.AgeStart || query.Age > *cached.ad.Conditions.AgeEnd {
				continue
			}
		}

		if query.Country != "" {
			if !cached.country[query.Country] {
				continue
			}
		}

		if query.Gender != "" {
			if !cached.gender[query.Gender] {
				continue
			}
		}

		if query.Platform != "" {
			if !cached.platform[query.Platform] {
				continue
			}
		}

		if skiped < query.Offset {
			skiped++
			continue
		}

		results = append(results, cached.ad)
		if int64(len(results)) >= query.Limit {
			break
		}
	}

	return results
}
