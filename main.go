package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/websocket"
)

type ScoreData struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Time  int    `json:"time"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Mise à niveau de la connexion HTTP en WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to set up WebSocket upgrade: %v", err)
	}
	//defer ws.Close()

	for {
		// Lecture des données JSON envoyées par le client
		var msg ScoreData
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			println("aaaaaaaaaaaa ")
			break
		}
		println("bbbbbbbbbbbbbbbbbb")

		// Validation des données reçues
		if msg.Name == "" || msg.Score < 0 || msg.Time < 0 {
			log.Println("Invalid data received : ", msg)
			break
		}

		// Écriture des données dans le fichier
		if err := writeJSONToFile(msg); err != nil {
			log.Printf("Error writing JSON to file: %v", err)
			break
		}

		fmt.Printf("Received: %+v\n", msg)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func writeJSONToFile(data ScoreData) error {
	var scores []ScoreData

	// Lecture du fichier existant
	_, err := os.Stat("score.json")
	if !os.IsNotExist(err) {
		content, err := os.ReadFile("score.json")
		if err != nil {
			return err
		}

		// Désérialisation des données existantes si le fichier n'est pas vide
		if len(content) > 0 {
			if err := json.Unmarshal(content, &scores); err != nil {
				return err
			}
		}
	}

	// Ajout des nouvelles données
	scores = append(scores, data)

	// Sérialisation des données mises à jour
	content, err := json.MarshalIndent(scores, "", " ")
	if err != nil {
		return err
	}

	// Ecriture des données mises à jour dans le fichier
	if err := os.WriteFile("score.json", content, 0644); err != nil {
		log.Printf("Error writing JSON to file: %v", err)
		return err
	}

	log.Println("Data successfully written to score.json")
	return nil
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/", home)
	fmt.Println("WebSocket server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
