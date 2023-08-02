package events

import (
	"fmt"
	"sync"
)

// EventQueue holds incoming events
type EventQueue struct {
	events []Event
	sync.Mutex
}

// NewEventQueue creates a new event queue
func NewEventQueue() *EventQueue {
	return &EventQueue{}
}

// AddEvent adds a new event to the queue
func (eq *EventQueue) AddEvent(event Event) {
	eq.Lock()
	eq.events = append(eq.events, event)
	eq.Unlock()
}

// Dispatcher processes events from the channel and calls the appropriate event handlers concurrently
func (eq *EventQueue) Dispatcher(
	wg *sync.WaitGroup,
	stop <-chan struct{},
) {
	fmt.Println("Processing messages")
	// Create a channel to receive the processed events
	done := make(chan bool)

	// Start the event processing in a separate goroutine
	go func() {
		for {
			var event Event
			eq.Lock()
			if len(eq.events) > 0 {
				event = eq.events[0]
				eq.events = eq.events[1:]
			}
			eq.Unlock()

			if event == nil {
				select {
				case <-stop:
					done <- true
					return
				}
			}

			switch event.Type() {
			case EquipmentForging:
				fmt.Println("received Forge event")
				wg.Add(1)
				go handleEquipmentForging(event.(EquipmentForgingEvent), wg)
			case BattleStart:
				fmt.Println("received battle start event")
				wg.Add(1)
				go handleBattleStart(event.(BattleStartEvent), wg)
			case BattleConclude:
				fmt.Println("received battle conclude event")
				wg.Add(1)
				go handleBattleConclude(event.(BattleConcludeEvent), wg)

			}
		}
	}()

	// Wait for the dispatcher to finish before exiting the EventQueue
	<-done
	fmt.Println("Done processing messages")
}
