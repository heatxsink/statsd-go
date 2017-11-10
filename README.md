statsd-go
=========
A statsd golang package. What is the difference between this version and others? This one has gauge support.

Setup
-----
Get the following golang packages.

	
	$ go get github.com/heatxsink/statsd-go

Run the example
---------------
1. Install [Temperature Monitor](http://www.bresink.com/osx/TemperatureMonitor.html) to get temp data from OS X.
1. Then ...

	$ go get github.com/heatxsink/go-osx-tempmonitor

1. Then ...

	$ cd examples
	$ go run example.go


Example
-------
```go
import(
	"github.com/heatxsink/statsd-go"
)

func main() {
	ss := statsd.New("127.0.0.1", 9121)
	ss.SetPrefix("omg")
	ss.Gauge("mbp.test.smc_cpu_a_diode", 75)
}
```
