# CO2Mini (RAD-0301)

I recently bought a CO2Mini (RAD-0301) CO₂ monitor for my apartment, with the intent of logging the data over time.

A quick Google search brought me to [several](https://github.com/vfilimonov/co2meter) [projects](https://github.com/dmage/co2mon) that interface with this model, including a decryption routine.

Once I got my hands on the monitor, it turns out that my model appears to be unencrypted? None of the projects I've surveyed make the encryption optional, so I'm not sure if this is a distinct product, or maybe something has changed? 

The [product website](https://www.co2meter.com/collections/indoor-air-quality/products/co2mini-co2-indoor-air-quality-monitor) even has a [protocol description](http://co2meters.com/Documentation/Other/AN_RAD_0301_USB_Communications_Revised8.pdf), but it doesn't mention anything about encryption. 

I've modified [larsp's project](https://github.com/larsp/co2monitor), which is a Go client to read the monitor and expose the values over a Prometheus endpoint, to remove the decryption routines.

Looking at the data, there also appear to be 9 undocumented ops that my device sends. I've also added logging for them, so I can see how they correlate with time, CO₂, and temperature to try to reverse-engineer them. I unfortunately had to break the nice abstraction that larsp's library made (and most libraries seem to make) of returning a `(CO₂, temperature)`` pair: some of these ops seem to come in at different rates? 

I'll try to keep this updated if I figure out what these other ops are for.

----

# Original Readme

[![Go Report Card](https://goreportcard.com/badge/github.com/zbanks/co2monitor)](https://goreportcard.com/report/github.com/zbanks/co2monitor)
[![GoDoc](https://godoc.org/github.com/zbanks/co2monitor/meter?status.svg)](https://godoc.org/github.com/zbanks/co2monitor/meter)

# CO₂ monitor

## Setup & Example
<img src="https://raw.githubusercontent.com/zbanks/co2monitor/img/monitor.jpg" alt="Setup" width="700">
<img src="https://raw.githubusercontent.com/zbanks/co2monitor/img/dashboard.png" alt="Dashboard" width="700">

## Motivation
Some time ago an [article](https://blog.wooga.com/woogas-office-weather-wow-67e24a5338) about a low cost CO₂ monitor 
came to our attention. A colleague quickly adopted the python [code](https://github.com/wooga/office_weather)
to fit in our prometheus setup. Since humans are sensitive to temperature and CO₂ level, we were now able to 
optimize HVAC settings in our office (Well, we mainly complained to our facility management).

For numerous reasons I wanted to replace the python code with a static Go binary.

## Hardware
- CO₂ meter: Can be found for around 70EUR/USD at [amazon.com](https://www.amazon.com/dp/B00H7HFINS) 
& [amazon.de](https://www.amazon.de/dp/B00TH3OW4Q/). Regardless of minor differences between both devices, both work.
- Some machine which can run the compiled Go binary, has USB and is reachable from your prometheus collector. 
A very first version of a raspberry pi is already sufficient.

## Software
You need prometheus to collect the metrics.

It might make things easier when you set up an `udev` rule e.g.
```bash
$ cat /etc/udev/rules.d/99-hidraw-permissions.rules 
KERNEL=="hidraw*", SUBSYSTEM=="hidraw", MODE="0664", GROUP="plugdev"
```

## Run & Collect

Help
```bash
$ ./co2monitor --help      
usage: co2monitor [<flags>] <device> [<listen-address>]

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Args:
  <device>            CO2 Meter device, such as /dev/hidraw2
  [<listen-address>]  The address to listen on for HTTP requests.
```

Starting the meter export
```bash
$ ./co2monitor /dev/hidraw2
2018/01/18 13:09:31 Serving metrics at ':8080/metrics'
2018/01/18 13:09:31 Device '/dev/hidraw2' opened

```

## Credit

[Henryk Plötz](https://hackaday.io/project/5301-reverse-engineering-a-low-cost-usb-co-monitor/log/17909-all-your-base-are-belong-to-us)
& [wooga](https://github.com/wooga/office_weather)
