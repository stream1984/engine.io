package example

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"flag"

	"testing"

	eio "github.com/jjeffcaii/engine.io"
)

var server eio.Engine

func init() {
	flag.Parse()
	server = eio.NewEngineBuilder().Build()
	http.HandleFunc("/conns", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(fmt.Sprintf("totals: %d", server.CountClients())))
	})

}

func TestNothing(t *testing.T) {
	server.OnConnect(func(socket eio.Socket) {
		log.Println("========> socket connect:", socket.Id())
		socket.OnMessage(func(data []byte) {
			// do nothing.
			log.Println("===> got message:", string(data))
		})
		socket.OnClose(func(reason string) {
			log.Println("========> socket closed:", socket.Id())
		})
	})
	http.HandleFunc(eio.DEFAULT_PATH, server.Router())
	log.Fatalln(http.ListenAndServe(":3000", nil))
}

func TestEcho(t *testing.T) {
	tick := time.NewTicker(5 * time.Second)
	kill := make(chan uint8)
	go func() {
		for {
			select {
			case <-tick.C:
				for _, it := range server.GetClients() {
					it.Send(fmt.Sprintf("AUTO_BRD: %s", time.Now().Format(time.RFC3339)))
				}
				break
			case <-kill:
				tick.Stop()
				return
			}
		}
	}()

	defer func() {
		close(kill)
		server.Close()
	}()
	server.OnConnect(func(socket eio.Socket) {
		log.Println("========> socket connect:", socket.Id())
		socket.OnMessage(func(data []byte) {
			log.Printf("got string data: %+v\n", data)
			socket.Send(fmt.Sprintf("ECHO1: %s", data))
			socket.Send(fmt.Sprintf("ECHO2: %s", data))
		})
		socket.OnClose(func(reason string) {
			log.Println("========> socket closed:", socket.Id())
		})
	})
	log.Fatalln(server.Listen(":3000"))
}