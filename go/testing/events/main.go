package main

import (
	"mf/events"
	"mf/models"
	"sync"
	"time"
)

func main() {
	// Test the events engine
	eventQueue := &events.EventQueue{}

	// Create a channel to handle early cancellation
	stop := make(chan struct{})

	// Create a wait group to wait for goroutines to finish
	var wg sync.WaitGroup

	// Process events
	go eventQueue.Dispatcher(&wg, stop)

	eventQueue.AddEvent(events.BattleStartEvent{
		Player: models.Player{
			Id:        1,
			Name:      "good_guy",
			CurrentHP: 100,
			MaxHP:     100,
			Attack:    100,
			Defense:   100,
			Speed:     1.0,
			Crit:      15.00,
			Dodge:     25.00,
			Block:     25.00,
			Ore:       10,
		},
		Enemy: models.Enemy{
			Id:        1,
			Name:      "bad_guy",
			CurrentHP: 100,
			MaxHP:     100,
			Attack:    100,
			Defense:   100,
			Speed:     1.0,
			Crit:      15.00,
			Dodge:     25.00,
			Block:     25.00,
		},
	})

	for {
		time.Sleep(1 * time.Second)
	}
}
