package main

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"net/url"
	"encoding/json"
	"io"
)

type Task struct {
	Id 	int `json:"id"`
	State	int `json:"state"`
}

var databaseLocation string
var storageLocation  string

func main() {
	if !registerInKVStore2() {
		return
	}
	keyValueStoreAddress := os.Args[2]

	response, err := http.Get("http://" + keyValueStoreAddress + "/get?key=databaseAddress")
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: can't get database address.")
		fmt.Println(response.Body)
		return
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	databaseLocation = string(data)

	response, err = http.Get("http://" + keyValueStoreAddress + "/get?key=storageAddress")
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: can't get storage address.")
		fmt.Println(response.Body)
		return
	}
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	storageLocation = string(data)

	http.HandleFunc("/new", newImage)
	http.HandleFunc("/get", getImage)
	http.HandleFunc("/isReady", isReady)
	http.HandleFunc("/getNewTask", getNewTask1)
	http.HandleFunc("/registerTaskFinished", registerTaskFinished)
	http.ListenAndServe(":3003", nil)
}

func registerInKVStore2() bool {
	if len(os.Args) < 3 {
		fmt.Println("Error: Too few arguments.")
		return false
	}
	databaseAddress := os.Args[1]
	keyValueStoreAddress := os.Args[2]

	response, err := http.Post("http://" + keyValueStoreAddress + "/set?key=databaseAddress&value=" + databaseAddress, "", nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: Failure when contacting key-value store: ", string(data))
		return false
	}
	return true
}

func newImage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		response, err := http.Post("http://" + databaseLocation + "/newTask", "text/plain", nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}
		id, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = http.Post("http://" + storageLocation + "/sendImage?id=" + string(id) + "&state=working", "image", r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}
		fmt.Fprint(w, string(id))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Only POST accepted")
	}
}

func getImage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}
		if len(values.Get("id")) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Wrong input")
			return
		}

		response, err := http.Get("http://" + storageLocation + "/getImage?id=" + values.Get("id") + "&state=finished")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}

		_, err = io.Copy(w, response.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Only GET accepted")
	}
}

func isReady(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		if len(values.Get("id")) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Wrong input")
			return
		}

		response, err := http.Get("http://" + databaseLocation + "/getById?id=" + values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		myTask := Task{}
		json.Unmarshal(data, &myTask)

		if myTask.State == 2 {
			fmt.Fprint(w, "1")
		} else {
			fmt.Fprint(w, "0")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Only GET accepted")
	}
}

func getNewTask1(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		response, err := http.Post("http://" + databaseLocation + "/getNewTask", "text/plain", nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}

		_, err = io.Copy(w, response.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Only Post accepted")
	}
}

func registerTaskFinished(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}
		if len(values.Get("id")) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Wrong input")
			return
		}

		response, err := http.Post("http://" + databaseLocation + "/finishTask?id=" + values.Get("id"), "text/plain", nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}

		_, err = io.Copy(w, response.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Only POST accepted")
	}
}