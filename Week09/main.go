package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var nextID uint64

// Session is the struct to store active client session
type Session struct {
	id       uint64
	wchannel chan string
}

var sessions sync.Map

func main() {
	var shutdownWaitGroup sync.WaitGroup
	serverPort := 1235

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {

		listen, err := net.Listen("tcp", fmt.Sprintf(":%v", serverPort))
		if err != nil {
			log.Fatalf("listen error:%v\n", err)
		}
		log.Printf("Server Listen:%v", serverPort)
		go func() {
			shutdownWaitGroup.Add(1)
			defer shutdownWaitGroup.Done()
			<-ctx.Done()
			if err := listen.Close(); err != nil {
				log.Printf("TCP Server forced to shutdown: %v", err)
			} else {
				log.Printf("Shutdown TCP")
			}

		}()
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Printf("accept error:%v\n", err)
				return err
			}
			go handleConn(conn)
		}
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
				broadcast("Server is going to shutdown in 5 secs\r\n", 0)
				time.Sleep(5 * time.Second)
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

func handleConn(conn net.Conn) {
	defer conn.Close()

	rd := bufio.NewReader(conn)

	session := &Session{
		id:       idGenerator(),
		wchannel: make(chan string, 8),
	}
	sessions.Store(session.id, session)
	session.wchannel <- fmt.Sprintf("You ID is %v", session.id)
	go sendMsg(session.id, conn, session.wchannel)

	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			log.Printf("Session %v read error:%v\n", session.id, err)
			removeSession(session.id)
			return
		}
		if strings.HasPrefix(string(line), "/l") {
			// List all active sessions
			message := "Active sessions:\r\n"
			sessions.Range(func(key interface{}, vaule interface{}) bool {
				message = message + fmt.Sprintf("%v\r\n", key)
				return true
			})
			session.wchannel <- message
		} else if strings.HasPrefix(string(line), "/c") {
			// DM with session {id}
			tokens := strings.Split(string(line), " ")
			sessionID, err := strconv.ParseUint(tokens[1], 10, 64)
			if err != nil {
				log.Printf("cannot parse uint64: %v", err)
				continue
			}
			targetSession, ok := sessions.Load(sessionID)
			if ok {
				targetSession.(*Session).wchannel <- strings.Join(tokens[2:], " ")
			}
		} else if strings.HasPrefix(string(line), "/b") {
			// Broadcast message
			broadcast(string(line)[3:], session.id)
		} else if strings.HasPrefix(string(line), "/q") {
			// Exit the session
			removeSession(session.id)
			break
		} else {
			session.wchannel <- "UNKNOWN Command"
		}
	}
}

func broadcast(message string, excluded uint64) {
	sessions.Range(func(key interface{}, value interface{}) bool {
		if key.(uint64) != excluded {
			session := value.(*Session)
			session.wchannel <- message
		}
		return true
	})
}

func sendMsg(sessionID uint64, conn net.Conn, ch <-chan string) {
	defer conn.Close()
	wr := bufio.NewWriter(conn)

	for msg := range ch {
		if strings.HasPrefix(msg, "/q") {
			break
		}
		_, err := wr.WriteString(msg + "\r\n")
		if err != nil {
			removeSession(sessionID)
			log.Printf("Write failed :%v", err)
			return
		}
		wr.Flush()
		if err != nil {
			removeSession(sessionID)
			log.Printf("Flush failed :%v", err)
			return
		}
	}
}

func removeSession(sessionID uint64) {
	session, ok := sessions.Load(sessionID)
	if ok {
		close(session.(*Session).wchannel)
		sessions.Delete(sessionID)
		broadcast(fmt.Sprintf("User %v exited", sessionID), 0)
	}
}

func idGenerator() uint64 {
	return atomic.AddUint64(&nextID, 1)
}
