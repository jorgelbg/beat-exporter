package collector

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//BeatInfo beat info json structure
type BeatInfo struct {
	Beat     string `json:"beat"`
	Hostname string `json:"hostname"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Version  string `json:"version"`
}

//Stats stats endpoint json structure
type Stats struct {
	System     System      `json:"system"`
	Beat       BeatStats   `json:"beat"`
	LibBeat    LibBeat     `json:"libbeat"`
	Registrar  Registrar   `json:"registrar"`
	Filebeat   Filebeat    `json:"filebeat"`
	Metricbeat Metricbeat  `json:"metricbeat"`
	Auditd     AuditdStats `json:"auditd"`
}

// Metric represents a prometheus metric (descriptor, evaluator, type)
type Metric struct {
	desc    *prometheus.Desc
	eval    func(stats *Stats) float64
	valType prometheus.ValueType
}

type exportedMetrics []struct {
	desc    *prometheus.Desc
	eval    func(stats *Stats) float64
	valType prometheus.ValueType
}

// FilebeatTime is a custom time struct to deal with the particularity of the
// format time exposed by filebeat
type FilebeatTime struct {
	*time.Time
}

// UnmarshalJSON custom logic that deals with an empty string or null
// values in filebeat time.Time fields
func (f *FilebeatTime) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	if str == "" {
		f.Time = nil
		return nil
	}

	t, err := time.Parse(time.RFC3339, str)
	f.Time = &t
	return err
}
