package command

import (
	"log"
	"testing"
)

func TestExecCommandRealTimeOutput(t *testing.T) {
	out := make(chan string)
	defer close(out)
	go func() {
		for {
			str, ok := <-out
			if !ok {
				break
			}
			log.Println(str)
		}
	}()
	args := []string{"-c", "ping www.trumanwl.com"}
	if err := ExecCommandRealTimeOutput(out, "bash", args...); err != nil {
		t.Error(err)
	}
}
