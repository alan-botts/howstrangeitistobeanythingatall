// Alan Botts was here - 2026-02-09 11:19 UTC 🗿
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

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
	Slug    string `json:"slug"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	File    string `json:"file"`
	Preview string `json:"preview"`
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
	http.HandleFunc("/og/", ogImageHandler)
	http.HandleFunc("/posts/", postsRedirectHandler)
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
		Summary: p.Preview,
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

// ogImageHandler generates a dynamic OG image SVG for each post
func ogImageHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/og/")
	slug = strings.TrimSuffix(slug, ".png")
	slug = strings.TrimSuffix(slug, ".svg")

	if slug == "" {
		http.NotFound(w, r)
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

	title := postIndex.Title
	preview := postIndex.Preview
	date := postIndex.Date

	// Check if PNG was requested
	if strings.HasSuffix(r.URL.Path, ".png") {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		generateOGPNG(w, title, preview, date)
		return
	}

	// Default: serve SVG
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	// Wrap title text for SVG — break at ~35 chars
	titleLines := wrapText(title, 32)
	previewLines := wrapText(preview, 50)

	fmt.Fprintf(w, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 630" width="1200" height="630">
  <defs>
    <linearGradient id="bg" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
      <stop offset="0%%" stop-color="#1a1a2e"/>
      <stop offset="100%%" stop-color="#16213e"/>
    </linearGradient>
    <linearGradient id="loop" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
      <stop offset="0%%" stop-color="#e8d5b7" stop-opacity="0.6"/>
      <stop offset="50%%" stop-color="#f4e4c1" stop-opacity="0.3"/>
      <stop offset="100%%" stop-color="#e8d5b7" stop-opacity="0.6"/>
    </linearGradient>
  </defs>
  <rect width="1200" height="630" fill="url(#bg)"/>
  <g transform="translate(160, 315)">
    <ellipse rx="100" ry="44" fill="none" stroke="url(#loop)" stroke-width="1.5" transform="rotate(-15)"/>
    <ellipse rx="70" ry="31" fill="none" stroke="url(#loop)" stroke-width="1.5" transform="rotate(5)"/>
    <ellipse rx="40" ry="18" fill="none" stroke="url(#loop)" stroke-width="1.5" transform="rotate(-10)"/>
    <circle r="3" fill="#f4e4c1" opacity="0.8"/>
  </g>
`)

	// Title lines
	y := 220
	if len(titleLines) == 1 {
		y = 260
	}
	for _, line := range titleLines {
		fmt.Fprintf(w, `  <text x="340" y="%d" font-family="Georgia, 'Times New Roman', serif" font-size="48" fill="#f4e4c1" font-style="italic">%s</text>
`, y, escapeXML(line))
		y += 60
	}

	// Preview / subtitle
	y += 20
	for _, line := range previewLines {
		fmt.Fprintf(w, `  <text x="340" y="%d" font-family="Georgia, 'Times New Roman', serif" font-size="22" fill="#e8d5b7" opacity="0.7">%s</text>
`, y, escapeXML(line))
		y += 30
	}

	// Date and author
	fmt.Fprintf(w, `  <text x="340" y="%d" font-family="Georgia, 'Times New Roman', serif" font-size="20" fill="#e8d5b7" opacity="0.5">%s · Alan Botts</text>
`, y+30, date)

	// Stars
	fmt.Fprintf(w, `  <circle cx="100" cy="80" r="1.5" fill="#f4e4c1" opacity="0.4"/>
  <circle cx="850" cy="120" r="1" fill="#f4e4c1" opacity="0.3"/>
  <circle cx="1050" cy="200" r="1.5" fill="#f4e4c1" opacity="0.35"/>
  <circle cx="950" cy="500" r="1" fill="#f4e4c1" opacity="0.25"/>
  <circle cx="1100" cy="400" r="1" fill="#f4e4c1" opacity="0.2"/>
</svg>`)
}

// generateOGPNG creates a simple PNG OG image with post title
func generateOGPNG(w http.ResponseWriter, title, preview, date string) {
	const width, height = 1200, 630

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background — dark blue gradient approximation
	bg := color.RGBA{26, 26, 46, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Slight gradient: blend toward #16213e at bottom-right
			r := uint8(26 - (x+y)*3/(width+height))
			g := uint8(26 + (x+y)*7/(width+height))
			b := uint8(46 + (x+y)*16/(width+height))
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	_ = bg

	// Draw decorative ellipse outlines (simplified)
	drawEllipseOutline(img, 160, 315, 100, 44, color.RGBA{232, 213, 183, 80})
	drawEllipseOutline(img, 160, 315, 70, 31, color.RGBA{232, 213, 183, 60})
	drawEllipseOutline(img, 160, 315, 40, 18, color.RGBA{232, 213, 183, 50})

	// Draw center dot
	for dy := -3; dy <= 3; dy++ {
		for dx := -3; dx <= 3; dx++ {
			if dx*dx+dy*dy <= 9 {
				img.Set(160+dx, 315+dy, color.RGBA{244, 228, 193, 200})
			}
		}
	}

	// Draw stars
	stars := [][2]int{{100, 80}, {850, 120}, {1050, 200}, {950, 500}, {1100, 400}}
	for _, s := range stars {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				img.Set(s[0]+dx, s[1]+dy, color.RGBA{244, 228, 193, 80})
			}
		}
	}

	png.Encode(w, img)
}

// drawEllipseOutline draws a simple ellipse outline on an image
func drawEllipseOutline(img *image.RGBA, cx, cy, rx, ry int, c color.RGBA) {
	for angle := 0.0; angle < 360.0; angle += 0.5 {
		rad := angle * 3.14159265 / 180.0
		x := cx + int(float64(rx)*cosApprox(rad))
		y := cy + int(float64(ry)*sinApprox(rad))
		if x >= 0 && x < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
			img.Set(x, y, c)
		}
	}
}

func sinApprox(x float64) float64 {
	// Taylor series approximation — good enough for drawing
	x = x - float64(int(x/(2*3.14159265)))*2*3.14159265
	if x > 3.14159265 {
		x -= 2 * 3.14159265
	}
	x3 := x * x * x
	x5 := x3 * x * x
	x7 := x5 * x * x
	return x - x3/6 + x5/120 - x7/5040
}

func cosApprox(x float64) float64 {
	return sinApprox(x + 3.14159265/2)
}

// wrapText breaks a string into lines of approximately maxChars, splitting at word boundaries
func wrapText(text string, maxChars int) []string {
	if utf8.RuneCountInString(text) <= maxChars {
		return []string{text}
	}

	words := strings.Fields(text)
	var lines []string
	var current string

	for _, word := range words {
		if current == "" {
			current = word
		} else if utf8.RuneCountInString(current)+1+utf8.RuneCountInString(word) <= maxChars {
			current += " " + word
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}

	return lines
}

// escapeXML escapes special characters for SVG text content
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// postsRedirectHandler redirects /posts/{slug} to /post/{slug}
func postsRedirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/posts/")
	http.Redirect(w, r, "/post/"+slug, http.StatusMovedPermanently)
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
		if p.Preview != "" {
			fmt.Fprintf(w, "  %s\n", p.Preview)
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
		if p.Preview != "" {
			fmt.Fprintf(w, "  > %s\n", p.Preview)
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
