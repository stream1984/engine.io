# Engine.IO (WARNING: STILL WORKING!!!!)
Server-side [Engine.IO](https://github.com/socketio/engine.io) in Golang.

## Example

``` golang
package main

import (
	"flag"
	"log"

	eio "engine.io"
)

func main() {
	flag.Parse()
	server := eio.NewEngineBuilder().Build()
	server.OnConnect(func(socket eio.Socket) {
		socket.Send("Hello world!")
		socket.OnMessage(func(data []byte) {
			log.Println("recieve:", string(data))
		})
		socket.OnClose(func(reason string) {
			log.Println("socket closed: ", socket.Id())
		})
	})
	log.Fatalln(server.Listen(":3000"))
}

```
