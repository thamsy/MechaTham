package main

import (
	"time"
)

type FamilyMember struct {
	Name     string
	Id       int
	BornYear int
	Schedule bool
}

type DinnerStatus struct {
	Date   time.Time
	Coming bool
}
