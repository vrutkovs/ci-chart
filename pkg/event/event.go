package event

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

type Input struct {
	group       string
	label       string
	value       string
	description string
	timestamp   time.Time
}

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
	Add(i Input)
	JSONHandler(w http.ResponseWriter, r *http.Request)
}

func NewStore() Store {
	return &store{events: map[string]map[string][]event{}}
}

func (s *store) Add(i Input) {
	s.Lock()
	defer s.Unlock()
	groupevents, ok := s.events[i.group]
	if !ok {
		s.events[i.group] = map[string][]event{}
		groupevents, _ = s.events[i.group]
	}
	labelevents, ok := groupevents[i.label]
	if !ok || labelevents[len(labelevents)-1].value != i.value {
		event := event{value: i.value, description: i.description, timestamp: i.timestamp}
		klog.Infof("adding event for %s/%s: %#v", i.group, i.label, event)
		groupevents[i.label] = append(groupevents[i.label], event)
	} else {
		klog.Infof("duplicate event dropped for %s/%s", i.group, i.label)
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

func PodEventToInput(ev corev1.Event) *Input {
	if ev.EventTime.Time.IsZero() {
		return nil
	}
	return &Input{
		group:     ev.InvolvedObject.Namespace,
		label:     ev.InvolvedObject.Name,
		value:     ev.Reason,
		timestamp: ev.EventTime.Time,
	}
}

func ClusterOperatorEventToInput(ev corev1.Event) *Input {
	if ev.EventTime.Time.IsZero() {
		return nil
	}
	return &Input{
		group:     ev.InvolvedObject.Namespace,
		label:     ev.InvolvedObject.Name,
		value:     ev.Reason,
		timestamp: ev.EventTime.Time,
	}
}
