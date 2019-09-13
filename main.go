package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jltorresm/peerserver/middleware"
	"github.com/jltorresm/peerserver/types"
	"github.com/nu7hatch/gouuid"
	"log"
	"net/http"
	"os"
)

const configFilename string = "config.json"

var topics = make(map[string]types.Topic)
var config = types.Config{}

func init() {
	config = getConfig()
}

func getConfig() types.Config {
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("[FATAL] config file %s does not exist.", configFilename))
	}

	file, _ := os.Open(configFilename)
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := types.Config{}
	err := decoder.Decode(&configuration)

	if err != nil {
		fmt.Println("error:", err)
	}

	return configuration
}

func PostTopic(writer http.ResponseWriter, request *http.Request) {
	// Generate a uuid for the new topic
	uuid, err := uuid.NewV4()

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Create the topic struct and save it in memory
	topics[uuid.String()] = types.Topic{
		Uuid:      uuid.String(),
		Content:   &types.Content{Title: "", Data: ""},
		Viewport:  &types.Viewport{X: 0, Y: 0},
		Selection: &types.Selections{},
	}

	// Send
	json.NewEncoder(writer).Encode(topics[uuid.String()])
}

func GetTopic(writer http.ResponseWriter, request *http.Request) {
	var response []types.Topic

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
	var content types.Content
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
	var viewport types.Viewport
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

func PutTopicSelection(writer http.ResponseWriter, request *http.Request) {
	var selections types.Selections
	params := mux.Vars(request)
	json.NewDecoder(request.Body).Decode(&selections)

	*(topics[params["uuid"]].Selection) = selections

	writer.WriteHeader(http.StatusNoContent)
}

func GetTopicSelection(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	json.NewEncoder(writer).Encode(topics[params["uuid"]].Selection)
}

func GetStatus(writer http.ResponseWriter, request *http.Request) {
	type sessionSummary struct {
		Uuid  string
		Title string
	}

	status := struct {
		Count    int
		Sessions []sessionSummary
	}{
		Count:    len(topics),
		Sessions: []sessionSummary{},
	}

	for _, t := range topics {
		status.Sessions = append(status.Sessions, sessionSummary{Uuid: t.Uuid, Title: t.Content.Title})
	}

	json.NewEncoder(writer).Encode(status)
}

// Main function
func main() {
	router := mux.NewRouter()

	// Set the routes
	router.HandleFunc("/status", GetStatus).Methods("GET")
	router.HandleFunc("/topic", PostTopic).Methods("POST")
	router.HandleFunc("/topic", GetTopic).Methods("GET")
	router.HandleFunc("/topic/{uuid}", DeleteTopic).Methods("DELETE")
	router.HandleFunc("/topic/{uuid}/content", PutTopicContent).Methods("PUT")
	router.HandleFunc("/topic/{uuid}/content", GetTopicContent).Methods("GET")
	router.HandleFunc("/topic/{uuid}/viewport", PutTopicViewport).Methods("PUT")
	router.HandleFunc("/topic/{uuid}/viewport", GetTopicViewport).Methods("GET")
	router.HandleFunc("/topic/{uuid}/selection", PutTopicSelection).Methods("PUT")
	router.HandleFunc("/topic/{uuid}/selection", GetTopicSelection).Methods("GET")

	// Add a few middlewares
	router.Use(middleware.HeaderNormalizerMiddleware)
	router.Use(middleware.LoggingMiddleware)

	log.Println(fmt.Sprintf("[INFO] Listening in %s", config.Server.Url))
	log.Fatal(http.ListenAndServe(config.Server.Url, router))
}
