package daemon

import (
	"log"
	"net"
	"net/http"
	"os"
)

func ServeIPC(sockAddr string) {
	if err := os.RemoveAll(sockAddr); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("unix", sockAddr)
	if err != nil {
		log.Fatalln("neoman: could not listen to the Unix Domain Socket", err)
	}
	defer listener.Close()
	log.Println("Listening to Unix Domain Socket at", sockAddr)

	handlers := Handlers{}
	http.Handle("GET /ping", http.HandlerFunc(handlers.Pong))
	http.Handle("GET /is-trusted", http.HandlerFunc(handlers.CheckTrustedRemoteAccount))

	if err := http.Serve(listener, nil); err != nil {
		log.Fatalln("neoman: could not serve IPC Unix Socket", err)
	}
}
