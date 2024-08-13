package server

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	server := NewServer()
	assert.NotNil(t, server)
	assert.Equal(t, 10, server.subscriberMessageBuffer)
	assert.NotNil(t, server.subscribers)
	assert.NotNil(t, server.mux)
}

func TestServer_Start(t *testing.T) {
	server := NewServer()
	go server.MonitorHardware()

	ts := httptest.NewServer(server.mux)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "test complete")

	msgChan := make(chan []byte)
	go func() {
		defer close(msgChan)
		for {
			_, msg, err := conn.Read(ctx)
			if err != nil {
				return
			}
			msgChan <- msg
		}
	}()

	select {
	case msg := <-msgChan:
		assert.NotEmpty(t, msg)
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestServer_addSubscriber(t *testing.T) {
	server := NewServer()
	subscriber := &Subscriber{
		msgs: make(chan []byte, server.subscriberMessageBuffer),
	}

	server.addSubscriber(subscriber)

	assert.Contains(t, server.subscribers, subscriber)
}

func TestServer_broadcast(t *testing.T) {
	server := NewServer()
	subscriber1 := &Subscriber{
		msgs: make(chan []byte, server.subscriberMessageBuffer),
	}

	subscriber2 := &Subscriber{
		msgs: make(chan []byte, server.subscriberMessageBuffer),
	}

	server.addSubscriber(subscriber1)
	server.addSubscriber(subscriber2)

	msg := []byte("test message")
	server.broadcast(msg)

	assert.Equal(t, msg, <-subscriber1.msgs)
	assert.Equal(t, msg, <-subscriber2.msgs)
}

func TestServer_MonitorHardware(t *testing.T) {
	server := NewServer()
	server.MonitorHardware()

	time.Sleep(3 * time.Second)

	subscriber := &Subscriber{
		msgs: make(chan []byte, server.subscriberMessageBuffer),
	}
	server.addSubscriber(subscriber)

	select {
	case msg := <-subscriber.msgs:
		assert.NotEmpty(t, msg)
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for broadcast message")
	}
}

