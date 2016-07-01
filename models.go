package main

import (
	"time"
)

type FamilyMember struct {
	Name     string
	BornYear int
}

type DinnerStatus struct {
	Date   time.Time
	Coming bool
}
