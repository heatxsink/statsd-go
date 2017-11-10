package statsd

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type Statsd struct {
	host       string
	port       int
	prefix     string
	connection net.Conn
}

func NewWithPrefix(host string, port int, prefix string) *Statsd {
	ss := &Statsd{host: host, port: port, prefix: prefix}
	ss.SetPrefix(prefix)
	return ss
}

func New(host string, port int) *Statsd {
	return &Statsd{host: host, port: port}
}

func (s *Statsd) Open() error {
	connectionString := fmt.Sprintf("%s:%d", s.host, s.port)
	conn, err := net.Dial("udp", connectionString)
	if err != nil {
		return err
	}
	s.connection = conn
	return nil
}

func (s *Statsd) Close() {
	s.connection.Close()
}

func (s *Statsd) SetPrefix(prefix string) {
	if strings.HasSuffix(prefix, ".") {
		prefix = strings.TrimSuffix(prefix, ".")
	}
	s.prefix = prefix
}

func (s *Statsd) Timing(stat string, time int64) error {
	updateString := fmt.Sprintf("%d|ms", time)
	stats := map[string]string{stat: updateString}
	return s.send(stats, 1)
}

func (s *Statsd) TimingWithSampleRate(stat string, time int64, sampleRate float32) error {
	updateString := fmt.Sprintf("%d|ms", time)
	stats := map[string]string{stat: updateString}
	return s.send(stats, sampleRate)
}

func (s *Statsd) Increment(stat string) error {
	stats := []string{stat}
	return s.updateStats(stats, 1, 1, "c")
}

func (s *Statsd) IncrementWithSampling(stat string, sampleRate float32) error {
	stats := []string{stat}
	return s.updateStats(stats[:], 1, sampleRate, "c")
}

func (s *Statsd) Decrement(stat string) error {
	stats := []string{stat}
	return s.updateStats(stats[:], -1, 1, "c")
}

func (s *Statsd) DecrementWithSampling(stat string, sampleRate float32) error {
	stats := []string{stat}
	return s.updateStats(stats[:], -1, sampleRate, "c")
}

func (s *Statsd) Counter(stat string, value int) error {
	stats := []string{stat}
	return s.updateStats(stats[:], 1, 1, "c")
}

func (s *Statsd) Gauge(stat string, value int) error {
	stats := []string{stat}
	return s.updateStats(stats[:], value, 1, "g")
}

func (s *Statsd) GaugeWithSampling(stat string, value int, sampleRate float32) error {
	stats := []string{stat}
	return s.updateStats(stats[:], value, sampleRate, "g")
}

func (s *Statsd) updateStats(stats []string, delta int, sampleRate float32, metric string) error {
	statsToSend := make(map[string]string)
	for _, stat := range stats {
		updateString := fmt.Sprintf("%d|%s", delta, metric)
		statsToSend[stat] = updateString
	}
	return s.send(statsToSend, sampleRate)
}

func (s *Statsd) injectPrefix(key string) string {
	if len(s.prefix) > 0 {
		return fmt.Sprintf("%s.%s", s.prefix, key)
	}
	return key
}

func (s *Statsd) send(data map[string]string, sampleRate float32) error {
	sampledData := make(map[string]string)
	if sampleRate < 1 {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		if rNum := r.Float32(); rNum <= sampleRate {
			for stat, value := range data {
				sampledUpdateString := fmt.Sprintf("%s|@%f", value, sampleRate)
				sampledData[stat] = sampledUpdateString
			}
		}
	} else {
		sampledData = data
	}

	for k, v := range sampledData {
		k = s.injectPrefix(k)
		updateString := fmt.Sprintf("%s:%s", k, v)
		_, err := fmt.Fprintf(s.connection, updateString)
		if err != nil {
			return err
		}
	}
	return nil
}
