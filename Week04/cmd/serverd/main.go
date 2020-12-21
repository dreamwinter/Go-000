package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	server := InitializeGRPCServer()
	g, ctx := errgroup.WithContext(context.Background())
	err := server.Start(ctx, g)
	if err != nil {
		log.Fatalf("Server exit unexpectedly due to %v", err)
	}
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			select {
			case <-ctx.Done():
				return nil
			case s := <-c:
				log.Printf("get a signal %s", s.String())
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					return errors.New("Close by signal " + s.String())
				case syscall.SIGHUP:
				default:
					return errors.New("Undefined signal")
				}
			}
		}
	})
	if err = g.Wait(); err != nil {
		log.Printf("Server Error:%v\n", err)
	}
}
