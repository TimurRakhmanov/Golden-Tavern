package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RakhmanovTimur/bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type postedData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	// {"make-res", "/make-reservation", "GET",  http.StatusOK},
	// {"post-search-availability", "/search-availability", "Post", []postedData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"post-search-availability-json", "/search-availability-json", "Post", []postedData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"make-reservation", "/make-reservation", "Post", []postedData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "me@here.com"},
	// 	{key: "phone", value: "555-555-5555"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_MakeReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Secret HQ",
		},
	}

	req, err := http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Log(err)
	}
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.MakeReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("start_date", "2030-01-01")
	postedData.Add("end_date", "2030-01-02")
	postedData.Add("first_name", "Tim")
	postedData.Add("last_name", "Timii")
	postedData.Add("email", "ewim@ddcs.com")
	postedData.Add("phone", "1231-2123-1211")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("postreservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("postreservation handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid start date
	postedData = url.Values{}
	postedData.Add("start_date", "invalid")
	postedData.Add("end_date", "2030-01-02")
	postedData.Add("first_name", "Tim")
	postedData.Add("last_name", "Timii")
	postedData.Add("email", "ewim@ddcs.com")
	postedData.Add("phone", "1231-2123-1211")
	postedData.Add("room_id", "1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("postreservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// Test for invalid end date
	postedData = url.Values{}
	postedData.Add("start_date", "2030-01-01")
	postedData.Add("end_date", "invalid")
	postedData.Add("first_name", "Tim")
	postedData.Add("last_name", "Timii")
	postedData.Add("email", "ewim@ddcs.com")
	postedData.Add("phone", "1231-2123-1211")
	postedData.Add("room_id", "1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("postreservation handler returned wrong response code for invalid start end: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid room id
	postedData = url.Values{}
	postedData.Add("start_date", "2030-01-01")
	postedData.Add("end_date", "2030-01-02")
	postedData.Add("first_name", "Tim")
	postedData.Add("last_name", "Timii")
	postedData.Add("email", "ewim@ddcs.com")
	postedData.Add("phone", "1231-2123-1211")
	postedData.Add("room_id", "invalid")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("postreservation handler returned wrong response code for invalid room id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid data
	postedData = url.Values{}
	postedData.Add("start_date", "2030-01-01")
	postedData.Add("end_date", "2030-01-02")
	postedData.Add("first_name", "T")
	postedData.Add("last_name", "Timii")
	postedData.Add("email", "ewim@ddcs.com")
	postedData.Add("phone", "1231-2123-1211")
	postedData.Add("room_id", "1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("postreservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert the reservation into database
	postedData = url.Values{}
	postedData.Add("start_date", "2030-01-01")
	postedData.Add("end_date", "2030-01-02")
	postedData.Add("first_name", "Tim")
	postedData.Add("last_name", "Timii")
	postedData.Add("email", "ewim@ddcs.com")
	postedData.Add("phone", "1231-2123-1211")
	postedData.Add("room_id", "2")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("postreservation handler failed when tried to validate inserting reseravtion response code for invalid data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Fail to insert restriction
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "Tim")
	postedData.Add("last_name", "Timii")
	postedData.Add("email", "ewim@ddcs.com")
	postedData.Add("phone", "1231-2123-1211")
	postedData.Add("room_id", "1000")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("postreservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// case 1: rooms are not available
	reqBody := "start=2070-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2070-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get response recorder
	rr := httptest.NewRecorder()

	// make handler handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// make request to our handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
	if j.OK {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}

	// case 2: no request body

	// create our request without body
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if j.OK || j.Message != "Internal Server Error" {
		t.Error("Got availability when request body was empty")
	}

	// case 3: error querying the database
	reqBody = "start=2100-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2100-01-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create our request without body
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if j.OK || j.Message != "Error Querying Database" {
		t.Error("Got availability when should have got the error querying database")
	}
}

func TestRepository_PostAvailability(t *testing.T) {

	// case 1: Correct request

	// Create the body
	reqBody := "start=2010-10-10"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2010-10-11")

	// Create a new request
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))

	// Add a context
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// Set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Simulate a handler
	handler := http.HandlerFunc(Repo.PostAvailability)

	// make the request to handler
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("Expected Status OK")
	}

	// case 2: no request body

	// Create a new request
	req, _ = http.NewRequest("POST", "/search-availability", nil)

	// Add a context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// Set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder
	rr = httptest.NewRecorder()

	// Simulate a handler
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to handler
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Error("Did not get status Temporary Redirect without body")
	}

	// case 3: no availability

	reqBody = "start=2100-10-10"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2100-10-11")

	// Create a new request
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))

	// Add a context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// Set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder
	rr = httptest.NewRecorder()

	// Simulate a handler
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to handler
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Error("Did not get status Temporary Redirect without body")
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	// Case 1: Correct request
	// Make a reservation
	res := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Secret HQ",
		},
	}

	// Make a request with context
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", res)
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// case 2: no session

	// Make a request with context
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	// Case 1: Correct request
	// Make a reservation
	res := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Secret HQ",
		},
	}
	// Make a request with context
	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/choose-room/1"
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", res)
	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChoosRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

}

func TestRepository_BookRoom(t *testing.T) {
	// Case 1: correct request
	// Make a reservation
	res := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Secret HQ",
		},
	}
	// Make a request with context
	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", res)
	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)

	}

	// Case 2: database failed

	// Make a request with context
	req, _ = http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=10", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)

	}

}
