package signals

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

var onlyOne = make(chan struct{})

func SetupSignalHandler() <-chan struct{} {
	// panics when called twice
	close(onlyOne)
	stopCh := make(chan struct{})
	ch := make(chan os.Signal, 2)
	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM}
	signal.Notify(ch, shutdownSignals...)
	go func() {
		s := <-ch
		glog.Infof("receive signal: %s\n", s.String())
		glog.Infoln("start close stopCh")
		close(stopCh)
		s = <-ch
		glog.Infof("receive signal: %s again. start exit.\n", s.String())
	}()
	return stopCh
}
