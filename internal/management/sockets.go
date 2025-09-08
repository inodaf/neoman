package management

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/inodaf/neoman/packages/config"
)

var UnixSockClient http.Client = http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", config.AppSockAddr)
		},
	},
}

// SocketClientPing tries to ping the daemon through the Unix socket.
// If the daemon is not running, it returns an error.
func SocketClientPing() error {
	resource := url.URL{Host: "unix", Scheme: "http", Path: "/ping"}
	_, err := UnixSockClient.Get(resource.String())

	return err
}

// SocketServeIPC serves the IPC Unix Domain Socket (UDS) for communication
// between nman and the nmand (daemon).
func SocketServeIPC(db *sql.DB) {
	if err := os.RemoveAll(config.AppSockAddr); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("unix", config.AppSockAddr)
	if err != nil {
		log.Fatalln("neoman: could not listen to the socket", err)
	}

	defer listener.Close()
	log.Println("Listening to socket at", config.AppSockAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("pong"))
	})

	if err := http.Serve(listener, mux); err != nil {
		log.Fatalln("neoman: could not serve IPC Unix Socket", err)
	}
}

// Handle
// - neoman.local/:owner
// - neoman.local/:localRepo
// - neoman.local/:owner/:repo/
// - neoman.local/:owner/:repo/:[introduction/getting-started]
// - neoman.local/?q=MyQuery
// - neoman.local/
// - -
//
// edge cases
// owner and repo sanitized + URL escaped/unescaped
// owner is not found
// repo is not found (fetch it)
// default .md file to display (index.md) or fallback to root's README.md
// -- CLI (ask which page want to open by default for that repo)
// what happens if owner name is same as a local repo?
func SocketServeTCP(db *sql.DB) {
	mux := http.NewServeMux()
	log.Println("Listening to web at", config.AppWebAppPort)

	mux.HandleFunc("GET /{localRepo}", func(w http.ResponseWriter, r *http.Request) {
		repo := r.PathValue("localRepo")

		dir, err := config.DocsRegistryDir()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, path.Join(dir, "local", repo, "docs", "README.md"))
	})

	mux.HandleFunc("GET /{owner}/{repo}", func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")

		dir, err := config.DocsRegistryDir()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, path.Join(dir, "remote", owner, repo, "README.md"))
	})

	if err := http.ListenAndServe(config.AppWebAppPort, mux); err != nil {
		log.Fatalln("neoman: could not serve web application", err)
	}
}
