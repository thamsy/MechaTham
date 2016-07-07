package main

import (
	"time"
)

type FamilyMember struct {
	Name            string
	Id              int
	BornYear        int
	DisableNotifTil time.Time
	PrevCommand     string
}

type DinnerStatus struct {
	Date   time.Time
	Coming bool
	Remark string
}
