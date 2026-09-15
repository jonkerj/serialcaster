package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jonkerj/serialcaster/internal/distributor"
)

var rootCmd = &cobra.Command{
	Use:   "serialcaster",
	Short: "Forward serial port RX bytes to all connected TCP clients",
	RunE:  run,
}

// Execute is the entry point called by main. It runs the root command and
// exits with a non-zero status on failure.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	f := rootCmd.PersistentFlags()
	f.StringP("port", "p", "/dev/ttyUSB0", "serial device path")
	f.IntP("baud", "b", 9600, "baud rate")
	f.Int("databits", 8, "data bits (5, 6, 7, 8)")
	f.String("stopbits", "1", "stop bits (1, 1.5, 2)")
	f.String("parity", "none", "parity mode (none, odd, even, mark, space)")
	f.StringP("listen", "l", ":8888", "TCP listen address")

	// SD_PORT, SD_BAUD, SD_DATABITS, SD_STOPBITS, SD_PARITY, SD_LISTEN
	viper.SetEnvPrefix("SD")
	viper.AutomaticEnv()
	viper.BindPFlags(f) //nolint:errcheck
}

func run(_ *cobra.Command, _ []string) error {
	mode, err := distributor.BuildMode(
		viper.GetInt("baud"),
		viper.GetInt("databits"),
		viper.GetString("stopbits"),
		viper.GetString("parity"),
	)
	if err != nil {
		slog.Error("invalid serial config", "err", err)
		os.Exit(1)
	}

	h := distributor.NewHub()
	go h.Run()
	go distributor.ReadSerial(h, viper.GetString("port"), mode)
	distributor.ServeTCP(h, viper.GetString("listen"))
	return nil
}
