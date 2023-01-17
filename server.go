package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type AppMetadata struct {
	Title       string   `yaml:"title"`
	Version     string   `yaml:"version"`
	Maintainers []Person `yaml:"maintainers"`
	Company     string   `yaml:"company"`
	Website     string   `yaml:"website"`
	Source      string   `yaml:"source"`
	License     string   `yaml:"license"`
	Description string   `yaml:"description"`
}

var appMetadata []AppMetadata

type Person struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/metadata", handleMetadataPost).Methods("POST")
	router.HandleFunc("/metadata/search", handleMetadataSearch).Methods("GET")

	http.ListenAndServe(":8080", router)
}

func handleMetadataPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var data AppMetadata
	if err := yaml.Unmarshal(body, &data); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	if data.Title == "" || data.Version == "" || len(data.Maintainers) == 0 || data.Company == "" || data.Website == "" || data.Source == "" || data.License == "" || data.Description == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	for _, maintainer := range data.Maintainers {
		if maintainer.Name == "" || maintainer.Email == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if !valideEmail(maintainer.Email) {
			http.Error(w, "Invalid Email Address", http.StatusBadRequest)
			return
		}
	}

	// Persist data
	appMetadata = append(appMetadata, data)

	w.WriteHeader(http.StatusCreated)
}

// Supports only one value as query parameter It performans and of all the mentioned fields
// Maintainer is matched with both email id and name just to make search only one level
func handleMetadataSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if len(query) == 0 {
		http.Error(w, "search query is empty", http.StatusBadRequest)
		return
	}

	var results []AppMetadata

	for _, m := range appMetadata {
		match := true
		for key, value := range query {
			switch key {
			case "title":
				if m.Title != value[0] {
					match = false
				}
				break
			case "version":
				if m.Version != value[0] {
					match = false
				}
				break
			case "maintainer":
				// matching both name and email id
				matchMaintainer := false
				for _, p := range m.Maintainers {
					if strings.Contains(p.Name, value[0]) {
						matchMaintainer = true
						break
					}
					if strings.Contains(p.Email, value[0]) {
						matchMaintainer = true
						break
					}
				}
				if !matchMaintainer {
					match = false
				}
				break
			case "company":
				if m.Company != value[0] {
					match = false
				}
				break
			case "website":
				if m.Website != value[0] {
					match = false
				}
				break
			case "source":
				if m.Source != value[0] {
					match = false
				}
				break
			case "license":
				if m.License != value[0] {
					match = false
				}
				break
			case "description":
				if !strings.Contains(m.Description, value[0]) {
					match = false
				}
				break
			default:
				match = false
				break
			}
		}
		if match == true {
			results = append(results, m)
		}
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func valideEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
