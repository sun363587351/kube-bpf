package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"

	operator "github.com/bpftools/kube-bpf"
)

// Version is the version of the Kubernetes BPF Operator
const Version ="0.0.0.dev"

var options operator.Config

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	sigs := make(chan os.Signal, 1)
	stop := make(chan struct{})
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM) // Push signals into channel

	wg := &sync.WaitGroup{} // Goroutines can add themselves to this to be waited on

	options.Labels = map[string]string{
		"operator": "bpf-operator",
		"version":  Version,
	}

	op := operator.New(options, zap.NewNop())
	op.Run(stop, wg)

	<-sigs // Wait for signals (this hangs until a signal arrives)
	fmt.Println("operator shutting down")

	close(stop) // Tell goroutines to stop themselves
	wg.Wait()   // Wait for all to be stopped

	return nil
}
