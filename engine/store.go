package main

import (
	"sync"
	"time"
)

type ReconcileEvent struct {
	Timestamp time.Time `json:"timestamp"`
	ServiceName string `json:"service_name"`
	WhatChanged string `json:"what_changed"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

type Store struct {
	mu sync.Mutex 
	Events []ReconcileEvent `json:"events"`
}

func (s *Store) AddEvent(event ReconcileEvent) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Events = append(s.Events, event)
	if len(s.Events) > 20 {
    s.Events = s.Events[len(s.Events)-20:]
	}
}

func (s *Store) GetEvents() []ReconcileEvent {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.Events
}