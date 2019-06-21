package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/felts94/queues"
	"github.com/felts94/stacks"
)

var (
	stack *stacks.Stack
	queue = queues.NewQueuePointer()
)

func main() {
	startHTTPServer()
}

func startHTTPServer() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8080"
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	// http.HandleFunc("/s/push", pushStack)
	// http.HandleFunc("/s/pop", popStack)
	// http.HandleFunc("/s/view", viewStack)
	http.HandleFunc("/q/enqueue", enqueue)
	http.HandleFunc("/q/dequeue", dequeue)
	http.HandleFunc("/q/view", viewQueue)

	log.Printf("Starting Server %s", fmt.Sprintf("%s:%s", "127.0.0.1", PORT))

	http.ListenAndServe(fmt.Sprintf("%s:%s", "127.0.0.1", PORT), nil)

}

// func pushStack(w http.ResponseWriter, req *http.Request) {

// }

// func popStack(w http.ResponseWriter, req *http.Request) {

// }

// add someone from the queue
func enqueue(w http.ResponseWriter, req *http.Request) {

	// read the json from the request
	var t interface{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)
	defer req.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Could not decode json body`))
		return
	}

	t.(map[string]interface{})["time"] = time.Now().Unix()

	// enqueue the item
	queue.Enqueue(t)

	// add position
	// return val
	RVal := map[string]interface{}{}
	RVal["Position"] = len(queue.S)
	RVal["echo"] = t

	//return status ok 200
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(RVal)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error enqueueing %v", t)))
	}
	w.Write(b)
	w.Header().Set("Content-Type", "application/json")

	log.Printf("Enqueued %v", t)
}

// remove someone from the queue
func dequeue(w http.ResponseWriter, req *http.Request) {
	var item interface{}
	item = queue.Dequeue()
	if item == nil {
		item = *new(interface{})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
	log.Printf("Dequeued %v", item)
	w.Header().Set("Content-Type", "application/json")
}

func viewQueue(w http.ResponseWriter, req *http.Request) {
	items := queue.View(nil, nil)
	w.WriteHeader(http.StatusOK)

	if items == nil {
		w.Write([]byte(`[]`))
		return
	}
	json.NewEncoder(w).Encode(items)
	w.Header().Set("Content-Type", "application/json")
}
