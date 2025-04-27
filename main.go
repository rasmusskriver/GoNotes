package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Note struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

var (
	notes      []Note
	notesMutex sync.Mutex
	filePath   = "notes.json"
)

func readNotes() error {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Filen findes ikke endnu, opret en tom liste
			notes = []Note{}
			return nil
		}
		return err
	}

	return json.Unmarshal(file, &notes)
}

func writeNotes() error {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	// Konverter notes til JSON
	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return err
	}

	// Skriv til filen
	return os.WriteFile(filePath, data, 0644)
}

func getNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Returner alle noter som JSON
	if r.Method != http.MethodGet {
		http.Error(w, "Metode ikke tilladt", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var newNote Note
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Sæt ID for den nye note
	newNote.ID = len(notes) + 1

	// Tilføj den nye note til listen
	notes = append(notes, newNote)

	// Gem opdaterede noter i filen
	if err := writeNotes(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Returner den oprettede note som JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newNote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var updatedNote Note
	if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find og opdater den eksisterende note
	for i, note := range notes {
		if note.ID == updatedNote.ID {
			notes[i].Content = updatedNote.Content
			break
		}
	}

	// Gem opdaterede noter i filen
	if err := writeNotes(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Returner den opdaterede note som JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedNote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Hent ID fra URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ugyldigt ID", http.StatusBadRequest)
		return
	}

	// Find og slet noten
	for i, note := range notes {
		if note.ID == id {
			if fmt.Sprintf("%d", note.ID) == id {
				notes = append(notes[:i], notes[i+1:]...)
				break
			}
		}
	}

	// Gem opdaterede noter i filen
	if err := writeNotes(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Læs noter fra fil ved serverstart
	if err := readNotes(); err != nil {
		fmt.Println("Fejl ved læsning af noter:", err)
	}
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Serverfejl:", err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/notes", getNotesHandler)
	http.HandleFunc("/notes/create", createNoteHandler)
	http.HandleFunc("/notes/update", updateNoteHandler)
	http.HandleFunc("/notes/delete", deleteNoteHandler)

	fmt.Println("Server kører på http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
