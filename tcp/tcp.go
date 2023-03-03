package tcp

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
	d, err := net.Dial("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer d.Close()

	return d.Write(packet)
}

func ListenAndServe(ctx context.Context, addr string, fn func(net.Conn, error)) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	logger.Info("tcp server listen on tcp://%s", l.Addr())

	w := sync.WaitGroup{}
	go func() {
		defer crash.Recover("crashed by tcp listen")

		for {
			conn, err := l.Accept()
			if err == io.EOF {
				logger.Info("tcp server EOF")
				return
			}

			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					logger.Info("tcp server closed")
					return
				}
			}

			w.Add(1)
			go func(c net.Conn, err error) {
				defer w.Done()
				defer crash.Recover("crashed by tcp accept")

				fn(c, err)
			}(conn, err)
		}
	}()

	<-ctx.Done()
	err = l.Close()
	w.Wait()

	return err
}
