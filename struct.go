package main

import (
	"time"
)

type AdConditions struct {
	AgeStart *int      `json:"ageStart,omitempty" binding:"omitempty,gte=0,lte=100"`
	AgeEnd   *int      `json:"ageEnd,omitempty" binding:"omitempty,gte=0,lte=100"`
	Gender   *[]string `json:"gender,omitempty" binding:"omitempty,unique,dive,oneof=M F"`
	Country  *[]string `json:"country,omitempty" binding:"omitempty,unique,dive,len=2,iso3166_1_alpha2"`
	Platform *[]string `json:"platform,omitempty" binding:"omitempty,unique,dive,oneof=android ios web"`
}

type Ad struct {
	Title      string       `json:"title" binding:"required,min=1,max=100"`
	StartAt    time.Time    `json:"startAt" binding:"required"`
	EndAt      time.Time    `json:"endAt" binding:"required"`
	Conditions AdConditions `json:"conditions,omitempty" binding:"omitempty"`
}

type AdQuery struct {
	Offset   int64  `form:"offset,default=0" binding:"gte=0"`
	Limit    int64  `form:"limit,default=5" binding:"gte=1,lte=100"`
	Age      int    `form:"age,default=0" binding:"gte=0,lte=100"`
	Gender   string `form:"gender" binding:"omitempty,oneof=M F"`
	Country  string `form:"country" binding:"omitempty,len=2,iso3166_1_alpha2"`
	Platform string `form:"platform" binding:"omitempty,oneof=android ios web"`
}
