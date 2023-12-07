package command

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"runtime/debug"
	"sync"
)

// readLine is a function that reads lines from a reader and sends them to a channel.
// It takes a pointer to a sync.WaitGroup, a channel of strings, and an io.ReadCloser as parameters.
// The function creates a new bufio.Reader with the io.ReadCloser, and then reads lines from the reader in a loop.
// If the reader returns io.EOF or an error, the function returns.
// Otherwise, the function sends the line to the channel.
// The function also recovers from any panics and logs the panic and the stack trace.
func readLine(wg *sync.WaitGroup, out chan string, reader io.ReadCloser) {
	// Recover from any panics and log the panic and the stack trace.
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, string(debug.Stack()))
		}
	}()
	// Signal that this function is done when it returns.
	defer wg.Done()
	// Create a new bufio.Reader with the io.ReadCloser.
	r := bufio.NewReader(reader)
	// Read lines from the reader in a loop.
	for {
		line, _, err := r.ReadLine()
		// If the reader returns io.EOF or an error, return.
		if err == io.EOF || err != nil {
			return
		}
		// Send the line to the channel.
		out <- string(line)
	}
}

// ExecCommandRealTimeOutput is a function that executes a command and sends its real-time output to a channel.
// It takes a channel of strings, a string representing the command, and a variadic parameter of strings representing the arguments of the command as parameters.
// The function creates a new exec.Cmd with the command and the arguments, and gets its stdout and stderr pipes.
// The function then starts the command, and creates two goroutines that read lines from the stdout and stderr pipes and send them to the channel.
// The function waits for the command to finish, and returns any error that occurred.
func ExecCommandRealTimeOutput(out chan string, name string, arg ...string) error {
	// Create a new exec.Cmd with the command and the arguments.
	cmd := exec.Command(name, arg...)
	// Get the stdout and stderr pipes of the command.
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	// Start the command.
	if err := cmd.Start(); err != nil {
		return err
	}
	// Create a new sync.WaitGroup.
	wg := sync.WaitGroup{}
	// Wait for the goroutines to finish when this function returns.
	defer wg.Wait()
	// Add two to the WaitGroup counter.
	wg.Add(2)
	// Create a goroutine that reads lines from the stdout pipe and sends them to the channel.
	go readLine(&wg, out, stdout)
	// Create a goroutine that reads lines from the stderr pipe and sends them to the channel.
	go readLine(&wg, out, stderr)
	// Wait for the command to finish, and return any error that occurred.
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
