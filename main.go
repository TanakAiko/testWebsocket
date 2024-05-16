package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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
	ReadBufferSize:  1024 * 1014,
	WriteBufferSize: 1024 * 1014,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	// Mise à niveau de la connexion HTTP en WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	for {

		messageType, data, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading msg: %v", err)
			return
		}

		fmt.Printf("Received message: %s\n", data)

		if err := ws.WriteMessage(messageType, data); err != nil {
			fmt.Println(err)
			return
		}

		var msg ScoreData

		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Printf("Error when Unmarshal data: %v", err)
			return
		}

		fmt.Println(msg, "HERREEE msg")

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

func home(w http.ResponseWriter, r *http.Request) {
	/* if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	} */

	if r.Method != http.MethodGet {
		http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error server", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error server", http.StatusInternalServerError)
		return
	}
}

func writeJSONToFile(data ScoreData) error {
	// Lecture du fichier existant
	content, err := os.ReadFile("score.json")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Désérialisation des données existantes si le fichier n'est pas vide
	var scores []ScoreData
	if len(content) > 0 {
		if err := json.Unmarshal(content, &scores); err != nil {
			return err
		}
	}

	// Ajout des nouvelles données
	scores = append(scores, data)

	// Sérialisation des données mises à jour
	content, err = json.MarshalIndent(scores, "", " ")
	if err != nil {
		return err
	}

	// Ecriture des données mises à jour dans le fichier
	if err := os.WriteFile("score.json", content, 0644); err != nil {
		return err
	}

	log.Println("Data successfully written to score.json")
	return nil
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/ws", handleWebsocket)

	// Gestion des fichiers statiques
	fs := http.FileServer(http.Dir("."))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("WebSocket server starting on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
