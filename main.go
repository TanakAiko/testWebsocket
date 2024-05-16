package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
		//log.Fatalf("Failed to set up WebSocket upgrade: %v", err)
		fmt.Println(err)
		return
	}
	defer ws.Close()

	for {

		messageType, data, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading msg: %v", err)
			fmt.Println("aaaaaaaaaaaa => messageType : ", messageType)
			return
		}

		if ws.WriteMessage(messageType, data); err != nil {
			fmt.Println(err)
			return
		}

		var msg ScoreData

		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Printf("Error when Unmarshal data: %v", err)
			return
		}

		// Validation des données reçues
		if msg.Name == "" || msg.Score < 0 || msg.Time < 0 {
			fmt.Println("Invalid data received !")
			return
		}

		if err := writeJSONToFile(msg); err != nil {
			fmt.Printf("Error writing JSON to file: %v", err)
			return
		}
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
		fmt.Printf("Error writing JSON to file: %v", err)
		return err
	}

	fmt.Println("Data successfully written to score.json")
	return nil
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	/* http.HandleFunc("/", home) */
	fmt.Println("WebSocket server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
