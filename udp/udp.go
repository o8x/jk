package udp

import (
	"io"
	"net"
	"strings"
	"sync"

	"golang.org/x/net/context"

	"github.com/o8x/jk/v2/crash"
	"github.com/o8x/jk/v2/logger"
)

func WritePacket(addr string, packet []byte) (int, error) {
	d, err := net.Dial("udp", addr)
	if err != nil {
		return 0, err
	}

	defer d.Close()

	return d.Write(packet)
}

func ListenAndServe(ctx context.Context, addr string, fn func([]byte, net.Addr, error)) error {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}

	logger.Info("udp server listen on udp://%s", l.LocalAddr())

	w := sync.WaitGroup{}
	go func() {
		defer crash.Recover("crashed by udp listen")

		for {
			bs := make([]byte, 8192)

			n, addr, err := l.ReadFrom(bs)
			if err == io.EOF {
				logger.Info("udp server EOF")
				return
			}

			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					logger.Info("udp server closed")
					return
				}
			}

			w.Add(1)
			go func(data []byte, addr net.Addr, err error) {
				defer w.Done()
				defer crash.Recover("crashed by udp accept")

				fn(data, addr, err)
			}(bs[:n], addr, err)
		}
	}()

	<-ctx.Done()
	err = l.Close()
	w.Wait()

	return err
}
