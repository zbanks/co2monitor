package main

import (
	"log"
	"net/http"

	"github.com/zbanks/co2monitor/meter"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	device     = kingpin.Arg("device", "CO2 Meter device, such as /dev/hidraw2").Required().String()
	listenAddr = kingpin.Arg("listen-address", "The address to listen on for HTTP requests.").Default(":8080").String()
)

var (
	temperature = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_temperature_celsius",
		Help: "Current temperature in Celsius",
	})
	co2 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_co2_ppm",
		Help: "Current CO2 level (ppm)",
	})
	unknown41 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_unknown_41",
		Help: "Unknown value 0x41 (u16)",
	})
	unknown43 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_unknown_43",
		Help: "Unknown value 0x43 (u16)",
	})
	unknown4f = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_unknown_4f",
		Help: "Unknown value 0x4f (u16)",
	})
	unknown52 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_unknown_52",
		Help: "Unknown value 0x52 (u16)",
	})
	unknown56 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_unknown_56",
		Help: "Unknown value 0x56 (u16)",
	})
	unknown57 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_unknown_57",
		Help: "Unknown value 0x57 (u16)",
	})
	unknown6d = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_unknown_6d",
		Help: "Unknown value 0x6d (u16)",
	})
	unknown6e = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_unknown_6e",
		Help: "Unknown value 0x6e (u16)",
	})
	unknown71 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_unknown_71",
		Help: "Unknown value 0x71 (u16)",
	})
)

func init() {
	prometheus.MustRegister(temperature)
	prometheus.MustRegister(co2)
	prometheus.MustRegister(unknown41)
	prometheus.MustRegister(unknown43)
	prometheus.MustRegister(unknown4f)
	prometheus.MustRegister(unknown52)
	prometheus.MustRegister(unknown56)
	prometheus.MustRegister(unknown57)
	prometheus.MustRegister(unknown6d)
	prometheus.MustRegister(unknown6e)
	prometheus.MustRegister(unknown71)
}

func main() {
	kingpin.Parse()
	http.Handle("/metrics", promhttp.Handler())
	go measure()
	log.Printf("Serving metrics at '%v/metrics'", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func measure() {
	m := new(meter.Meter)
	err := m.Open(*device)
	if err != nil {
		log.Fatalf("Could not open '%v'", *device)
		return
	}

	for {
		op, value, err := m.ReadOne()
		if err != nil {
			log.Fatalf("Something went wrong: '%v'", err)
		}
		switch op {
		case meter.MeterTemp:
			temperature.Set(meter.ConvertTemp(value))
		case meter.MeterCO2:
			co2.Set(float64(value))
		case 0x41:
			unknown41.Set(float64(value))
		case 0x43:
			unknown43.Set(float64(value))
		case 0x4f:
			unknown4f.Set(float64(value))
		case 0x52:
			unknown52.Set(float64(value))
		case 0x56:
			unknown56.Set(float64(value))
		case 0x57:
			unknown57.Set(float64(value))
		case 0x6d:
			unknown6d.Set(float64(value))
		case 0x6e:
			unknown6e.Set(float64(value))
		case 0x71:
			unknown71.Set(float64(value))
		default:
			log.Printf("Unknown operation: %x", op)
		}
	}
}
