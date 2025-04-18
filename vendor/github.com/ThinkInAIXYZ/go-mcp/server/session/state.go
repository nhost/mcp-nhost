package session

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"

	"github.com/ThinkInAIXYZ/go-mcp/pkg"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

type State struct {
	lastActiveAt time.Time

	mu       sync.RWMutex
	sendChan chan []byte

	requestID int64

	reqID2respChan cmap.ConcurrentMap[string, chan *protocol.JSONRPCResponse]

	// cache client initialize request info
	clientInfo         *protocol.Implementation
	clientCapabilities *protocol.ClientCapabilities

	// subscribed resources
	subscribedResources cmap.ConcurrentMap[string, struct{}]

	receivedInitRequest *pkg.AtomicBool
	ready               *pkg.AtomicBool
	closed              *pkg.AtomicBool
}

func NewState() *State {
	return &State{
		lastActiveAt:        time.Now(),
		sendChan:            make(chan []byte, 64),
		reqID2respChan:      cmap.New[chan *protocol.JSONRPCResponse](),
		subscribedResources: cmap.New[struct{}](),
		receivedInitRequest: pkg.NewAtomicBool(),
		ready:               pkg.NewAtomicBool(),
		closed:              pkg.NewAtomicBool(),
	}
}

func (s *State) SetClientInfo(ClientInfo *protocol.Implementation, ClientCapabilities *protocol.ClientCapabilities) {
	s.clientInfo = ClientInfo
	s.clientCapabilities = ClientCapabilities
}

func (s *State) SetReceivedInitRequest() {
	s.receivedInitRequest.Store(true)
}

func (s *State) GetReceivedInitRequest() bool {
	return s.receivedInitRequest.Load()
}

func (s *State) SetReady() {
	s.ready.Store(true)
}

func (s *State) GetReady() bool {
	return s.ready.Load()
}

func (s *State) IncRequestID() int64 {
	return atomic.AddInt64(&s.requestID, 1)
}

func (s *State) GetReqID2respChan() cmap.ConcurrentMap[string, chan *protocol.JSONRPCResponse] {
	return s.reqID2respChan
}

func (s *State) GetSubscribedResources() cmap.ConcurrentMap[string, struct{}] {
	return s.subscribedResources
}

func (s *State) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.closed.Store(true)
	close(s.sendChan)
}

func (s *State) updateLastActiveAt() {
	s.lastActiveAt = time.Now()
}

func (s *State) sendMessage(ctx context.Context, message []byte) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed.Load() {
		return errors.New("session already closed")
	}

	select {
	case s.sendChan <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *State) getMessageForSend(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg, ok := <-s.sendChan:
		if msg == nil && !ok {
			// There are no new messages and the chan has been closed, indicating that the request may need to be terminated.
			return nil, pkg.ErrSendEOF
		}
		return msg, nil
	}
}
