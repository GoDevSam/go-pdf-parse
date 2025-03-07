/*
Create an API that scans a PDF and returns the file metadata. Each endpoint should return a json response.

This can use any language/framework of your choice and any popular libraries you wish to use. Database can be used, or just stored in memory on the program.

LLMs(ChatGPT, Deep seek, Claude, etc) can be used for research purposes, but not to complete the challenge.

/scan



The request to this endpoint should accept a PDF file.

The endpoint should generate an SHA256 hash synchronously, but process collecting the metadata asynchronously:

If the file type is a pdf, return a successful response to the client with the SHA256 hash in the response.

If the file type is not a PDF return a 400 error for invalid file type.

Collect the following metadata:

PDF version
Producer
Author
Created date
Updated date
Date scan was submitted in UTC time.
Save the results to a database or in memory map of your choice using the file SHA256 hash as the key.


Let the user lookup scan results using a SHA256 hash as the key.

The response will be the following.

SHA256 hash
PDF version
Producer
Author
Created date
Updated date
Submission date
If a record doesn't exist, handle this gracefully with a 404 error.

PDF parsing library: github.com/pdfcpu/pdfcpu




import "github.com/pdfcpu/pdfcpu/pkg/api"

ctx, err := api.ReadContext(bytes.NewReader(data), model.NewDefaultConfiguration())

err = api.ValidateContext(ctx)

ctx.XRefTable.

ctx.XRefTable.<> will contain the metadata

*/

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFMetadata struct {
	Hash         string `json:"hash"`
	Version      string `json:"pdf_version"`
	Producer     string `json:"producer"`
	Author       string `json:"author"`
	CreationDate string `json:"created_date"`
	ModifiedDate string `json:"updated_date"`
	UploadedAt   string `json:"uploaded_at"`
}

var dataStore sync.Map

func main() {
	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received request at /scan")
		handleUpload(w, r)
	})

	http.HandleFunc("/lookup", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received request at /lookup")
		handleLookup(w, r)
	})

	fmt.Println("server running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("internal server err:", err)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "error uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "application/pdf" {
		http.Error(w, "invalid file: only PDF allowed", http.StatusUnsupportedMediaType)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "couldn't reading file", http.StatusInternalServerError)
		return
	}

	hash := sha256.Sum256(content)
	hashString := hex.EncodeToString(hash[:])

	entry := PDFMetadata{
		Hash:       hashString,
		UploadedAt: time.Now().UTC().Format(time.RFC3339),
	}
	dataStore.Store(hashString, entry)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"hash": hashString})

	go processMetadata(content, hashString)
}

func processMetadata(content []byte, hash string) {
	tmpFile, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		fmt.Println("Temp file creation failed:", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(content)
	if err != nil {
		fmt.Println("Error writing to temp file:", err)
		return
	}
	tmpFile.Seek(0, io.SeekStart)

	conf := model.NewDefaultConfiguration()
	pdfInfo, err := api.PDFInfo(tmpFile, tmpFile.Name(), nil, conf)
	if err != nil {
		fmt.Println("err reading PDF metadata:", err)
		return
	}

	entryRaw, exists := dataStore.Load(hash)
	if !exists {
		fmt.Println("hash not found in dataStore:", hash)
		return
	}

	entry := entryRaw.(PDFMetadata)
	entry.Version = pdfInfo.Version
	entry.Producer = pdfInfo.Producer
	entry.Author = pdfInfo.Author
	entry.CreationDate = pdfInfo.CreationDate
	entry.ModifiedDate = pdfInfo.ModificationDate
	dataStore.Store(hash, entry)
}

func handleLookup(w http.ResponseWriter, r *http.Request) {
	queryHash := r.URL.Query().Get("hash")
	if queryHash == "" {
		http.Error(w, "hash not found", http.StatusBadRequest)
		return
	}

	entryRaw, exists := dataStore.Load(queryHash)
	if !exists {
		http.Error(w, "no data found", http.StatusNotFound)
		return
	}

	entry := entryRaw.(PDFMetadata)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}
