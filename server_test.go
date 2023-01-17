package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMetadataPost(t *testing.T) {
	// Test valid payload
	validPayload := []byte(`title: Valid App 1
version: 0.0.1
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app1
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`)

	req, err := http.NewRequest("POST", "/metadata", bytes.NewBuffer(validPayload))
	req.Header.Set("Content-Type", "application/x-yaml")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleMetadataPost)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	// Test invalid payload with invalid email
	invalidPayload := []byte(`title: App w/ Invalid maintainer email
version: 1.0.1
maintainers:
- name: Firstname Lastname
  email: apptwohotmail.com
company: Upbound Inc.
website: https://upbound.io
source: https://github.com/upbound/repo
license: Apache-2.0
description: |
 ### blob of markdown
 More markdown
`)

	req, err = http.NewRequest("POST", "/metadata", bytes.NewBuffer(invalidPayload))
	req.Header.Set("Content-Type", "application/x-yaml")
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid Email Address")

	// Test invalid payload with missing version
	invalidPayload = []byte(`title: App w/ missing version
maintainers:
- name: first last
  email: email@hotmail.com
- name: first last
  email: email@gmail.com
company: Company Inc.
website: https://website.com
source: https://github.com/company/repo
license: Apache-2.0
description: |
 ### blob of markdown
 More markdown
`)

	req, err = http.NewRequest("POST", "/metadata", bytes.NewBuffer(invalidPayload))

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "All fields are required")

}

func TestHandleMetadataSearch(t *testing.T) {
	appMetadata = nil
	appMetadata = append(appMetadata, AppMetadata{
		Title:       "Test App 1",
		Version:     "0.0.1",
		Maintainers: []Person{{Name: "John Doe", Email: "johndoe@example.com"}},
		Company:     "Test Company Inc.",
		Website:     "https://testcompany.com",
		Source:      "https://github.com/testcompany/testapp",
		License:     "MIT",
		Description: "This is a test app.",
	})
	appMetadata = append(appMetadata, AppMetadata{
		Title:       "Test App 2",
		Version:     "0.0.2",
		Maintainers: []Person{{Name: "Jane Doe", Email: "janedoe@example.com"}},
		Company:     "Test Company Inc.",
		Website:     "https://testcompany.com",
		Source:      "https://github.com/testcompany/testapp2",
		License:     "MIT",
		Description: "This is another test app.",
	})

	// Test searching by title
	req, err := http.NewRequest("GET", "/metadata/search?title=Test App 1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleMetadataSearch)
	handler.ServeHTTP(rr, req)
	expected := []AppMetadata{{
		Title:       "Test App 1",
		Version:     "0.0.1",
		Maintainers: []Person{{Name: "John Doe", Email: "johndoe@example.com"}},
		Company:     "Test Company Inc.",
		Website:     "https://testcompany.com",
		Source:      "https://github.com/testcompany/testapp",
		License:     "MIT",
		Description: "This is a test app.",
	}}

	jsonData, err := json.Marshal(expected)
	assert.Equal(t, rr.Result().StatusCode, http.StatusOK)
	assert.Equal(t, rr.Body.String(), string(jsonData))

	// Test searching by company
	req, err = http.NewRequest("GET", "/metadata/search?company=Test Company Inc.", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	expected = append(expected, AppMetadata{
		Title:       "Test App 2",
		Version:     "0.0.2",
		Maintainers: []Person{{Name: "Jane Doe", Email: "janedoe@example.com"}},
		Company:     "Test Company Inc.",
		Website:     "https://testcompany.com",
		Source:      "https://github.com/testcompany/testapp2",
		License:     "MIT",
		Description: "This is another test app.",
	})
	jsonData, err = json.Marshal(expected)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, string(jsonData), rr.Body.String())
}
