package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	auth "./auth"
	dirutils "./utils"
	"github.com/getlantern/systray"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	// path = filepath.Join(h.staticPath, path) uncomment for non-windows
	path = filepath.Join(h.staticPath, r.URL.Path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	onExit := func() {
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}
	systray.Run(onReady, onExit)
}

func serverStart() {

	auth.InitBaseForUser()
	router := mux.NewRouter()

	// TODO automate environment change
	// currDir := dirutils.CurrentDir()
	currDir := "dist"
	fmt.Println("current dir is " + currDir)
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/validateUser", auth.ValidateUser).Methods("GET")
	spa := spaHandler{staticPath: currDir, indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func onReady() {
	// systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetIcon(dirutils.GetIcon("icon/iconwin.ico"))
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Lantern")
	mUrl := systray.AddMenuItem("Open Server", "ui server")
	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetIcon(dirutils.GetIcon("icon/iconwin.ico"))
		//systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Any Camera Controller")
		systray.SetTooltip("Any Open Server")

		for {
			select {
			case <-mUrl.ClickedCh:
				open.Run("http://127.0.0.1:8000")
			}
		}
	}()

	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
	serverStart()
}
