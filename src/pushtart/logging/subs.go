package logging

// This package implements a means for other packages to recieve the log messages, in a pub-sub like architecture.
// Implementers call Subscribe(chan LogMessage), and log messages will be pushed to the given channel IF THERE IS SPACE.
// Unsubscribe is called when done.


import (
  "time"
  "sync"
)

type LogMessage struct {
  Component string
  Type string
  Message string
  Created int64
}

var subscribers  = map[chan LogMessage]bool{}
var subStructLock sync.Mutex

func Subscribe(in chan LogMessage){ //DO NOT LOG WITHIN THIS MSG - DEADLOCK
  subStructLock.Lock()
  defer subStructLock.Unlock()

  subscribers[in] = true
}

func Unsubscribe(in chan LogMessage){ //DO NOT LOG WITHIN THIS MSG - DEADLOCK
  subStructLock.Lock()
  defer subStructLock.Unlock()

  delete(subscribers, in)
}

func publishMessage(component, typ, msg string){ //DO NOT LOG WITHIN THIS MSG - DEADLOCK
  pkt := LogMessage{
    Component: component,
    Type: typ,
    Message: msg,
    Created: time.Now().Unix(),
  }

  subStructLock.Lock()
  defer subStructLock.Unlock()
  for ch, _ := range subscribers {
    select { //prevents blocking if a channel is full
      case ch <- pkt:
      default:
    }
  }
}
