package main

import (
	"mf/events"
	"mf/models"
)

func main() {
	// Test the events engine
	eventQueue := &events.EventQueue{}
	// Process events
	go eventQueue.EventDispatcher()

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
}
