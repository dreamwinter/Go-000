package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var server1, server2 *http.Server
	var shutdownWaitGroup sync.WaitGroup
	serverPort1 := 9000
	serverMux1 := http.NewServeMux()
	serverPort2 := 9001
	serverMux2 := http.NewServeMux()

	serverMux1.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(15 * time.Second)
		fmt.Fprint(res, "HELLO")
	})

	serverMux2.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(10 * time.Second)
		fmt.Fprint(res, "WORLD!")
	})

	g, ctx := errgroup.WithContext(context.Background())
	// start server 1
	g.Go(func() error {
		server1 = &http.Server{
			Addr:    fmt.Sprintf(":%v", serverPort1),
			Handler: serverMux1,
		}
		go func() {
			shutdownWaitGroup.Add(1)
			defer shutdownWaitGroup.Done()
			<-ctx.Done()
			shutdownctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server1.Shutdown(shutdownctx); err != nil {
				log.Printf("Server 1 forced to shutdown: %v", err)
			} else {
				log.Printf("Shutdown Server 1")
			}

		}()
		log.Printf("Start Server1")
		return server1.ListenAndServe()
	})
	// start server 2
	g.Go(func() error {
		server2 = &http.Server{
			Addr:    fmt.Sprintf(":%v", serverPort2),
			Handler: serverMux2,
		}
		go func() {
			shutdownWaitGroup.Add(1)
			defer shutdownWaitGroup.Done()
			<-ctx.Done()
			shutdownctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server2.Shutdown(shutdownctx); err != nil {
				log.Printf("Server 2 forced to shutdown: %v", err)
			} else {
				log.Printf("Shutdown Server 2")
			}
		}()
		log.Printf("Start Server2")
		return server2.ListenAndServe()
	})

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case s := <-quit:
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

	if err := g.Wait(); err != nil {
		shutdownWaitGroup.Wait()
		log.Printf("Server Error:%v\n", err)
	}
}
