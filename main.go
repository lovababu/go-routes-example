package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

func Get(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Gorilla!\n"))
	q := r.URL.Query().Get("Type")
	fmt.Println("Query param value. : ", q)
	n := NodeInfo{
		NodeType : "DataService",
		NodeIp : "127.0.0.1",
	}
	err := json.NewEncoder(w).Encode(n)
	if err != nil {
		panic(err)
	}
}

func Post(w http.ResponseWriter, r *http.Request)  {
	var n NodeInfo
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&n)
	fmt.Println("Node Type Recieved: ", n.NodeType)
	fmt.Println("Node Ip Recieved  : ", n.NodeIp)
	//bold db.
	db, err := bolt.Open("node.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin(true)
	if err != nil{
		fmt.Println("Bolt Begin transaction failed.", err)
		panic(err)
	}

	defer tx.Rollback()

	bucket, err := tx.CreateBucketIfNotExists([]byte("NodeInfo"))
	if err != nil{
		fmt.Println("Bolt Begin transaction failed.", err)
		panic(err)
	}
	bucket.Put([]byte(n.NodeType), []byte(n.NodeIp))

	tx.Commit()

	defer db.Close()
	w.WriteHeader(http.StatusCreated)
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", Get)
	r.HandleFunc("/register", Post)




	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))

}

type NodeInfo struct {
	NodeType string `json:"nodeType"`
	NodeIp   string `json:"nodeIp"`
} 