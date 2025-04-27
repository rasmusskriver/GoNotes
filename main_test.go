package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test for at hente noter
func TestGetNotesHandler(t *testing.T) {
	// Opret en tom note
	notes = []Note{
		{ID: 1, Content: "Note 1"},
		{ID: 2, Content: "Note 2"},
	}

	// Opret en testserver
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(getNotesHandler)

	// Kør handleren
	handler.ServeHTTP(rec, req)

	// Tjek statuskoden
	if rec.Code != http.StatusOK {
		t.Errorf("Forventede statuskode 200, men fik %d", rec.Code)
	}

	// Tjek om output er korrekt JSON
	var notesResponse []Note
	if err := json.NewDecoder(rec.Body).Decode(&notesResponse); err != nil {
		t.Fatalf("Kunne ikke dekode JSON: %v", err)
	}

	if len(notesResponse) != 2 {
		t.Errorf("Forventede 2 noter, men fik %d", len(notesResponse))
	}
}

// Test for at oprette en note
func TestCreateNoteHandler(t *testing.T) {
	// Opret en testnote som JSON
	newNote := Note{Content: "Test note"}
	body, err := json.Marshal(newNote)
	if err != nil {
		t.Fatalf("Kunne ikke marshale note: %v", err)
	}

	// Opret en testserver
	req := httptest.NewRequest(http.MethodPost, "/notes/create", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createNoteHandler)

	// Kør handleren
	handler.ServeHTTP(rec, req)

	// Tjek statuskoden
	if rec.Code != http.StatusOK {
		t.Errorf("Forventede statuskode 200, men fik %d", rec.Code)
	}

	// Tjek om den oprettede note er korrekt
	var createdNote Note
	if err := json.NewDecoder(rec.Body).Decode(&createdNote); err != nil {
		t.Fatalf("Kunne ikke dekode JSON: %v", err)
	}

	if createdNote.Content != "Test note" {
		t.Errorf("Forventede 'Test note', men fik '%s'", createdNote.Content)
	}

	// Også sikre på, at listen af noter er opdateret
	if len(notes) != 1 {
		t.Errorf("Forventede 1 note i notes-listen, men fandt %d", len(notes))
	}
}

// Test for at opdatere en note
func TestUpdateNoteHandler(t *testing.T) {
	// Først, opret en note
	notes = append(notes, Note{ID: 1, Content: "Old content"})

	// Opret en opdatering
	updatedNote := Note{ID: 1, Content: "Updated content"}
	body, err := json.Marshal(updatedNote)
	if err != nil {
		t.Fatalf("Kunne ikke marshale note: %v", err)
	}

	// Opret en testserver
	req := httptest.NewRequest(http.MethodPut, "/notes/update", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(updateNoteHandler)

	// Kør handleren
	handler.ServeHTTP(rec, req)

	// Tjek statuskoden
	if rec.Code != http.StatusOK {
		t.Errorf("Forventede statuskode 200, men fik %d", rec.Code)
	}

	// Tjek om den opdaterede note er korrekt
	var result Note
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("Kunne ikke dekode JSON: %v", err)
	}

	if result.Content != "Updated content" {
		t.Errorf("Forventede 'Updated content', men fik '%s'", result.Content)
	}
}

// Test for at slette en note
func TestDeleteNoteHandler(t *testing.T) {
	// Først, opret en note
	notes = append(notes, Note{ID: 1, Content: "Note to delete"})

	// Opret en testserver
	req := httptest.NewRequest(http.MethodDelete, "/notes/delete?id=1", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteNoteHandler)

	// Kør handleren
	handler.ServeHTTP(rec, req)

	// Tjek statuskoden
	if rec.Code != http.StatusNoContent {
		t.Errorf("Forventede statuskode 204, men fik %d", rec.Code)
	}

	// Tjek om listen af noter er opdateret
	if len(notes) != 0 {
		t.Errorf("Forventede 0 noter i notes-listen, men fandt %d", len(notes))
	}
}
