package distributor

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"go.bug.st/serial"
)

// BuildMode constructs a serial.Mode from the given parameters, returning an
// error if any value is outside the set accepted by the serial library.
func BuildMode(baud, dataBits int, stopBits, parity string) (*serial.Mode, error) {
	var p serial.Parity
	switch strings.ToLower(parity) {
	case "none", "n":
		p = serial.NoParity
	case "odd", "o":
		p = serial.OddParity
	case "even", "e":
		p = serial.EvenParity
	case "mark", "m":
		p = serial.MarkParity
	case "space", "s":
		p = serial.SpaceParity
	default:
		return nil, fmt.Errorf("unknown parity %q", parity)
	}

	var sb serial.StopBits
	switch stopBits {
	case "1":
		sb = serial.OneStopBit
	case "1.5":
		sb = serial.OnePointFiveStopBits
	case "2":
		sb = serial.TwoStopBits
	default:
		return nil, fmt.Errorf("unknown stop bits %q", stopBits)
	}

	switch dataBits {
	case 5, 6, 7, 8:
	default:
		return nil, fmt.Errorf("invalid data bits %d: must be 5, 6, 7, or 8", dataBits)
	}

	return &serial.Mode{
		BaudRate: baud,
		DataBits: dataBits,
		Parity:   p,
		StopBits: sb,
	}, nil
}

// ReadSerial opens port, reads incoming bytes, and broadcasts each frame to h.
// It exits the process on open or read failure.
func ReadSerial(h *Hub, port string, mode *serial.Mode) {
	slog.Info("opening serial port", "port", port, "baud", mode.BaudRate)
	sp, err := serial.Open(port, mode)
	if err != nil {
		slog.Error("serial open failed", "port", port, "err", err)
		os.Exit(1)
	}
	defer sp.Close() //nolint:errcheck

	buf := make([]byte, 4096)
	for {
		n, err := sp.Read(buf)
		if n > 0 {
			frame := make([]byte, n)
			copy(frame, buf[:n])
			h.Broadcast(frame)
		}
		if err != nil {
			slog.Error("serial read failed", "err", err)
			os.Exit(1)
		}
	}
}
