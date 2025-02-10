package logger

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var logChannel = make(chan string, 100)

func StartLogListener(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case logMsg := <-logChannel:
				fmt.Println(logMsg)
			case <-ctx.Done():
				fmt.Println("Log Listener shutting down...")
				return
			}
		}
	}()
}

func LogAction(action string, bookID string) {
	logChannel <- fmt.Sprintf("[%s] %s - Book ID: %s", time.Now().Format(time.RFC3339), action, bookID)
}
