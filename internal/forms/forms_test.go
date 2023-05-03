package forms

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Got valid when should have been invalid")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Got invalid when should have been valid")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")
	form := New(postedData)
	if !form.Has("a") {
		t.Error("Got false when should have been true")
	}
	if form.Has("e") {
		t.Error("Got true when should have been false")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	field := "a"
	postedData.Add(field, "aaa")
	form := New(postedData)
	form.MinLength(field, 3)
	if !form.Valid() {
		t.Error(fmt.Sprintf("Expected length of 3, but got %d instead", len(form.Get(field))))
	}
	isError := form.Errors.Get(field)
	if isError != "" {
		t.Error("Should not have an error, but got one")
	}

	field = "b"
	postedData.Add(field, "bb")
	form = New(postedData)
	form.MinLength(field, 3)
	if form.Valid() {
		t.Error("Expected false, but got true")
	}
	isError = form.Errors.Get(field)
	if isError == "" {
		t.Error("Should have an error, but did not get one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("email", "email@gmai.com")
	form := New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Expected valid email, but got invalid instead")
	}
	postedData = url.Values{}
	postedData.Add("email", "32efds@")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("Expected invalid email, but got valid instead")
	}
}
