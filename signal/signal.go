//go:build !windows

package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func Wait() os.Signal {
	return WaitSignal(
		syscall.SIGUSR1, syscall.SIGUSR2,
		syscall.SIGTTIN, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,
	)
}

func WaitSignal(s ...os.Signal) os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, s...)
	return <-c
}
