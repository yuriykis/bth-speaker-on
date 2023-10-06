package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yuriykis/bth-speaker-on/log"

	"github.com/yuriykis/bth-speaker-on/system"
)

func printStartupInfo() {
	fmt.Println(asciBanner)
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Arch: %s\n", runtime.GOARCH)
	fmt.Printf("Version: %s\n", runtime.Version())
	fmt.Println("Log file: /tmp/bth-speaker-on.log")
}

var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"start"},
	Short:   "Start bth-speaker-on",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.ClearLogFile()
		log.Println(asciBanner)
		printStartupInfo()

		upIntervalFlag := cmd.Flags().Lookup("up-interval")
		if upIntervalFlag == nil {
			log.Fatal("up-interval flag is not set")
		}
		upIntervalValue, err := strconv.Atoi(upIntervalFlag.Value.String())
		if err != nil {
			log.Fatal(err)
		}
		upInterval := time.Duration(upIntervalValue) * time.Second // TODO: minutes
		flag.Parse()

		Run(upInterval)
	},
}

var rootCmd = &cobra.Command{
	Use:     "bth-speaker-on",
	Aliases: []string{"bth-speaker-on"},
	Short:   "bth-speaker-on is a tool to keep bluetooth speaker on",
	Long:    "bth-speaker-on is a tool to keep bluetooth speaker on",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	startCmd.Flags().IntP(
		"up-interval",
		"u",
		5,
		"Interval in minutes to check if device is up",
	)
	rootCmd.AddCommand(startCmd)
}

func Run(upInterval time.Duration) {
	var (
		dm  system.DeviceManager
		err error
	)
	switch system.SystemType(runtime.GOOS) {
	case system.MacSystemType:
		dm, err = system.NewMacDeviceManager()
	case system.LinuxSystemType:
		dm, err = system.NewLinuxDeviceManager()
	case system.WindowsSystemType:
		dm, err = system.NewWindowsDeviceManager()
	default:
		log.Fatal("Unknown system type")
	}

	dm = system.NewLoggingDeviceManagerMiddleware(dm)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	dm.Start(ctx, upInterval)
}

var asciBanner = `
  ___  _____  _  _   ___                   _                ___   _  _ 
 | _ )|_   _|| || | / __| _ __  ___  __ _ | |__ ___  _ _   / _ \ | \| |
 | _ \  | |  | __ | \__ \| '_ \/ -_)/ _' || / // -_)| '_| | (_) || .' |
 |___/  |_|  |_||_| |___/| .__/\___|\__,_||_\_\\___||_|    \___/ |_|\_|
                         |_|                                           
`
