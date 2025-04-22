package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"

	"github.com/mileusna/masinilica/assets"
	"github.com/mileusna/srb"
	"golang.org/x/image/font/opentype"
)

var (
	masinilicaFont   *opentype.Font
	homepageTemplate *template.Template
)

func main() {

	var err error
	masinilicaFont, err = getFont("fonts/Masinci-Regular.ttf")
	if err != nil {
		log.Fatal("Greška prilikom učitavanja fonta", err)
	}

	homepageTemplate, err = template.ParseFS(assets.Templates, "templates/index.gohtml")
	if err != nil {
		log.Fatal("Template parse error:", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/masinilica.png", generateHandler)
	fs := http.FS(assets.Images)
	http.Handle("/images/", http.FileServer(fs))

	fmt.Println("Server started at http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

// homeHandler homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := homepageTemplate.Execute(w, nil); err != nil {
		http.Error(w, "Greška pri renderovanju templejta", http.StatusInternalServerError)
		log.Println("Template execute error:", err)
	}
}

// generateHadle handles http request to generate image
func generateHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Text nije unet", http.StatusInternalServerError)
		return
	}

	fontSize := r.FormValue("fs")
	align := r.FormValue("al")
	darkMode := r.FormValue("dm") == "1"

	png, err := CreateImage(srb.ToCyr(text), fontSize, align, darkMode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if r.FormValue("dl") == "1" {
		w.Header().Set("Content-Disposition", `attachment; filename="`+randFilename(12)+`.png"`)
	}
	w.Write(png)
}

// getFont loads font from assets
func getFont(fontPath string) (*opentype.Font, error) {
	fontBytes, err := assets.Fonts.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	return opentype.Parse(fontBytes)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randFilename(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
