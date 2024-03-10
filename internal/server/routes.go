package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"goAuth/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)
	r.Get("/auth/{provider}", s.beginAuthProvideCallback)
	r.Get("/logout/{provider}", s.logOutProvider)

	r.Post("/upload", s.uploadFile)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	// jsonResp, _ := json.Marshal(s.db.Health())
	// _, _ = w.Write(jsonResp)
	jsonText := `{"status": "ok"}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonText))
	// jsonResp, _ := json.Marshal(jsonText)
	// w.Write(jsonResp)
	// w.Write([]byte("OK"))
}

func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Println(user)
	postJsonBytes, err := json.MarshalIndent(user, "", "    ")
	// postJsonBytes, err := JSONMarshal(postJson, "", "    ")
	if err != nil {
		fmt.Fprintln(w, err)
	}
	fmt.Println(string(postJsonBytes))
	http.Redirect(w, r, "http://localhost:3000/movies/dashboard", http.StatusFound)

}

func (s *Server) beginAuthProvideCallback(w http.ResponseWriter, r *http.Request) {

	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	gothic.BeginAuthHandler(w, r)
}

func (s *Server) logOutProvider(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	gothic.Logout(w, r)
	w.Header().Set("Location", "http://localhost:3000/movies/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func JSONMarshal(t interface{}, prefix, indent string) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func (s *Server) uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	// Maximum upload of 10 MB files
	// r.ParseMultipartForm(10 << 20)

	const maxUploadSize = 55<<20 + 512 // 55 MB + 512 bytes

	// Add this line on top
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		// Handle errors consistently and provide informative error messages
		http.Error(w, "File exceeds maximum upload size or is invalid", http.StatusBadRequest)
		return
	}
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
		jsonText := fmt.Sprintf(`{"message": "Error uploading file: %v"}`, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonText))
		return
	}

	fileBuffer := make([]byte, handler.Size)
	_, err = file.Read(fileBuffer)
	if err != nil {
		// Handle error
	}

	defer file.Close()

	// Process the uploaded file
	err = processUploadedFile(fileBuffer, handler.Filename, handler.Size, handler.Header.Get("Content-Type"))
	if err != nil {
		// Handle error
	}

	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Indicate successful upload with a concise JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "file uploaded successfully"}`))
}

func processUploadedFile(fileBuffer []byte, fileName string, fileSize int64, contentType string) error {
	fmt.Println("Processing file")
	fmt.Printf("Uploaded File: %+v\n", fileName)
	fmt.Printf("File Size: %+v\n", fileSize)
	fmt.Printf("Content Type: %+v\n", contentType)

	storage.AzureFileUploadByBytes(fileBuffer, fileSize, fileName)
	return nil
}
