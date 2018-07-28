package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
)

type Topic struct {
	Uuid     string    `json:"uuid"`
	Content  *Content  `json:"content,omitempty"`
	Viewport *Viewport `json:"viewport,omitempty"`
}

type Content struct {
	Title string `json:"title"`
	Data string `json:"data"`
}

type Viewport struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

var topics = make(map[string]Topic)

func PostTopic(writer http.ResponseWriter, request *http.Request) {
	// Generate a uuid for the new topic
	uuid, err := uuid.NewV4()

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Create the topic struct and save it in memory
	topics[uuid.String()] = Topic{Uuid: uuid.String(), Content: &Content{Title: "", Data: ""}, Viewport: &Viewport{X: 0, Y: 0}}

	// Send
	json.NewEncoder(writer).Encode(topics[uuid.String()])
}

func GetTopic(writer http.ResponseWriter, request *http.Request) {
	var response []Topic

	for _, topic := range topics {
		response = append(response, topic)
	}

	json.NewEncoder(writer).Encode(response)
}

func DeleteTopic(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	delete(topics, params["uuid"])
	writer.WriteHeader(http.StatusNoContent)
}

func PutTopicContent(writer http.ResponseWriter, request *http.Request) {
	var content Content
	params := mux.Vars(request)
	json.NewDecoder(request.Body).Decode(&content)

	topics[params["uuid"]].Content.Data = content.Data
	topics[params["uuid"]].Content.Title = content.Title

	writer.WriteHeader(http.StatusNoContent)
}

func GetTopicContent(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	json.NewEncoder(writer).Encode(topics[params["uuid"]].Content)
}

func PutTopicViewport(writer http.ResponseWriter, request *http.Request) {
	var viewport Viewport
	params := mux.Vars(request)
	json.NewDecoder(request.Body).Decode(&viewport)

	topics[params["uuid"]].Viewport.X = viewport.X
	topics[params["uuid"]].Viewport.Y = viewport.Y

	writer.WriteHeader(http.StatusNoContent)
}

func GetTopicViewport(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	json.NewEncoder(writer).Encode(topics[params["uuid"]].Viewport)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		log.Printf("[%s] %s", request.Method, request.RequestURI)

		next.ServeHTTP(writer, request)
	})
}

func headerNormalizerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/json; charset=utf-8")
		next.ServeHTTP(writer, request)
	})
}

// Main function
func main() {
	topics["030fffa0-8ebb-11e8-b26b-926d69cab634"] = Topic{Uuid: "030fffa0-8ebb-11e8-b26b-926d69cab634", Content: &Content{Title: "Test Topic", Data: "This is the data of the remote buffer"}, Viewport: &Viewport{X: 1, Y: 1}}

	router := mux.NewRouter()

	// Set the routes
	router.HandleFunc("/topic", PostTopic).Methods("POST")
	router.HandleFunc("/topic", GetTopic).Methods("GET")
	router.HandleFunc("/topic/{uuid}", DeleteTopic).Methods("DELETE")
	router.HandleFunc("/topic/{uuid}/content", PutTopicContent).Methods("PUT")
	router.HandleFunc("/topic/{uuid}/content", GetTopicContent).Methods("GET")
	router.HandleFunc("/topic/{uuid}/viewport", PutTopicViewport).Methods("PUT")
	router.HandleFunc("/topic/{uuid}/viewport", GetTopicViewport).Methods("GET")

	// Add a few middlewares
	router.Use(headerNormalizerMiddleware)
	router.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":8000", router))
}
