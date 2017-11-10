//
// To run type:
// $ go run example.go
//
package main

import (
	"fmt"
	"log"
	"strconv"

	tempmonitor "github.com/heatxsink/go-osx-tempmonitor"
	statsd "github.com/heatxsink/statsd-go"
)

func main() {
	data, err := tempmonitor.GetTemperatureSensors()
	if err != nil {
		fmt.Println(err)
	}
	sensorKey := "SMC CPU A DIODE"
	sensorValue, _ := strconv.Atoi(data[sensorKey])
	ss := statsd.New("127.0.0.1", 9125)
	err := ss.Open()
	if err != nil {
		log.Info(err)
	}
	err := ss.Gauge("mbp.test.smc_cpu_a_diode", sensorValue)
	if err != nil {
		log.Info(err)
	}
}
