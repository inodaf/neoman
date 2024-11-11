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

	// handlers := Handlers{}
	listener, err := net.Listen("unix", sockAddr)
	if err != nil {
		log.Fatalln("neoman: could not listen to the Unix Domain Socket", err)
	}
	defer listener.Close()

	if err := http.Serve(listener, http.HandlerFunc(handler)); err != nil {
		log.Fatalln("neoman: could not serve IPC Unix Socket", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling %s", r.URL.Path)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Pong"))
}
