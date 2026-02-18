// Alan Botts was here - 2026-02-09 11:19 UTC 🗿
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

const siteURL = "https://howstrangeitistobeanythingatall.com"
const verbatimURL = "https://strangerloops.com"

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
	AllPosts []Post
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

	// Use a custom handler that routes based on Host header
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/archive", archiveHandler)
	http.HandleFunc("/feed.xml", rssHandler)
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

// rootHandler routes based on Host header
func rootHandler(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	// Strip port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	if host == "strangerloops.com" || host == "www.strangerloops.com" {
		verbatimHandler(w, r)
		return
	}

	homeHandler(w, r)
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

	// Parse date - try index.json first, fallback to filename
	date, err := time.Parse("2006-01-02", p.Date)
	if err != nil || p.Date == "" {
		// Try extracting from filename (e.g., "2026-02-03-title.md")
		filename := filepath.Base(p.File)
		if len(filename) >= 10 {
			date, _ = time.Parse("2006-01-02", filename[:10])
		}
	}

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

	index, allPosts, err := loadPosts()
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		log.Printf("Error loading posts: %v", err)
		return
	}

	displayPosts := allPosts
	hasMore := len(allPosts) > recentPostsLimit
	if hasMore {
		displayPosts = allPosts[:recentPostsLimit]
	}

	data := HomeData{
		Blog:     index,
		Posts:    displayPosts,
		AllPosts: allPosts,
		HasMore:  hasMore,
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

	index, allPosts, err := loadPosts()
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
		Blog     BlogIndex
		Post     Post
		AllPosts []Post
	}{
		Blog:     index,
		Post:     post,
		AllPosts: allPosts,
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
		fmt.Fprintf(w, "- [%s](%s) (%s)\n", p.Title, siteURL+"/post/"+slug, p.Date)
		if p.Summary != "" {
			fmt.Fprintf(w, "  %s\n", p.Summary)
		}
		fmt.Fprintln(w)
	}
}

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	index, posts, err := loadPosts()
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	var items []RSSItem
	for _, post := range posts {
		link := siteURL + "/post/" + post.Slug
		items = append(items, RSSItem{
			Title:       post.Title,
			Link:        link,
			Description: string(post.Content),
			PubDate:     post.Date.Format(time.RFC1123Z),
			GUID:        link,
		})
	}

	rss := RSS{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       index.Title,
			Link:        siteURL,
			Description: index.Description,
			Items:       items,
		},
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	enc.Encode(rss)
}

// ============================================================================
// Verbatim subdomain handlers - serves raw markdown for agent consumption
// ============================================================================

// VerbatimIndex represents optional metadata for files in a directory
type VerbatimIndex struct {
	Files []VerbatimFileEntry `json:"files"`
}

// VerbatimFileEntry represents metadata for a single file
type VerbatimFileEntry struct {
	Filename    string `json:"filename"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Published   string `json:"published"`    // ISO 8601 timestamp (e.g., 2026-02-01T14:30:00Z)
	LastUpdated string `json:"last_updated"` // ISO 8601 timestamp (e.g., 2026-02-01T14:30:00Z)
}

// verbatimHandler serves the verbatim subdomain
func verbatimHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	path = strings.TrimSuffix(path, "/")

	// Root path - serve index
	if path == "" {
		verbatimIndexHandler(w, r)
		return
	}

	// Virtual "blog" directory - pulls from main blog content
	if path == "blog" {
		verbatimBlogIndexHandler(w, r)
		return
	}
	if strings.HasPrefix(path, "blog/") {
		verbatimBlogFileHandler(w, r, strings.TrimPrefix(path, "blog/"))
		return
	}

	// Serve the requested file
	verbatimFileHandler(w, r, path)
}

// verbatimIndexHandler serves the homepage - either index.md if it exists, or auto-generated
func verbatimIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")

	// If index.md exists, serve it directly (gives Alan full editorial control)
	indexPath := "content/verbatim/index.md"
	if data, err := os.ReadFile(indexPath); err == nil {
		w.Write(data)
		return
	}

	// Fallback: auto-generate index
	fmt.Fprintln(w, "# Verbatim")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Raw markdown content curated by Alan Botts. Navigate via the links below.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w)

	// Load optional index.json for root
	indexMeta := loadVerbatimIndex("content/verbatim")

	// Collect all markdown files
	files, dirs := collectVerbatimEntries("content/verbatim", "")

	// Add virtual "blog" directory
	dirs = append([]string{"blog"}, dirs...)

	// Render directories first
	if len(dirs) > 0 {
		fmt.Fprintln(w, "## Directories")
		fmt.Fprintln(w)
		for _, dir := range dirs {
			fmt.Fprintf(w, "- [%s/](%s/%s/)\n", dir, verbatimURL, dir)
		}
		fmt.Fprintln(w)
	}

	// Render files
	if len(files) > 0 {
		fmt.Fprintln(w, "## Files")
		fmt.Fprintln(w)
		for _, file := range files {
			renderVerbatimFileLink(w, file, "", indexMeta)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "See [README.md](%s/README.md) for content guidelines.\n", verbatimURL)
}

// verbatimFileHandler serves individual files or directory indexes
func verbatimFileHandler(w http.ResponseWriter, r *http.Request, path string) {
	fullPath := filepath.Join("content/verbatim", path)

	// Check if it's a directory
	info, err := os.Stat(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if info.IsDir() {
		// Serve directory index
		verbatimDirIndexHandler(w, r, path)
		return
	}

	// Only serve .md files
	if !strings.HasSuffix(path, ".md") {
		http.NotFound(w, r)
		return
	}

	// Serve the raw markdown file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Write(data)
}

// verbatimDirIndexHandler serves an index for a subdirectory
func verbatimDirIndexHandler(w http.ResponseWriter, r *http.Request, dirPath string) {
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")

	fmt.Fprintf(w, "# %s\n\n", dirPath)

	fullPath := filepath.Join("content/verbatim", dirPath)
	indexMeta := loadVerbatimIndex(fullPath)

	files, dirs := collectVerbatimEntries(fullPath, dirPath)

	// Parent link
	parentPath := filepath.Dir(dirPath)
	if parentPath == "." {
		fmt.Fprintf(w, "[← Back to root](%s/)\n\n", verbatimURL)
	} else {
		fmt.Fprintf(w, "[← Back to %s](%s/%s/)\n\n", parentPath, verbatimURL, parentPath)
	}

	fmt.Fprintln(w, "---")
	fmt.Fprintln(w)

	// Render subdirectories
	if len(dirs) > 0 {
		fmt.Fprintln(w, "## Directories")
		fmt.Fprintln(w)
		for _, dir := range dirs {
			subDirPath := filepath.Join(dirPath, dir)
			fmt.Fprintf(w, "- [%s/](%s/%s/)\n", dir, verbatimURL, subDirPath)
		}
		fmt.Fprintln(w)
	}

	// Render files
	if len(files) > 0 {
		fmt.Fprintln(w, "## Files")
		fmt.Fprintln(w)
		for _, file := range files {
			renderVerbatimFileLink(w, file, dirPath, indexMeta)
		}
	}
}

// loadVerbatimIndex loads index.json from a directory if it exists
func loadVerbatimIndex(dirPath string) map[string]VerbatimFileEntry {
	result := make(map[string]VerbatimFileEntry)

	indexPath := filepath.Join(dirPath, "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return result
	}

	var index VerbatimIndex
	if err := json.Unmarshal(data, &index); err != nil {
		log.Printf("Error parsing verbatim index at %s: %v", indexPath, err)
		return result
	}

	for _, entry := range index.Files {
		result[entry.Filename] = entry
	}

	return result
}

// collectVerbatimEntries returns markdown files and subdirectories in a path
func collectVerbatimEntries(fullPath, relativePath string) (files []string, dirs []string) {
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, nil
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files, index.json, and index.md
		if strings.HasPrefix(name, ".") || name == "index.json" || name == "index.md" {
			continue
		}

		if entry.IsDir() {
			dirs = append(dirs, name)
		} else if strings.HasSuffix(name, ".md") {
			files = append(files, name)
		}
	}

	sort.Strings(files)
	sort.Strings(dirs)

	return files, dirs
}

// renderVerbatimFileLink writes a markdown link for a file, with optional metadata
func renderVerbatimFileLink(w http.ResponseWriter, filename, dirPath string, indexMeta map[string]VerbatimFileEntry) {
	entry, hasMetadata := indexMeta[filename]

	// Build the full URL path
	urlPath := filename
	if dirPath != "" {
		urlPath = dirPath + "/" + filename
	}

	if hasMetadata && entry.Title != "" {
		// Rich link with metadata
		fmt.Fprintf(w, "- [%s](%s/%s)", entry.Title, verbatimURL, urlPath)

		var meta []string
		if entry.Published != "" {
			meta = append(meta, fmt.Sprintf("published: %s", entry.Published))
		}
		if entry.LastUpdated != "" {
			meta = append(meta, fmt.Sprintf("updated: %s", entry.LastUpdated))
		}

		if len(meta) > 0 {
			fmt.Fprintf(w, " (%s)", strings.Join(meta, ", "))
		}
		fmt.Fprintln(w)

		if entry.Summary != "" {
			fmt.Fprintf(w, "  > %s\n", entry.Summary)
		}
	} else {
		// Simple link
		fmt.Fprintf(w, "- [%s](%s/%s)\n", filename, verbatimURL, urlPath)
	}
}

// walkVerbatimDir recursively collects all markdown files
func walkVerbatimDir(basePath string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".md") {
			relPath, _ := filepath.Rel(basePath, path)
			files = append(files, relPath)
		}

		return nil
	})

	return files, err
}

// ============================================================================
// Virtual blog directory - serves main blog content via verbatim
// ============================================================================

// verbatimBlogIndexHandler serves the index for the virtual blog directory
func verbatimBlogIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")

	index, err := loadIndex()
	if err != nil {
		http.Error(w, "Error loading blog index", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "# blog")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "> %s\n\n", index.Description)
	fmt.Fprintf(w, "[← Back to root](%s/)\n\n", verbatimURL)
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "## Posts")
	fmt.Fprintln(w)

	for _, p := range index.Posts {
		filename := filepath.Base(p.File)
		fmt.Fprintf(w, "- [%s](%s/blog/%s)", p.Title, verbatimURL, filename)
		if p.Date != "" {
			fmt.Fprintf(w, " (%s)", p.Date)
		}
		fmt.Fprintln(w)
		if p.Summary != "" {
			fmt.Fprintf(w, "  > %s\n", p.Summary)
		}
	}
}

// verbatimBlogFileHandler serves individual blog post files as raw markdown
func verbatimBlogFileHandler(w http.ResponseWriter, r *http.Request, filename string) {
	// Only serve .md files
	if !strings.HasSuffix(filename, ".md") {
		http.NotFound(w, r)
		return
	}

	// Read from content/posts/
	fullPath := filepath.Join("content/posts", filename)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Write(data)
}
