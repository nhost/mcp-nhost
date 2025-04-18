package transport

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ThinkInAIXYZ/go-mcp/pkg"
)

type mockClientTransport struct {
	receiver clientReceiver
	in       io.ReadCloser
	out      io.Writer

	logger pkg.Logger

	cancel          context.CancelFunc
	receiveShutDone chan struct{}
}

func NewMockClientTransport(in io.ReadCloser, out io.Writer) ClientTransport {
	return &mockClientTransport{
		in:              in,
		out:             out,
		logger:          pkg.DefaultLogger,
		receiveShutDone: make(chan struct{}),
	}
}

func (t *mockClientTransport) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go func() {
		defer pkg.Recover()

		t.receive(ctx)

		close(t.receiveShutDone)
	}()

	return nil
}

func (t *mockClientTransport) Send(_ context.Context, msg Message) error {
	if _, err := t.out.Write(append(msg, mcpMessageDelimiter)); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

func (t *mockClientTransport) SetReceiver(receiver clientReceiver) {
	t.receiver = receiver
}

func (t *mockClientTransport) Close() error {
	t.cancel()

	if err := t.in.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	<-t.receiveShutDone

	return nil
}

func (t *mockClientTransport) receive(ctx context.Context) {
	s := bufio.NewScanner(t.in)

	for s.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			if err := t.receiver.Receive(ctx, s.Bytes()); err != nil {
				t.logger.Errorf("receiver failed: %v", err)
				return
			}
		}
	}

	if err := s.Err(); err != nil {
		if !errors.Is(err, io.ErrClosedPipe) { // This error occurs during unit tests, suppressing it here
			t.logger.Errorf("unexpected error reading input: %v", err)
		}
		return
	}
}
