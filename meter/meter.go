package meter

import (
	"os"
	"syscall"
	"unsafe"

	"log"
	"math"

	"sync/atomic"

	"github.com/pkg/errors"
)

const (
	MeterTemp       byte    = 0x42
	MeterCO2        byte    = 0x50
	hidiocsfeature9 uintptr = 0xc0094806
)

var key = [8]byte{}

// Meter gives access to the CO2 Meter. Make sure to call Open before Read.
type Meter struct {
	file   *os.File
	opened int32
}

// Measurement is the result of a Read operation.
type Measurement struct {
	Temperature float64
	Co2		 int
}

// Open will open the device file specified in the path which is usually something like /dev/hidraw2.
func (m *Meter) Open(path string) (err error) {
	atomic.StoreInt32(&m.opened, 1)

	m.file, err = os.OpenFile(path, os.O_RDWR, 0644)

	if err != nil || m.file == nil {
		return errors.Wrapf(err, "Failed to open '%v'", path)
	}

	log.Printf("Device '%v' opened", m.file.Name())
	return m.ioctl()
}

// ioctl writes into the device file. We need to write 9 bytes where the first byte specifies the report number.
// In this case 0x00.
func (m *Meter) ioctl() error {
	data := [9]byte{}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, m.file.Fd(), hidiocsfeature9, uintptr(unsafe.Pointer(&data)))

	if err != 0 {
		return errors.Wrap(syscall.Errno(err), "ioctl failed")
	}
	return nil
}

// Read will read one record from the device
func (m *Meter) ReadOne() (byte, int, error) {
	if atomic.LoadInt32(&m.opened) != 1 {
		return 0, 0, errors.New("Device needs to be opened")
	}
	result := make([]byte, 8)

	_, err := m.file.Read(result)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "Could not read from: '%v'", m.file.Name())
	}

	operation := result[0]
	value := int(result[1])<<8 | int(result[2])
	// result[3] is a checksum
	// result[4] is 0x0D
	// result[5:7] are 0x00

	return operation, value, nil
}

func ConvertTemp (value int) (float64) {
	return math.Round((float64(value)/16.0 - 273.15) * 10.0) / 10.0
}

// Read will read from the device file until it finds a temperature and co2 measurement. Before it can be used the
// device file needs to be opened via Open.
func (m *Meter) Read() (*Measurement, error) {
	measurement := &Measurement{Co2: 0, Temperature: 0}

	for {
		operation, value, err := m.ReadOne()
		if err != nil {
			return nil, errors.Wrapf(err, "Could not read from: '%v'", m.file.Name())
		}


		switch byte(operation) {
		case MeterCO2:
			measurement.Co2 = int(value)
		case MeterTemp:
			measurement.Temperature = ConvertTemp(value)
		}

		if measurement.Co2 != 0 && measurement.Temperature != 0 {
			return measurement, nil
		}
	}
}

// Close will close the device file.
func (m *Meter) Close() error {
	log.Printf("Closing '%v'", m.file.Name())
	atomic.StoreInt32(&m.opened, 0)
	return m.file.Close()
}

