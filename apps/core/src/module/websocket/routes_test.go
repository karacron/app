package websocket

import (
	"encoding/json"
	"testing"
)

func TestExtractPromptFromString(t *testing.T) {
	raw := json.RawMessage(`"hola kara"`)
	got := extractPrompt(raw)
	if got != "hola kara" {
		t.Fatalf("se esperaba prompt 'hola kara', se obtuvo '%s'", got)
	}
}

func TestExtractPromptFromObject(t *testing.T) {
	raw := json.RawMessage(`{"message":"  prueba ws  "}`)
	got := extractPrompt(raw)
	if got != "prueba ws" {
		t.Fatalf("se esperaba prompt 'prueba ws', se obtuvo '%s'", got)
	}
}

func TestExtractPromptEmptyOnInvalidPayload(t *testing.T) {
	raw := json.RawMessage(`{"other":"value"}`)
	got := extractPrompt(raw)
	if got != "" {
		t.Fatalf("se esperaba prompt vacio, se obtuvo '%s'", got)
	}
}
