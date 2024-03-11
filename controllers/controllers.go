package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ramnath.1998/apica-app/handlers"
	"github.com/ramnath.1998/apica-app/models"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func GetCache(w http.ResponseWriter, r *http.Request) {
	outPutList := handlers.AppCache.GetTheListFromHeadToTail()

	outPutListJson, err := json.Marshal(outPutList)
	if err != nil {
		fmt.Printf("error in marshalling output map: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(string(outPutListJson)))
}

func UpdateCache(w http.ResponseWriter, r *http.Request) {

	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf(err.Error())
	}
	node := models.Node{}

	err = json.Unmarshal(content, &node)

	if err != nil {
		log.Fatalf(err.Error())
	}

	node = handlers.AppCache.UpdateCache(node.Key, node.Value, node.Expiration)

	json, err := json.Marshal(node)
	if err != nil {
		log.Fatalf(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(string(json)))

}
