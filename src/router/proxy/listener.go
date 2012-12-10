package proxy

import (
	"net"
	"sync"
)

type TrackedConn struct {
	net.Conn

	listener *TrackingListener
	once     sync.Once
}

func (tc *TrackedConn) Close() error {
	tc.once.Do(func() { tc.listener.wg.Done() })
	return tc.Conn.Close()
}

type TrackingListener struct {
	net.Listener

	wg sync.WaitGroup
}

func NewTrackingListener(l net.Listener) *TrackingListener {
	return &TrackingListener{
		Listener: l,
	}
}

func (tl *TrackingListener) Accept() (net.Conn, error) {
	c, err := tl.Listener.Accept()
	if err != nil {
		return c, err
	}
	tl.wg.Add(1)
	tc := &TrackedConn{
		Conn:     c,
		listener: tl,
	}

	return tc, err
}

func (tl *TrackingListener) WaitForConnectionsClosed() {
	tl.wg.Wait()
}
