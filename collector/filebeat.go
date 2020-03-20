package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//Filebeat json structure
type Filebeat struct {
	Events struct {
		Active float64 `json:"active"`
		Added  float64 `json:"added"`
		Done   float64 `json:"done"`
	} `json:"events"`

	Harvester struct {
		Closed    float64 `json:"closed"`
		OpenFiles float64 `json:"open_files"`
		Running   float64 `json:"running"`
		Skipped   float64 `json:"skipped"`
		Started   float64 `json:"started"`
		Files     map[string]struct {
			LastEventPublishedTime *FilebeatTime `json:"last_event_published_time,omitempty"`
			LastEventTimestamp     *FilebeatTime `json:"last_event_timestamp"`
			Name                   string        `json:"name"`
			ReadOffset             float64       `json:"read_offset"`
			Size                   float64       `json:"size"`
			StartTime              *time.Time    `json:"start_time"`
		} `json:"files"`
	} `json:"harvester"`

	Input struct {
		Log struct {
			Files struct {
				Renamed   float64 `json:"renamed"`
				Truncated float64 `json:"truncated"`
			} `json:"files"`
		} `json:"log"`
	} `json:"input"`
}

type filebeatCollector struct {
	beatInfo *BeatInfo
	stats    *Stats
	metrics  exportedMetrics
}

func (c *filebeatCollector) buildFileInfoMetrics() exportedMetrics {
	metrics := exportedMetrics{}

	paths := make(map[string]bool)
	for _, info := range c.stats.Filebeat.Harvester.Files {
		if _, seen := paths[info.Name]; seen {
			continue
		}

		paths[info.Name] = true
		metrics = append(metrics, []struct {
			desc    *prometheus.Desc
			eval    func(stats *Stats) float64
			valType prometheus.ValueType
		}{
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(c.beatInfo.Beat, "harvester", "size_bytes"),
					"filebeat.harvester.file.size",
					nil, prometheus.Labels{"path": info.Name},
				),
				eval:    func(stats *Stats) float64 { return info.Size },
				valType: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(c.beatInfo.Beat, "harvester", "last_event_published_time"),
					"filebeat.harvester.file.last_event_published_time",
					nil, prometheus.Labels{"path": info.Name},
				),
				eval: func(stats *Stats) float64 {
					if info.LastEventPublishedTime.Time != nil {
						return float64(info.LastEventPublishedTime.Time.UTC().Unix())
					}

					return 0
				},
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(c.beatInfo.Beat, "harvester", "last_event_timestamp"),
					"filebeat.harvester.file.last_event_timestamp",
					nil, prometheus.Labels{"path": info.Name},
				),
				eval: func(stats *Stats) float64 {
					if info.LastEventTimestamp.Time != nil {
						return float64(info.LastEventTimestamp.Time.UTC().Unix())
					}

					return 0
				},
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(c.beatInfo.Beat, "harvester", "start_time"),
					"filebeat.harvester.file.start_time",
					nil, prometheus.Labels{"path": info.Name},
				),
				eval:    func(stats *Stats) float64 { return float64(info.StartTime.UTC().Unix()) },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(c.beatInfo.Beat, "harvester", "file_read_offset"),
					"filebeat.harvester.file.file_read_offset",
					nil, prometheus.Labels{"path": info.Name},
				),
				eval:    func(stats *Stats) float64 { return info.ReadOffset },
				valType: prometheus.GaugeValue,
			},
		}...)
	}

	return metrics
}

// NewFilebeatCollector constructor
func NewFilebeatCollector(beatInfo *BeatInfo, stats *Stats) prometheus.Collector {
	c := &filebeatCollector{
		beatInfo: beatInfo,
		stats:    stats,
		metrics: exportedMetrics{
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "events"),
					"filebeat.events",
					nil, prometheus.Labels{"event": "active"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Events.Active },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "events"),
					"filebeat.events",
					nil, prometheus.Labels{"event": "added"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Events.Added },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "events"),
					"filebeat.events",
					nil, prometheus.Labels{"event": "done"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Events.Done },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "harvester"),
					"filebeat.harvester",
					nil, prometheus.Labels{"harvester": "closed"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.Closed },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "harvester"),
					"filebeat.harvester",
					nil, prometheus.Labels{"harvester": "open_files"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.OpenFiles },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "harvester"),
					"filebeat.harvester",
					nil, prometheus.Labels{"harvester": "running"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.Running },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "harvester"),
					"filebeat.harvester",
					nil, prometheus.Labels{"harvester": "skipped"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.Skipped },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "harvester"),
					"filebeat.harvester",
					nil, prometheus.Labels{"harvester": "started"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.Started },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "input_log"),
					"filebeat.input_log",
					nil, prometheus.Labels{"files": "renamed"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Input.Log.Files.Renamed },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "filebeat", "input_log"),
					"filebeat.input_log",
					nil, prometheus.Labels{"files": "truncated"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Input.Log.Files.Truncated },
				valType: prometheus.UntypedValue,
			},
		},
	}

	return c
}

// Describe returns all descriptions of the collector.
func (c *filebeatCollector) Describe(ch chan<- *prometheus.Desc) {

	for _, metric := range c.metrics {
		ch <- metric.desc
	}

}

// Collect returns the current state of all metrics of the collector.
func (c *filebeatCollector) Collect(ch chan<- prometheus.Metric) {
	for _, f := range c.buildFileInfoMetrics() {
		ch <- prometheus.MustNewConstMetric(f.desc, f.valType, f.eval(c.stats))
	}

	for _, i := range c.metrics {
		ch <- prometheus.MustNewConstMetric(i.desc, i.valType, i.eval(c.stats))
	}

}
