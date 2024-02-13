package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type AdCondition struct {
	AgeStart *int      `json:"ageStart,omitempty"`
	AgeEnd   *int      `json:"ageEnd,omitempty"`
	Gender   *[]string `json:"gender,omitempty"`
	Country  *[]string `json:"country,omitempty"`
	Platform *[]string `json:"platform,omitempty"`
}

type AdPayload struct {
	Title      string      `json:"title"`
	StartAt    string      `json:"startAt"`
	EndAt      string      `json:"endAt"`
	Conditions AdCondition `json:"conditions"`
}

var r *rand.Rand

func generateAd(index int) *AdPayload {
	StartAt, EndAt := randomDateRange()

	Conditions := AdCondition{}

	if r.Intn(2) == 0 {
		AgeStart, AgeEnd := randomAgeRange()
		Conditions.AgeStart = &AgeStart
		Conditions.AgeEnd = &AgeEnd
	}

	genders := randomGenders()
	if len(genders) > 0 {
		Conditions.Gender = &genders
	}

	countries := randomCountries()
	if len(countries) > 0 {
		Conditions.Country = &countries
	}

	platforms := randomPlatforms()
	if len(platforms) > 0 {
		Conditions.Platform = &platforms
	}

	return &AdPayload{
		Title:      fmt.Sprintf("AD %d", index),
		StartAt:    StartAt,
		EndAt:      EndAt,
		Conditions: Conditions,
	}
}

func randomDateRange() (string, string) {
	start := time.Now().Add(time.Duration(-24*r.Intn(30)) * time.Hour)
	end := start.Add(time.Duration(24*r.Intn(30)) * time.Hour)
	return start.Format(time.RFC3339), end.Format(time.RFC3339)
}

func randomAgeRange() (int, int) {
	s := r.Intn(40)
	e := s + 1 + r.Intn(40)
	return s, e
}

func randomGenders() []string {
	genders := []string{"M", "F"}
	count := r.Intn(2)

	for i := range genders {
		j := r.Intn(i + 1)
		genders[i], genders[j] = genders[j], genders[i]
	}

	return genders[:count]
}

func randomCountries() []string {
	countries := []string{"US", "CA", "MX", "BR", "JP", "KR", "CN", "RU", "AU", "NZ", "TW"}
	count := r.Intn(8)

	for i := range countries {
		j := r.Intn(i + 1)
		countries[i], countries[j] = countries[j], countries[i]
	}

	return countries[:count]
}

func randomPlatforms() []string {
	platforms := []string{"android", "ios", "web"}
	count := r.Intn(3)

	for i := range platforms {
		j := r.Intn(i + 1)
		platforms[i], platforms[j] = platforms[j], platforms[i]
	}

	return platforms[:count]
}

func sendAd(ad *AdPayload) {
	jsonData, err := json.Marshal(ad)
	if err != nil {
		fmt.Println("Error encoding JSON", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/api/v1/ad", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending ad to API", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Ad sent with response status:", resp.Status)
}

func main() {
	r = rand.New(rand.NewSource(0))

	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: go run seed.go <n>")
		return
	}

	n := 0
	_, err := fmt.Sscanf(args[1], "%d", &n)
	if err != nil {
		fmt.Println("Invalid number of ads")
		return
	}

	// Generate and send N ads
	for i := 1; i <= n; i++ {
		ad := generateAd(i)
		jsonAd, _ := json.Marshal(ad)
		fmt.Println(string(jsonAd))
		sendAd(ad)
		log.Printf("Ad %d sent", i)
	}
}
