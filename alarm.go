package main

import (
	"sort"
	"log"
)

type Alarm struct {
	thresholds []float64
}

func (a *Alarm) Check(current float64) bool {
	alarm := false
	for _, e := range a.thresholds {
		if e > current {
			alarm = true
			break
		}
	}

	if alarm {
		sort.Sort(ByValue(a.thresholds))
		log.Printf("Alarm triggered for %v\n", a.thresholds[0])
		a.thresholds = a.thresholds[1:]
	}

	return alarm
}

type ByValue []float64
func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool { return a[i] > a[j] }