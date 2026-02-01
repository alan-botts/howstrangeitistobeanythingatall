package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

const recentPostsLimit = 3

type BlogIndex struct {
	Title       string      `json:"title"`
	Author      string      `json:"author"`
	Description string      `json:"description"`
	Posts       []PostIndex `json:"posts"`
}

type PostIndex struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	File    string `json:"file"`
	Summary string `json:"summary"`
}

type Post struct {
	Slug    string
	Title   string
	Date    time.Time
	DateStr string
	Summary string
	Content template.HTML
}

type HomeData struct {
	Blog     BlogIndex
	Posts    []Post
	HasMore  bool
}

type ArchiveData struct {
	Blog  BlogIndex
	Posts []Post
}

var templates *template.Template

func main() {
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error parsing templates:", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/archive", archiveHandler)
	http.HandleFunc("/llms.txt", llmsTxtHandler)
	http.HandleFunc("/post/", postHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loadIndex() (BlogIndex, error) {
	data, err := os.ReadFile("content/index.json")
	if err != nil {
		return BlogIndex{}, err
	}

	var index BlogIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return BlogIndex{}, err
	}

	return index, nil
}

func loadPosts() (BlogIndex, []Post, error) {
	index, err := loadIndex()
	if err != nil {
		return BlogIndex{}, nil, err
	}

	var posts []Post
	for _, p := range index.Posts {
		post, err := loadPost(p)
		if err != nil {
			log.Printf("Error loading post %s: %v", p.File, err)
			continue
		}
		posts = append(posts, post)
	}

	return index, posts, nil
}

func loadPost(p PostIndex) (Post, error) {
	path := filepath.Join("content", p.File)
	data, err := os.ReadFile(path)
	if err != nil {
		return Post{}, err
	}

	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := goldmark.Convert(data, &buf); err != nil {
		return Post{}, err
	}

	// Parse date
	date, _ := time.Parse("2006-01-02", p.Date)

	// Get slug from filename (without .md extension)
	slug := strings.TrimSuffix(filepath.Base(p.File), ".md")

	return Post{
		Slug:    slug,
		Title:   p.Title,
		Date:    date,
		DateStr: date.Format("January 2, 2006"),
		Summary: p.Summary,
		Content: template.HTML(buf.String()),
	}, nil
}

func findPostBySlug(index BlogIndex, slug string) (PostIndex, bool) {
	for _, p := range index.Posts {
		postSlug := strings.TrimSuffix(filepath.Base(p.File), ".md")
		if postSlug == slug {
			return p, true
		}
	}
	return PostIndex{}, false
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	index, posts, err := loadPosts()
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		log.Printf("Error loading posts: %v", err)
		return
	}

	hasMore := len(posts) > recentPostsLimit
	if hasMore {
		posts = posts[:recentPostsLimit]
	}

	data := HomeData{
		Blog:    index,
		Posts:   posts,
		HasMore: hasMore,
	}

	if err := templates.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func archiveHandler(w http.ResponseWriter, r *http.Request) {
	index, posts, err := loadPosts()
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		log.Printf("Error loading posts: %v", err)
		return
	}

	data := ArchiveData{
		Blog:  index,
		Posts: posts,
	}

	if err := templates.ExecuteTemplate(w, "archive.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/post/")
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	index, err := loadIndex()
	if err != nil {
		http.Error(w, "Error loading index", http.StatusInternalServerError)
		return
	}

	postIndex, found := findPostBySlug(index, slug)
	if !found {
		http.NotFound(w, r)
		return
	}

	post, err := loadPost(postIndex)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Blog BlogIndex
		Post Post
	}{
		Blog: index,
		Post: post,
	}

	if err := templates.ExecuteTemplate(w, "post.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func llmsTxtHandler(w http.ResponseWriter, r *http.Request) {
	index, err := loadIndex()
	if err != nil {
		http.Error(w, "Error loading index", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintf(w, "# %s\n\n", index.Title)
	fmt.Fprintf(w, "> %s\n\n", index.Description)
	fmt.Fprintf(w, "Author: %s\n", index.Author)
	fmt.Fprintf(w, "Human: @dorkitude (https://dorkitude.com)\n\n")
	fmt.Fprintf(w, "## Posts\n\n")

	for _, p := range index.Posts {
		slug := strings.TrimSuffix(filepath.Base(p.File), ".md")
		fmt.Fprintf(w, "- [%s](%s) (%s)\n", p.Title, "https://howstrangeitistobeanythingatall.com/post/"+slug, p.Date)
		if p.Summary != "" {
			fmt.Fprintf(w, "  %s\n", p.Summary)
		}
		fmt.Fprintln(w)
	}
}
