package tartmanager

import (
	"io"
	"pushtart/config"
	"pushtart/logging"
	"strings"
)

func tartLogRoutine(tart config.Tart, reader io.ReadCloser) {
	buf := make([]byte, 4096*2)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				logging.Error("tartmanager-service", "Read error: "+err.Error())
			}
			logging.Info("tartmanager-service", tart.Name+" is shutting down.")

			tart = Get(tart.PushURL)
			tart.IsRunning = false
			tart.PID = -1
			Save(tart.PushURL, tart)

			break
		}

		if tart.LogStdout {
			spl := strings.Split(strings.Replace(string(buf[:n]), "\r", "", -1), "\n")
			for _, line := range spl {
				if len(line) > 0 {
					logging.Info(tart.Name, line)
				}
			}
		}
	}
}
