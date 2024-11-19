package daemon

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/inodaf/neoman/internal"
	"github.com/inodaf/neoman/internal/handlers"
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

	http.Handle("/ping", handlers.PingHandler{})
	http.Handle("/trust", handlers.TrustHandler{DB: internal.DB})

	if err := http.Serve(listener, nil); err != nil {
		log.Fatalln("neoman: could not serve IPC Unix Socket", err)
	}
}
