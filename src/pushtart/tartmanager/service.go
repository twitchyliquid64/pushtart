package tartmanager

import (
	"io"
	"pushtart/config"
	"pushtart/logging"
	"strings"
	"time"
)

func tartLogRoutine(tart config.Tart, reader io.ReadCloser, errReader io.ReadCloser) {
	buf := make([]byte, 4096*2)

	go func() {
		buf2 := make([]byte, 4096*2)
		for {
			n, err := errReader.Read(buf2)
			if err != nil {
				if err != io.EOF {
					logging.Error("tartmanager-service", "Stderr read error: "+err.Error())
				}
				break
			} else {
				if tart.LogStdout {
					spl := strings.Split(strings.Replace(string(buf2[:n]), "\r", "", -1), "\n")
					for _, line := range spl {
						if len(line) > 0 {
							logging.Info(tart.Name, line)
						}
					}
				}
			}
		}
	}()

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				logging.Error("tartmanager-service", "Read error: "+err.Error())
			}
			logging.Info("tartmanager-service", tart.Name+" is shutting down.")

			tart = Get(tart.PushURL)

			if tart.RestartOnStop {
				time.Sleep(time.Duration(tart.RestartDelaySecs) * time.Second)
				tart = Get(tart.PushURL)
			}

			tart.IsRunning = false
			tart.PID = -1
			Save(tart.PushURL, tart)
			if tart.RestartOnStop {
				logging.Info("tartmanager-service", tart.Name+" is restarting.")
				Start(tart.PushURL)
			}

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
