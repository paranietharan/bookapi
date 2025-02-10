package logger

import (
	"fmt"
	"time"
)

var logChannel = make(chan string, 100)

func StartLogListener() {
	go func() {
		for logMessage := range logChannel {
			fmt.Println(logMessage)
		}
	}()
}

func LogAction(method, bookID string) {
	logChannel <- fmt.Sprintf("[%s] %s /books/%s - Book ID: %s",
		time.Now().Format(time.RFC3339), method, bookID, bookID)
}
