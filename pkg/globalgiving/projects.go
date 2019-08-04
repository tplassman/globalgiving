package globalgiving

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type project struct {
	Funding     float64
	Goal        float64
	ID          int
	ImageLink   string
	ProjectLink string
	Status      string
	Summary     string
	Title       string
}

type SearchResults struct {
	Projects struct {
		List []project `json:"project"`
	}
	NumberFound int `json:"numberFound"`
}

type ProjectSearch struct {
	Limit  int
	Page   int
	Theme  string
	Search string
	Start  int
}

func getResults(c *client, r *http.Request) (SearchResults, error) {
	// Get API response
	res, err := c.httpClient.Do(r)
	if err != nil {
		return SearchResults{}, err
	}
	defer res.Body.Close()
	// Decode JSON from response body in search results
	sr := struct {
		Search struct {
			SearchResults SearchResults `json:"response"`
		}
	}{}
	dec := json.NewDecoder(res.Body)
	for {
		if err := dec.Decode(&sr); err == io.EOF {
			break
		} else if err != nil {
			return SearchResults{}, err
		}
	}

	return sr.Search.SearchResults, nil
}

func getRequest(c *client, v url.Values) (*http.Request, error) {
	// Build out projects request from client and values
	u := c.baseURL.ResolveReference(&url.URL{
		Path:     "/api/public/services/search/projects",
		RawQuery: v.Encode(),
	})
	// Make request from URL string
	r, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	// Get JSON
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")

	return r, nil
}

func getValues(c *client, s ProjectSearch) url.Values {
	// Build query values from client and search
	v := url.Values{}
	// Add API key from client
	v.Add("api_key", c.apiKey)
	// Add search
	if s.Search == "" {
		s.Search = "*"
	}
	v.Add("q", s.Search)
	// Add start
	v.Add("start", strconv.Itoa(s.Start))
	// Add filters (e.g. theme, organization, etc.)
	filters := []string{}
	if s.Theme != "" {
		filters = append(filters, fmt.Sprintf("theme:%s", s.Theme))
	}
	if len(filters) > 0 {
		v.Add("filter", strings.Join(filters[:], ","))
	}

	return v
}

type response struct {
	index   int
	results SearchResults
}

func GetProjects(s ProjectSearch) (SearchResults, error) {
	// Global Giving API limits results set to maximum of 10
	limit := 10
	// If search limit from front end is less than 10 just make one request
	count := int(math.Max(float64(s.Limit)/float64(limit), float64(1)))
	// Create client to use for requests
	c := newClient()
	// Create channels to pass API responses and errors
	ch := make(chan response, count)
	errch := make(chan error)

	fmt.Printf("Making %d requests real quick\n", count)
	for i := 0; i < count; i++ {
		go func(c *client, s ProjectSearch, i int) {
			v := getValues(c, s)
			r, err := getRequest(c, v)
			if err != nil {
				errch <- err
				return
			}
			sr, err := getResults(c, r)
			if err != nil {
				errch <- err
				return
			}
			ch <- response{i, sr}
		}(c, s, i)

		// Increase starting position for next request
		s.Start = s.Start + limit
	}

	// Need to store total number of results returned to break if we have enough
	total := 0
	// Create map to store results in w/ index as key
	responseMap := make(map[int]SearchResults)

	// Wait for return from all porject requests
	for i := 0; i < count; i++ {
		select {
		case err := <-errch:
			fmt.Printf("Error fetching projects: %s", err)
		case res := <-ch:
			responseMap[res.index] = res.results

			// Check if we got what we were looking for
			total += len(res.results.Projects.List)
			if total == res.results.NumberFound {
				break
			}
		}
	}

	// Close open channels
	close(ch)
	close(errch)

	// Create empty results set to populate and return
	results := SearchResults{}
	// Order results using map keys
	for i := 0; i < count; i++ {
		sr := responseMap[i]
		results.Projects.List = append(results.Projects.List, sr.Projects.List...)
		results.NumberFound = sr.NumberFound
	}

	return results, nil
}
