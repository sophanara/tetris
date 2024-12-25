package signal

import (
	"fmt"
	"sync"
)

func Test() {
	nameSignal := NewSignal("so")
	nameChannel := nameSignal.Subscribe()

	go func() {
		for name := range nameChannel {
			fmt.Printf("New name %s \n", name)
		}
	}()

	// set new name
	nameSignal.Set("Bob")
	nameSignal.Set("Jason")
}

type Signal[T any] struct {
	value     T
	listeners []chan T
	mu        sync.RWMutex
}

// function to create new listener
func NewSignal[T any](initialValue T) *Signal[T] {
	return &Signal[T]{
		value: initialValue,
	}
}

// function to set the value in the listener
func (s *Signal[T]) Set(newValue T) {
	s.mu.Lock()
	s.value = newValue

	// copy listener to avoid the mutex contention
	currentListener := make([]chan T, len(s.listeners))
	copy(currentListener, s.listeners)
	s.mu.Unlock()

	// brodcast to each listener
	for _, listener := range currentListener {
		listener <- newValue
	}
}

func (s *Signal[T]) Subscribe() <-chan T {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan T, 1)
	s.listeners = append(s.listeners, ch)

	// send current value to new subscribe
	// like rsjs operator startWith
	ch <- s.value
	return ch
}
