package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

func convertToPDFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Retrieve text content from the request body
	textContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}

	// Create new PDF instance
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add a page
	pdf.AddPage()

	// Write text content to PDF
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, string(textContent), "", "", false)

	// Output PDF as response
	err = pdf.Output(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error generating PDF: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/convertToPDF", convertToPDFHandler)

	fmt.Println("Server is running on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}
