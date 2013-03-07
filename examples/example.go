//
// To run type:
// $ go run example.go
//
package main

import "github.com/heatxsink/statsd-go"

import (
	"os/exec"
	"bytes"
	"log"
	"strings"
	"strconv"
)

type temperature_sensors map[string] string

func get_computer_temperature() temperature_sensors {
	command := "/Applications/TemperatureMonitor.app/Contents/MacOS/tempmonitor"
	cmd := exec.Command(command, "-c", "-l", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	temp_data := temperature_sensors {}
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		datum := strings.Split(line, ": ")
		if len(datum) > 1 {
			key := datum[0]
			value := strings.TrimRight(datum[1], " C")
			temp_data[key] = value
		}
	}
	return temp_data
}

func main() {
	sensor_data :=get_computer_temperature()
	sensor_key := "SMC CPU A DIODE"
	sensor_value, _ := strconv.Atoi(sensor_data[sensor_key])
	hostname := "r1-4-dvm.rwc03.kabamdco.net"
	port_number := 9121
	client := statsd.New(hostname, port_number)
	client.Gauge("koc.test.smc_cpu_a_diode", sensor_value)
}
