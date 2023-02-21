package signal

import (
	"os"
)

func Wait() os.Signal {
	<-make(chan any, 0)
	return nil
}

func WaitSignal(s ...os.Signal) os.Signal {
	<-make(chan any, 0)
	return nil
}
