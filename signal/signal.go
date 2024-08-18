package signal

import (
	"context"
	"os"
	"os/signal"
)

func ContextClosableOnSignals(signals ...os.Signal) context.Context {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, signals...)

	ctx, cancel := context.WithCancel(context.Background())

	go func(signalCh <-chan os.Signal, cancel context.CancelFunc) {
		<-signalCh

		cancel()
	}(signalCh, cancel)

	return ctx
}
