package daemon

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/inodaf/neoman/internal"
	"github.com/inodaf/neoman/internal/handlers"
)

func ServeSocket(sockAddr string, db *sql.DB) {
	if err := os.RemoveAll(sockAddr); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("unix", sockAddr)
	if err != nil {
		log.Fatalln("neoman: could not listen to the socket", err)
	}
	defer listener.Close()
	log.Println("Listening to socket at", sockAddr)

	mux := http.NewServeMux()
	mux.Handle("/ping", handlers.PingHandler{})
	mux.Handle("/trust", handlers.TrustHandler{DB: db})

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
func ServeWeb(port string, db *sql.DB) {
	mux := http.NewServeMux()
	log.Println("Listening to web at", port)

	mux.HandleFunc("GET /{localRepo}", func(w http.ResponseWriter, r *http.Request) {
		repo := r.PathValue("localRepo")

		dir, err := internal.DocsRegistryDir()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, path.Join(dir, "local", repo, "docs", "index.md"))
	})

	mux.HandleFunc("GET /{owner}/{repo}", func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")

		dir, err := internal.DocsRegistryDir()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, path.Join(dir, "remote", owner, repo, "README.md"))
	})

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalln("neoman: could not serve web application", err)
	}
}