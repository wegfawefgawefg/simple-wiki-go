package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index", http.StatusFound)
}

type IndexPageListing struct {
	Title         string
	InternalTitle string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pageIndexHandler start")
	files, err := os.ReadDir(saved_pages_path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var page_listings []IndexPageListing
	for _, file := range files {
		page_listing := IndexPageListing{
			InternalTitle: file.Name(),
			Title:         internalTitleToTitle(file.Name()),
		}
		page_listings = append(page_listings, page_listing)
	}

	template_err := templates.ExecuteTemplate(w,
		"index.html", page_listings)
	if template_err != nil {
		http.Error(w, template_err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("pageIndexHandler end")
}

// as-is renders the add new page template
func addNewPageHandler(w http.ResponseWriter, r *http.Request) {
	template_err := templates.ExecuteTemplate(w, "add_new_page.html", nil)
	if template_err != nil {
		http.Error(w, template_err.Error(), http.StatusInternalServerError)
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
		templates.ExecuteTemplate(w, "add_new_page.html", data)
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
	fmt.Println("viewHandler start")
	internalTitle := extractInternalTitleFromPageRequest(r)
	p, err := loadPage(internalTitle)
	if err != nil {
		http.Redirect(w, r, "/edit/"+p.InternalTitle, http.StatusFound)
		return
	}
	template_err := templates.ExecuteTemplate(w, "view.html", p)
	if template_err != nil {
		http.Error(w, template_err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("viewHandler end")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("editHandler start")
	title := extractInternalTitleFromPageRequest(r)
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	template_err := templates.ExecuteTemplate(w, "edit.html", p)
	if template_err != nil {
		http.Error(w, template_err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("editHandler end")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("saveHandler start")
	internalTitle := extractInternalTitleFromPageRequest(r)
	fmt.Println("internalTitle: ", internalTitle)
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
	fmt.Println("saveHandler end")
}
