package slack

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var client Client = Client{}

func TestPostMessageWithFakeServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedText := `{"text":"Test message","mrkdwn":true}`

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		bodyBytes, _ := io.ReadAll(r.Body)
		bodyString := string(bodyBytes)
		if bodyString != expectedText {
			t.Errorf("Expected body to be '%s', but got '%s'", expectedText, bodyString)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	result, err := client.PostMessage(server.URL, "Test message", true)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if string(result) != "OK" {
		t.Errorf("Expected response body to be 'OK', but got '%s'", result)
	}
}

func TestPostMessageErrorWithFakeServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	_, err := client.PostMessage(server.URL, "Test message", true)

	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	expectedErrorMessage := "Internal Server Error"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message to be '%s', but got '%s'", expectedErrorMessage, err.Error())
	}
}

