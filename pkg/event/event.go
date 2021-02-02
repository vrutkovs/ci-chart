package event

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

type event struct {
	value       string
	description string
	timestamp   time.Time
}

type store struct {
	sync.Mutex
	events map[string]map[string][]event
}

type Store interface {
	Add(group, label, value, extended string)
	JSONHandler(w http.ResponseWriter, r *http.Request)
}

func NewStore() Store {
	return &store{events: map[string]map[string][]event{}}
}

func (s *store) Add(group, label, value, description string) {
	s.Lock()
	defer s.Unlock()
	groupevents, ok := s.events[group]
	if !ok {
		s.events[group] = map[string][]event{}
		groupevents, _ = s.events[group]
	}
	labelevents, ok := groupevents[label]
	if !ok || labelevents[len(labelevents)-1].value != value {
		event := event{value: value, description: description, timestamp: time.Now()}
		klog.Infof("adding event for %s/%s: %#v", group, label, event)
		groupevents[label] = append(groupevents[label], event)
	} else {
		klog.Infof("duplicate event dropped for %s/%s", group, label)
	}
}

type LabelData struct {
	TimeRange [2]time.Time `json:"timeRange,omitempty"`
	Val       string       `json:"val,omitempty"`
	Extended  string       `json:"extended,omitempty"`
}

type GroupData struct {
	Label string      `json:"label,omitempty"`
	Data  []LabelData `json:"data,omitempty"`
}

type Group struct {
	Group string      `json:"group,omitempty"`
	Data  []GroupData `json:"data,omitempty"`
}

func (s *store) JSONHandler(w http.ResponseWriter, r *http.Request) {
	var groups []Group
	for group, labels := range s.events {
		g := Group{Group: group}
		for label, events := range labels {
			gd := GroupData{Label: label}
			for i, event := range events {
				ld := LabelData{TimeRange: [2]time.Time{event.timestamp, time.Now()}, Val: event.value, Extended: event.description}
				gd.Data = append(gd.Data, ld)
				if i > 0 {
					gd.Data[i-1].TimeRange[1] = event.timestamp
				}
			}
			g.Data = append(g.Data, gd)
		}
		groups = append(groups, g)
	}
	js, err := json.Marshal(groups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
