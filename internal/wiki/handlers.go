package wiki

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

const (
	savedPagesPath = "data/pages"
	templatesPath  = "web/templates"
	staticPath     = "web/static"
)

type Page struct {
	Title         string
	InternalTitle string
	Body          []byte
}

var funcMap = template.FuncMap{"nl2br": nl2br}

var templates = template.Must(
	template.New("").Funcs(funcMap).ParseFiles(
		filepath.Join(templatesPath, "edit.html"),
		filepath.Join(templatesPath, "view.html"),
		filepath.Join(templatesPath, "index.html"),
		filepath.Join(templatesPath, "add_new_page.html"),
	),
)

func NewServer() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/view/", viewHandler)
	mux.HandleFunc("/edit/", editHandler)
	mux.HandleFunc("/save/", saveHandler)
	mux.HandleFunc("/add_new_page", addNewPageHandler)
	mux.HandleFunc("/internal_add_new_page", internalAddNewPageHandler)
	mux.HandleFunc("/index", indexHandler)
	mux.HandleFunc("/", rootHandler)

	return mux
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index", http.StatusFound)
}

type IndexPageListing struct {
	Title         string
	InternalTitle string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(savedPagesPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			files = nil
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var pageListings []IndexPageListing
	for _, file := range files {
		pageListing := IndexPageListing{
			InternalTitle: file.Name(),
			Title:         internalTitleToTitle(file.Name()),
		}
		pageListings = append(pageListings, pageListing)
	}

	templateErr := templates.ExecuteTemplate(w, "index.html", pageListings)
	if templateErr != nil {
		http.Error(w, templateErr.Error(), http.StatusInternalServerError)
	}
}

// as-is renders the add new page template
func addNewPageHandler(w http.ResponseWriter, r *http.Request) {
	templateErr := templates.ExecuteTemplate(w, "add_new_page.html", nil)
	if templateErr != nil {
		http.Error(w, templateErr.Error(), http.StatusInternalServerError)
	}
}

// called by the add new page form
// makes a new blank page
func isValidPageTitle(title string) bool {
	// Regex to allow alphanumeric and spaces, adjust your rules as needed
	re := regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	return re.MatchString(title) && title != ""
}

func internalAddNewPageHandler(w http.ResponseWriter, r *http.Request) {
	// if its not valid, redirect to the add new page form, with a name error message
	displayTitle := r.FormValue("displayTitle")
	if !isValidPageTitle(displayTitle) {
		// Prepare data to pass to the template, including the error message
		data := map[string]interface{}{
			"Error":          "Title must be alphanumeric and cannot contain special characters or underscores.",
			"AttemptedTitle": displayTitle,
		}
		// Render the form template with the data
		templateErr := templates.ExecuteTemplate(w, "add_new_page.html", data)
		if templateErr != nil {
			http.Error(w, templateErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	p := &Page{
		Title:         displayTitle,
		InternalTitle: titleToInternalTitle(displayTitle),
		Body:          []byte(""),
	}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/edit/"+p.InternalTitle, http.StatusFound)
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9_]+)$")

func extractInternalTitleFromPageRequest(r *http.Request) string {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return ""
	}
	return m[2]
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	internalTitle := extractInternalTitleFromPageRequest(r)
	if internalTitle == "" {
		http.NotFound(w, r)
		return
	}

	p, err := loadPage(internalTitle)
	if err != nil {
		http.Redirect(w, r, "/edit/"+internalTitle, http.StatusFound)
		return
	}
	templateErr := templates.ExecuteTemplate(w, "view.html", p)
	if templateErr != nil {
		http.Error(w, templateErr.Error(), http.StatusInternalServerError)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	internalTitle := extractInternalTitleFromPageRequest(r)
	if internalTitle == "" {
		http.NotFound(w, r)
		return
	}

	p, err := loadPage(internalTitle)
	if err != nil {
		p = &Page{
			Title:         internalTitleToTitle(internalTitle),
			InternalTitle: internalTitle,
		}
	}
	templateErr := templates.ExecuteTemplate(w, "edit.html", p)
	if templateErr != nil {
		http.Error(w, templateErr.Error(), http.StatusInternalServerError)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	internalTitle := extractInternalTitleFromPageRequest(r)
	if internalTitle == "" {
		http.NotFound(w, r)
		return
	}

	body := r.FormValue("body")

	p := &Page{
		Title:         internalTitleToTitle(internalTitle),
		InternalTitle: internalTitle,
		Body:          []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+internalTitle, http.StatusFound)
}
