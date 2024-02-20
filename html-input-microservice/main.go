package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

// var htmlData string // Initialize htmlData

var pdfConversionURL = "http://localhost:8081/convertToPDF" // URL of the PDF Conversion Microservice

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Retrieve HTML content from the form
		htmlContent := r.FormValue("htmlContent")
		if htmlContent == "" {
			http.Error(w, "Empty HTML content", http.StatusBadRequest)
			return
		}

		// fmt.Println(htmlContent)

		// Send HTML data to the PDF Conversion Microservice
		response, err := http.Post(pdfConversionURL, "text/html", bytes.NewBufferString(htmlContent))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error sending POST request: %v", err), http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		// Read PDF content from response body
		pdfContent, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
			return
		}

		// Save PDF content to a file
		err = ioutil.WriteFile("output.pdf", pdfContent, 0644)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing PDF content to file: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Println("PDF file saved as output.pdf")

		// Provide a link for the user to download the PDF file
		w.Header().Set("Content-Disposition", "attachment; filename=output.pdf")
		w.Header().Set("Content-Type", "application/pdf")
		w.Write(pdfContent)

		http.Redirect(w, r, "/success", http.StatusSeeOther)
		return
	}

	// If the request method is GET, serve the HTML page with a form
	tmpl, err := template.ParseFiles("input.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing HTML template: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HTML data received and sent to PDF Conversion Microservice successfully.")
}

func main() {
	http.HandleFunc("/send", htmlHandler)
	http.HandleFunc("/success", successHandler)

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}
