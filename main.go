package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gamedevja/backend/git"
)

const (
	FILE_LIMIT_BYTES = 100000
	TEXT_LIMIT       = 60
)

var templates = template.Must(template.ParseFiles("tmpl/upload.html"))

func main() {
	http.HandleFunc("/", index)
	fmt.Println("listening...")
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}

func gitpush(blobs []git.Blob) error {
	loc, _ := time.LoadLocation("Japan")
	message := "remote: " + time.Now().In(loc).String()
	if err := git.Push(blobs, &message); err != nil {
		return err
	}
	return nil
}

func display(w http.ResponseWriter, tmpl string, data interface{}) {
	d := map[string]interface{}{
		"Message":   data,
		"Limit":     fmt.Sprintf("(~%dkB)", FILE_LIMIT_BYTES/1000),
		"TextLimit": fmt.Sprintf("(~%d)", TEXT_LIMIT),
	}
	templates.ExecuteTemplate(w, tmpl+".html", d)
}

func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		display(w, "upload", nil)

	case "POST":
		err := r.ParseMultipartForm(FILE_LIMIT_BYTES)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fh := r.MultipartForm.File["file"][0]
		value := r.MultipartForm.Value["text"][0]

		if value == "" || len(value) > 60 {
			http.Error(w, "inalid text", http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(fh.Filename)
		if ext != ".png" && ext != ".jpg" && ext != ".svg" && ext != ".gif" && ext != ".mp3" && ext != ".ogg" {
			http.Error(w, "invalid file", http.StatusInternalServerError)
			return
		}
		f, err := fh.Open()
		defer f.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var b bytes.Buffer
		if _, err := io.Copy(&b, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		content := b.String()

		var blobs []git.Blob

		now := fmt.Sprintf("%d", time.Now().UnixNano())

		file := new(git.Blob)
		file.Blob.Content = &content
		file.Blob.Encoding = &git.BlobBase64Encode
		var testpath string
		if os.Getenv("PRODUCTION") != "true" {
			testpath = "test/"
		}
		file.Path = "assets/" + testpath + "images/" + now + ext
		blobs = append(blobs, *file)

		j, err := json.Marshal(map[string]string{
			"filepath": file.Path,
			"text":     value,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		value = string(j)
		text := new(git.Blob)
		text.Blob.Content = &value
		text.Blob.Encoding = &git.BlobUtf8Encode
		text.Path = "assets/" + testpath + "entries/" + now + ".json"
		blobs = append(blobs, *text)

		if err = gitpush(blobs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		display(w, "upload", "Upload successful.")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
