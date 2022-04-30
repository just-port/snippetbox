package main

import (
  "flag"
  "log"
  "net/http"
  "path/filepath"
  "os"
)

func main() {
  addr := flag.String("addr", ":4000", "HTTP network address")

  flag.Parse()

  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

  errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  mux := http.NewServeMux()

  fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("/", home)
  mux.HandleFunc("/snippet/view", snippetView)
  mux.HandleFunc("/snippet/create", snippetCreate)

  srv := &http.Serve{
    Addr: *addr,
    ErrorLog: errorLog,
    Handler: mux,
  }

  infoLog.Printf("Starting server on %s", *addr)
  err := srv.ListenAndServe(*addr, mux)
  errorLog.Fatal(err)
}

type neuteredFileSystem struct {
  fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
  f, err := nfs.fs.Open(path)
  if err != nil {
    return nil, err
  }

  s, err := f.Stat()
  if s.IsDir() {
    index := filepath.Join(path, "index.html")
    if _, err := nfs.fs.Open(index); err != nil {
      closeErr := f.Close()
      if closeErr != nil {
        return nil, closeErr
      }

      return nil, err
    }
  }

  return f, nil
}
