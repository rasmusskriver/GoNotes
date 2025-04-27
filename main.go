package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Note struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

var notes = []Note{
	{ID: 1, Content: "Min første note"},
}

// Henter alle noter
func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// Opretter en ny note
func createNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newNote Note
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &newNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newNote.ID = getNextID()
	notes = append(notes, newNote)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}

// Sletter en note
func deleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ugyldigt ID", http.StatusBadRequest)
		return
	}

	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Note ikke fundet", http.StatusNotFound)
}

// Opdaterer en note
func updateNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ugyldigt ID", http.StatusBadRequest)
		return
	}

	var updatedNote Note
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &updatedNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, note := range notes {
		if note.ID == id {
			notes[i].Content = updatedNote.Content
			json.NewEncoder(w).Encode(notes[i])
			return
		}
	}
	http.Error(w, "Note ikke fundet", http.StatusNotFound)
}

func getNextID() int {
	if len(notes) == 0 {
		return 1
	}
	return notes[len(notes)-1].ID + 1
}

func main() {
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getNotes(w, r)
		} else if r.Method == http.MethodPost {
			createNote(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			deleteNote(w, r)
		case http.MethodPut:
			updateNote(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Serveren kører på http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
