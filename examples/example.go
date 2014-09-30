//
// To run type:
// $ go run example.go
//
package main

import (
	"github.com/heatxsink/go-osx-tempmonitor"
	"github.com/heatxsink/statsd-go"
	"strconv"
)

func main() {
	data, err := tempmonitor.GetTemperatureSensors()
	if err != nil {
		fmt.Println(err)
	}
	sensor_key := "SMC CPU A DIODE"
	sensor_value, _ := strconv.Atoi(data[sensor_key])
	hostname := "127.0.0.1"
	port_number := 9121
	client := statsd.New(hostname, port_number)
	client.Gauge("mbp.test.smc_cpu_a_diode", sensor_value)
}
