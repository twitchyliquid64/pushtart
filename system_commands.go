package main

import (
	"fmt"
	"io"
	"pushtart/logging"
	"time"
)

func logMsgs(params map[string]string, w io.Writer) {
	logMsgs := logging.GetBacklog()

	for _, msg := range logMsgs {
		fmt.Fprintln(w, time.Unix(msg.Created, 0).Format(time.ANSIC), "["+msg.Component+"]", msg.Message)
	}
}
