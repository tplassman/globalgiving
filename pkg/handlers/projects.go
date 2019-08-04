package handlers

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"

	"../globalgiving"
)

type ViewData struct {
	Limit   int
	Limits  []int
	Page    int
	Results globalgiving.SearchResults
	Search  string
	Theme   string
	Themes  []globalgiving.Theme
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// Get search parameters from query
	q := r.URL.Query()
	limits := []int{10, 20, 50, 100, 1000}
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		limit = limits[0]
	}
	// Don't allow limits larger than our largest
	if limit > limits[len(limits)-1] {
		limit = limits[len(limits)-1]
	}
	themes := globalgiving.GetThemes()
	theme := q.Get("theme")
	search := q.Get("search")
	page, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		page = 1
	}
	start := (page - 1) * limit

	// Build project search from query params
	s := globalgiving.ProjectSearch{
		Limit:  limit,
		Page:   page,
		Search: search,
		Start:  start,
		Theme:  theme,
	}

	results, err := globalgiving.GetProjects(s)
	if err != nil {
		fmt.Printf("Error getting projects: %s", err)
	}

	// Define template functions
	funcMap := template.FuncMap{
		"excerpt": func(summary string) string {
			max := 139

			if len(summary) > max {
				summary = fmt.Sprintf("%s...", summary[:max])
			}

			return summary
		},
		"fundingProgress": func(funding, goal float64) float64 {
			if int(goal) == 0 {
				return float64(0)
			}

			progress := math.Floor(funding/goal*100*10) / 10

			return math.Min(progress, float64(100))
		},
		"getThemeName": func(id string) string {
			name := ""
			for _, t := range themes {
				if t.ID == id {
					name = t.Name
				}
			}

			return name
		},
		"hasPrevPage": func() bool { return page > 1 },
		"getPrevPage": func() int { return int(math.Min(float64(1), float64(page-1))) },
		"hasNextPage": func() bool { return (page-1)*limit+len(results.Projects.List) < results.NumberFound },
		"getNextPage": func() int { return page + 1 },
	}

	// Render template
	viewData := ViewData{
		Limit:   limit,
		Limits:  limits,
		Page:    page,
		Results: results,
		Search:  search,
		Theme:   theme,
		Themes:  themes,
	}
	t, err := template.New("main.tmpl").Funcs(funcMap).ParseFiles(
		"templates/_layouts/main.tmpl",
		"templates/projects.tmpl",
	)
	if err != nil {
		panic(err)
	}
	t.Execute(w, viewData)
}
