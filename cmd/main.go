package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/blaze-d83/hardware-monitor-go/internal/hardware"
	"github.com/coder/websocket"
)

type Server struct {
	subscriberMessageBuffer int
	mux                     http.ServeMux
	subscribersMutex        sync.Mutex
	subscribers             map[*Subscriber]struct{}
}

type Subscriber struct {
	msgs chan []byte
}

func NewServer() *Server {
	s := &Server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*Subscriber]struct{}),
	}
	s.mux.Handle("/", http.FileServer(http.Dir("./static")))
	s.mux.HandleFunc("/ws", s.subscriberHandler)
	return s
}

func (s *Server) subscriberHandler(w http.ResponseWriter, r *http.Request) {
	err := s.subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Server) addSubscriber(subscriber *Subscriber) {
	s.subscribersMutex.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscribersMutex.Unlock()
	fmt.Println("Added a new subscriber: ", subscriber)
}

func (s *Server) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var c *websocket.Conn
	subscriber := &Subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}
	s.addSubscriber(subscriber)

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)
	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			err := c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}

func (s *Server) broadcast(msg []byte) {
	s.subscribersMutex.Lock()
	for subscriber := range s.subscribers {
		subscriber.msgs <- msg
	}
	s.subscribersMutex.Unlock()
}

func main() {
	fmt.Println("Starting system monitor..")
	srv := NewServer()
	go func(s *Server) {
		for {
			systemSection, err := hardware.GetSystemSection()
			if err != nil {
				fmt.Println(err)
			}

			diskSection, err := hardware.GetDiskSection()
			if err != nil {
				fmt.Println(err)
			}

			cpuSection, err := hardware.GetCPUSection()
			if err != nil {
				fmt.Println(err)
			}

			timeStamp := time.Now().Format("2006-01-02 15:04:05")

			html :=
				`<div hx-swap-oob="innerHTML:#update-timestamp">` + timeStamp + `</div>` +
					`<div hx-swap-oob="innerHTML:#system-data">` + systemSection + `</div>` +
					`<div hx-swap-oob="innerHTML:#disk-data">` + diskSection + `</div>` +
					`<div hx-swap-oob="innerHTML:#cpu-data">` + cpuSection + `</div>`

			s.broadcast([]byte(html))

			time.Sleep(2 * time.Second)
		}

	}(srv)

	err := http.ListenAndServe(":8080", &srv.mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
