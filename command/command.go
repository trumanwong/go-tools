package command

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"runtime/debug"
	"sync"
)

func readLine(wg *sync.WaitGroup, out chan string, reader io.ReadCloser) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, string(debug.Stack()))
		}
	}()
	defer wg.Done()
	r := bufio.NewReader(reader)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF || err != nil {
			return
		}
		out <- string(line)
	}
}

// ExecCommandRealTimeOutput 读取命令行实时输出
// Exec command and read real time output
func ExecCommandRealTimeOutput(out chan string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(2)
	go readLine(&wg, out, stdout)
	go readLine(&wg, out, stderr)
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
