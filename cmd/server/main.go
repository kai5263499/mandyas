package main

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kai5263499/mandyas/domain"
	"github.com/kai5263499/mandyas/server"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Routine to reap zombies (it's the job of init)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go removeZombies(ctx, &wg)

	mainCmd := strings.Join(os.Args[1:], " ")

	// Launch main command
	var mainRC int
	log.Infof("Main command launched : %s", mainCmd)
	err := run(mainCmd)
	if err != nil {
		log.Errorf("Main command failed %#+v\n", err)
		mainRC = 1
	} else {
		log.Errorf("Main command exited")
	}

	// Wait removeZombies goroutine
	cleanQuit(cancel, &wg, mainRC)
}

func removeZombies(ctx context.Context, wg *sync.WaitGroup) {
	for {
		var status syscall.WaitStatus

		// Wait for orphaned zombie process
		pid, _ := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)

		if pid <= 0 {
			// PID is 0 or -1 if no child waiting
			// so we wait for 1 second for next check
			time.Sleep(1 * time.Second)
		} else {
			// PID is > 0 if a child was reaped
			// we immediately check if another one
			// is waiting
			continue
		}

		// Non-blocking test
		// if context is done
		select {
		case <-ctx.Done():
			// Context is done
			// so we stop goroutine
			wg.Done()
			return
		default:
		}
	}
}

func run(command string) error {
	// Register chan to receive system signals
	sigs := make(chan os.Signal, 1)
	defer close(sigs)
	signal.Notify(sigs)
	defer signal.Reset()

	config := &domain.Config{
		GrpcPort: 9000,
	}

	var err error

	// Define command and rebind
	// stdout and stdin
	cmd := exec.Command("sh", "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("error opening stdout pipe:", err)
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
		return err
	}

	grpcServer, _ := server.New(config, stdout, stdin)

	log.Infof("starting grpcServer with pipe %#+v", stdout)
	grpcServer.Start()

	// Create a dedicated pidgroup
	// used to forward signals to
	// main process and all children
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Goroutine for signals forwarding
	go func() {
		for sig := range sigs {
			// Ignore SIGCHLD signals since
			// thez are only usefull for cs-init
			if sig != syscall.SIGCHLD {
				// Forward signal to main process and all children
				syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
			}
		}
	}()

	log.Infof("starting command")
	err = cmd.Start()
	if err != nil {
		log.Fatalf("error starting command %#+v", err)
		return err
	}

	// Wait for command to exit
	return cmd.Wait()
}

func cleanQuit(cancel context.CancelFunc, wg *sync.WaitGroup, code int) {
	// Signal zombie goroutine to stop
	// and wait for it to release waitgroup
	cancel()
	wg.Wait()

	os.Exit(code)
}
