package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type DownloadServer struct {
	service_address string
	router          *Router
	downloadStack   *Stack
	canonical_ok    []byte
	download_root   string
	stop_worker     chan int
}

func (self *DownloadServer) createDownloadFromRequest(request *http.Request) (*Download, error) {
	var new_dowload *Download = new(Download)
	new_dowload.Name = request.FormValue("name")
	new_dowload.Directory = request.FormValue("path")

	err := json.Unmarshal([]byte(request.FormValue("files")), &(new_dowload.Files))
	if err != nil {
		fmt.Println(err)
	}
	return new_dowload, err
}

func (self *DownloadServer) composeAddr(port int) {
	var host string = os.Getenv("DOWNLOAD_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	self.service_address = fmt.Sprintf("%s:%d", host, port)
}

func (self *DownloadServer) downloadWorker() {
	for {
		select {
		case <-time.After(time.Second * 1):
			if self.downloadStack.count > 0 {
				var current_download *Download
				current_download, err := self.downloadStack.pop()
				fmt.Printf("Download '%s' is now been processed\n", current_download.Name)
				if err == nil {
					self.processDownload(current_download)
				}
			}
		case <-self.stop_worker:
			return
		}
	}
}

//Download route should have de form {name string, path string, files []string}
func (self *DownloadServer) handleDownloads(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		fmt.Printf("Got a Download request with name '%s'\n", request.FormValue("name"))
		var new_download *Download
		new_download, err := self.createDownloadFromRequest(request)
		if err != nil {
			return
		}

		self.downloadStack.push(new_download)

		response.WriteHeader(200)
		response.Write(self.canonical_ok)
	} else {
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(response, "Method '%s' is not allowed... gai", request.Method)
	}
}

func (self *DownloadServer) handleTests(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	fmt.Fprint(response, request.URL.Path)
}

func (self *DownloadServer) processDownload(download_obj *Download) {
	var fulldownload_path string = fmt.Sprintf("%s/%s", self.download_root, download_obj.Directory)
	if pathExists(fulldownload_path) {
		files_saving_point := fmt.Sprintf("%s/%s", fulldownload_path, download_obj.Name)
		if !pathExists(files_saving_point) {
			os.Mkdir(files_saving_point, 0777)
			for h, file_url := range download_obj.Files {
				response, err := http.Get(file_url)
				if err != nil {
					fmt.Printf("Error from '%s': %s", file_url, err.Error())
					continue
				}

				if response.StatusCode == 200 {
					file_name := fmt.Sprintf("%s/%s_%d%s", files_saving_point, download_obj.Name, h, filepath.Ext(file_url))
					if !pathExists(file_name) {
						file_desc, err := os.Create(file_name)
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						_, err = io.Copy(file_desc, response.Body)
						if err != nil {
							fmt.Printf("Error on file '%s': %s\n", file_name, err.Error())
						} else {
							fmt.Printf("Succes on file (%d)\r", h+1)
						}
					} else {
						fmt.Printf("File '%s' already exists, skipping but something must be wrong..\n", file_name)
					}
				} else {
					fmt.Printf("Response from '%s' was: %d\n", file_url, response.StatusCode)
				}
			}
			fmt.Printf("\nDownloaded %d files\n", len(download_obj.Files))
		} else {
			fmt.Printf("Saving point '%s' already exists, aborting...\n", files_saving_point)
		}
	} else {
		fmt.Printf("Cannot process download '%s' because path '%s' doesnt exist\n", download_obj.Name, fulldownload_path)
	}
}

func (self *DownloadServer) run() {
	go self.downloadWorker()

	self.router.registerRoute(regexp.MustCompile("/download"), self.handleDownloads)
	self.router.registerRoute(regexp.MustCompile(`/test-\d$`), self.handleTests)

	fmt.Printf("Lisiting on '%s'\n", self.service_address)
	http.ListenAndServe(self.service_address, self.router)
}

func createDownloadServer(port int) *DownloadServer {
	var new_dserver *DownloadServer = new(DownloadServer)
	new_dserver.router = createRouter()
	new_dserver.downloadStack = new(Stack)
	new_dserver.downloadStack.init()
	new_dserver.canonical_ok = []byte("{\"response\": \"ok\"}")
	new_dserver.composeAddr(port)
	new_dserver.download_root = os.Getenv("DOWNLOAD_ROOT")
	if new_dserver.download_root == "" {
		new_dserver.download_root = "."
	}

	return new_dserver
}

func main() {
	var download_server *DownloadServer = createDownloadServer(5006)
	download_server.run()
}
