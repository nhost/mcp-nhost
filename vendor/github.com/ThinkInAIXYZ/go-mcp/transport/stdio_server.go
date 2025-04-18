package transport

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ThinkInAIXYZ/go-mcp/pkg"
)

const stdioSessionID = "stdio"

type StdioServerTransportOption func(*stdioServerTransport)

func WithStdioServerOptionLogger(log pkg.Logger) StdioServerTransportOption {
	return func(t *stdioServerTransport) {
		t.logger = log
	}
}

type stdioServerTransport struct {
	receiver serverReceiver
	reader   io.ReadCloser
	writer   io.Writer

	sessionManager sessionManager

	logger pkg.Logger

	cancel          context.CancelFunc
	receiveShutDone chan struct{}
}

func NewStdioServerTransport(opts ...StdioServerTransportOption) ServerTransport {
	t := &stdioServerTransport{
		reader: os.Stdin,
		writer: os.Stdout,
		logger: pkg.DefaultLogger,

		receiveShutDone: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *stdioServerTransport) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	t.sessionManager.CreateSession(stdioSessionID)

	t.receive(ctx)

	close(t.receiveShutDone)
	return nil
}

func (t *stdioServerTransport) Send(_ context.Context, _ string, msg Message) error {
	if _, err := t.writer.Write(append(msg, mcpMessageDelimiter)); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

func (t *stdioServerTransport) SetReceiver(receiver serverReceiver) {
	t.receiver = receiver
}

func (t *stdioServerTransport) SetSessionManager(m sessionManager) {
	t.sessionManager = m
}

func (t *stdioServerTransport) Shutdown(userCtx context.Context, serverCtx context.Context) error {
	t.cancel()

	if err := t.reader.Close(); err != nil {
		return err
	}

	select {
	case <-t.receiveShutDone:
		return nil
	case <-serverCtx.Done():
		return nil
	case <-userCtx.Done():
		return userCtx.Err()
	}
}

func (t *stdioServerTransport) receive(ctx context.Context) {
	s := bufio.NewScanner(t.reader)

	for s.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			// filter empty messages
			// filter space messages and \t messages
			if len(bytes.TrimFunc(s.Bytes(), func(r rune) bool { return r == ' ' || r == '\t' })) == 0 {
				t.logger.Debugf("skipping empty message")
				continue
			}
			if err := t.receiver.Receive(ctx, stdioSessionID, s.Bytes()); err != nil {
				t.logger.Errorf("receiver failed: %v", err)
				continue
			}
		}
	}

	if err := s.Err(); err != nil {
		if !errors.Is(err, io.ErrClosedPipe) { // This error occurs during unit tests, suppressing it here
			t.logger.Errorf("server server unexpected error reading input: %v", err)
		}
		return
	}
}
