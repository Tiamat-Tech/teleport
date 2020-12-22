package api

import (
	"context"
	"sync"

	"github.com/gravitational/teleport/api/proto"
	"github.com/gravitational/teleport/api/types"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gravitational/trace/trail"
)

// NewKeepAliver returns a new instance of keep aliver.
// Run k.Close to release the keepAliver and its goroutines
func (c *Client) NewKeepAliver(ctx context.Context) (types.KeepAliver, error) {
	cancelCtx, cancel := context.WithCancel(ctx)
	stream, err := c.grpc.SendKeepAlives(cancelCtx)
	if err != nil {
		cancel()
		return nil, trail.FromGRPC(err)
	}
	k := &streamKeepAliver{
		stream:      stream,
		ctx:         cancelCtx,
		cancel:      cancel,
		keepAlivesC: make(chan types.KeepAlive),
	}
	go k.forwardKeepAlives()
	go k.recv()
	return k, nil
}

type streamKeepAliver struct {
	mu          sync.RWMutex
	stream      proto.AuthService_SendKeepAlivesClient
	ctx         context.Context
	cancel      context.CancelFunc
	keepAlivesC chan types.KeepAlive
	err         error
}

// KeepAlives returns the streamKeepAliver's channel of KeepAlives
func (k *streamKeepAliver) KeepAlives() chan<- types.KeepAlive {
	return k.keepAlivesC
}

func (k *streamKeepAliver) forwardKeepAlives() {
	for {
		select {
		case <-k.ctx.Done():
			return
		case keepAlive := <-k.keepAlivesC:
			err := k.stream.Send(&keepAlive)
			if err != nil {
				k.closeWithError(trail.FromGRPC(err))
				return
			}
		}
	}
}

// Error returns the streamKeepAliver's error after closing
func (k *streamKeepAliver) Error() error {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.err
}

// Done returns a channel that closes once the streamKeepAliver is Closed
func (k *streamKeepAliver) Done() <-chan struct{} {
	return k.ctx.Done()
}

// recv is necessary to receive errors from the
// server, otherwise no errors will be propagated
func (k *streamKeepAliver) recv() {
	err := k.stream.RecvMsg(&empty.Empty{})
	k.closeWithError(trail.FromGRPC(err))
}

func (k *streamKeepAliver) closeWithError(err error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.Close()
	k.err = err
}

// Close the streamKeepAliver
func (k *streamKeepAliver) Close() error {
	k.cancel()
	return nil
}
