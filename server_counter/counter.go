package server_counter

var (
	countChanger = make(chan bool)
	runningCount = 0
)

func CountServer(exit chan bool) {
	for {
		b := <-countChanger
		if b {
			runningCount++
		} else {
			runningCount--
			if runningCount <= 0 {
				exit <- true
				break
			}
		}
	}
}

func IncServerCount() {
	countChanger <- true
}

func DecServerCount() {
	countChanger <- false
}

func ServerRunning() bool {
	return runningCount > 0
}
