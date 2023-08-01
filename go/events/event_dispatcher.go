package events

// EventQueue holds incoming events
type EventQueue struct {
	events []Event
}

// AddEvent adds a new event to the queue
func (eq *EventQueue) AddEvent(event Event) {
	eq.events = append(eq.events, event)
}

// EventDispatcher processes events from the queue and calls the appropriate event handlers concurrently
func (eq *EventQueue) EventDispatcher() {
	// Create a channel to receive the processed events
	done := make(chan bool)

	// Process each event concurrently using goroutines
	for _, event := range eq.events {
		switch event.Type() {
		case EquipmentForging:
			go handleEquipmentForging(event.(EquipmentForgingEvent), done)
		case BattleStart:
			go handleBattleStart(event.(BattleStartEvent), done)
		case BattleConclude:
			go handleBattleConclude(event.(BattleConcludeEvent), done)
		case LevelUp:
			go handleLevelUp(event.(LevelUpEvent), done)
		case QuestCompletion:
			go handleQuestCompletion(event.(QuestCompletionEvent), done)
		}
	}

	// Wait for all goroutines to finish processing
	for range eq.events {
		<-done
	}

	// Clear processed events from the queue
	eq.events = nil
}
