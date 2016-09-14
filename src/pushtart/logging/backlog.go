package logging

import (
	"container/ring"
	"fmt"
	"time"
)

//implements a small history of log messages

const defaultBacklogSize = 15

var backlog *ring.Ring

func init() {
	backlog = ring.New(defaultBacklogSize)
}

func addToBacklog(component, typ, msg string) {
	nmsg := LogMessage{
		Component: component,
		Type:      typ,
		Message:   msg,
		Created:   time.Now().Unix(),
	}

	backlog.Value = nmsg
	backlog = backlog.Next()
}

//GetBacklog returns the last few log messages.
func GetBacklog() []LogMessage {
	var output []LogMessage

	cursor := backlog
	for i := 0; i < cursor.Len(); i++ {
		if cursor.Value == nil {
		} else {
			output = append(output, cursor.Value.(LogMessage))
		}
		cursor = cursor.Next()
	}
	return output
}

func debugPrintBacklog() {
	cursor := backlog
	for i := 0; i < cursor.Len(); i++ {
		if cursor.Value == nil {
		} else {
			fmt.Println(cursor.Value.(LogMessage).Message)
		}
		cursor = cursor.Next()
	}
}
