package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/blaze-d83/hardware-monitor-go/internal/hardware"
	"github.com/blaze-d83/hardware-monitor-go/internal/templates"
	"github.com/coder/websocket"
)

type Subscriber struct {
	msgs chan []byte
}

func (s *Server) subscriberHandler(w http.ResponseWriter, r *http.Request) {
	err := s.subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println("Error in subscription:", err)
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
	}
}

func (s *Server) addSubscriber(subscriber *Subscriber) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()
	s.subscribers[subscriber] = struct{}{}
	fmt.Println("Added a new subscriber: ", subscriber)
}

func (s *Server) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
			if err := c.Write(ctx, websocket.MessageText, msg); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Server) broadcast(msg []byte) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()
	for subscriber := range s.subscribers {
		select {
		case subscriber.msgs <- msg:
		default:
		}
	}
}

func (s *Server) MonitorHardware() {
	go func() {
		for {
			systemInfo, err := hardware.GetSystemInfo()
			if err != nil {
				fmt.Println("Error retrieving system info:", err)
				continue
			}

			diskInfo, err := hardware.GetDiskInfo()
			if err != nil {
				fmt.Println("Error retrieving disk info:", err)
				continue
			}

			cpuInfo, err := hardware.GetCPUInfo()
			if err != nil {
				fmt.Println("Error retrieving CPU info:", err)
				continue
			}

			timeStamp := time.Now().Format("2006-01-02 15:04:05")

			monitorData := &Metrics{
				HostName:       systemInfo.HostName,
				TotalMemory:    systemInfo.TotalMem,
				UsedMemory:     systemInfo.UsedMem,
				OS:             systemInfo.OS,
				TotalDiskSpace: diskInfo.TotalDiskSpace,
				FreeDiskSpace:  diskInfo.FreeDiskSpace,
				CPUModelName:   cpuInfo.ModelName,
				Cores:          uint8(cpuInfo.Cores),
			}

			monitorComponent := templates.MonitorSection(
				timeStamp,
				monitorData.HostName,
				monitorData.TotalMemory,
				monitorData.UsedMemory,
				monitorData.OS,
				monitorData.TotalDiskSpace,
				monitorData.FreeDiskSpace,
				monitorData.CPUModelName,
				monitorData.Cores)

			var buf bytes.Buffer
			err = monitorComponent.Render(context.Background(), &buf)
			if err != nil {
				fmt.Println("Error rendering the component:", err)
				continue
			}

			monitorBytes := buf.Bytes()

			s.broadcast(monitorBytes)
			time.Sleep(2 * time.Second)
		}
	}()
}
