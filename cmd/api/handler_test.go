package main

import (
	"context"
	"encoding/json"
	"github.com/3n0ugh/simple-crud-api/internal/data"
	"github.com/3n0ugh/simple-crud-api/internal/data/mock"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type in struct {
	urlPath string
	param   string
	reqBody data.Book
}

type out struct {
	envelopeName string
	statusCode   int
	body         any
}

type testCases map[string]struct {
	in       in
	expected out
}

type testCasesBook struct {
	Author string `json:"author,omitempty"`
	Name   string `json:"name,omitempty"`
}

func TestHandleAddBook(t *testing.T) {
	app := &application{model: mock.NewModel()}

	tc := testCases{
		"Must_Success": {
			in: in{
				urlPath: "localhost:8080/v1/book",
				reqBody: data.Book{
					Name:   "TestName",
					Author: "TestAuthor",
				},
			},
			expected: out{
				envelopeName: "book",
				statusCode:   http.StatusCreated,
				body: data.Book{
					ID:     mock.Book.ID,
					Name:   "TestName",
					Author: "TestAuthor",
				},
			},
		},
		"Empty_Request": {
			in: in{
				urlPath: "localhost:8080/v1/book",
				reqBody: data.Book{
					Author: "",
					Name:   "",
				},
			},
			expected: out{
				envelopeName: "error",
				statusCode:   http.StatusUnprocessableEntity,
				body: testCasesBook{
					Name:   "must be provided",
					Author: "must be provided",
				},
			},
		},
		"Long_Name": {
			in: in{
				urlPath: "localhost:8080/v1/book",
				reqBody: data.Book{
					Name:   string(make([]byte, 41)),
					Author: "TestAuthor",
				},
			},
			expected: out{
				envelopeName: "error",
				statusCode:   http.StatusUnprocessableEntity,
				body: testCasesBook{
					Name: "must not be more than 40 bytes long",
				},
			},
		},
		"Long_Author": {
			in: in{
				urlPath: "localhost:8080/v1/book",
				reqBody: data.Book{
					Name:   "TestName",
					Author: string(make([]byte, 41)),
				},
			},
			expected: out{
				envelopeName: "error",
				statusCode:   http.StatusUnprocessableEntity,
				body: testCasesBook{
					Author: "must not be more than 40 bytes long",
				},
			},
		},
	}

	for scenario, tt := range tc {
		t.Run(scenario, func(t *testing.T) {
			reqBodyJSON, err := json.Marshal(tt.in.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			tt.expected.body, err = app.prettyJSON(envelope{tt.expected.envelopeName: tt.expected.body})
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, tt.in.urlPath, strings.NewReader(string(reqBodyJSON)))

			app.handleAddBook(w, r)

			body, err := ioutil.ReadAll(w.Result().Body)
			if err != nil {
				t.Fatal(err)
			}

			if tt.expected.statusCode != w.Result().StatusCode {
				t.Errorf("Status Code -> want: %d; got: %d", tt.expected.statusCode, w.Result().StatusCode)
			}

			if !reflect.DeepEqual(tt.expected.body, body) {
				t.Errorf("body -> want: \n%q; got: \n%q", tt.expected.body, body)
			}
		})
	}
}

func TestHandleGetBook(t *testing.T) {
	app := &application{model: mock.NewModel()}

	tc := testCases{
		"Must_Success": {
			in: in{
				urlPath: "localhost:8080/v1/book",
			},
			expected: out{
				envelopeName: "books",
				statusCode:   http.StatusOK,
				body:         mock.Books,
			},
		},
	}

	for scenario, tt := range tc {
		t.Run(scenario, func(t *testing.T) {
			reqBodyJSON, err := json.Marshal(tt.in.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			tt.expected.body, err = app.prettyJSON(envelope{tt.expected.envelopeName: tt.expected.body})
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.in.urlPath, strings.NewReader(string(reqBodyJSON)))

			app.handleGetBook(w, r)

			body, err := ioutil.ReadAll(w.Result().Body)
			if err != nil {
				t.Fatal(err)
			}

			if tt.expected.statusCode != w.Result().StatusCode {
				t.Errorf("Status Code -> want: %d; got: %d", tt.expected.statusCode, w.Result().StatusCode)
			}

			if !reflect.DeepEqual(tt.expected.body, body) {
				t.Errorf("body -> want: \n%q; got: \n%q", tt.expected.body, body)
			}
		})
	}
}

func TestHandleDeleteBook(t *testing.T) {
	app := &application{model: mock.NewModel()}

	tc := testCases{
		"Valid_ID": {
			in: in{
				urlPath: "localhost:8080/v1/book/1",
				param:   "1",
			},
			expected: out{
				envelopeName: "message",
				statusCode:   http.StatusOK,
				body:         "book successfully deleted",
			},
		},
		"Negative_ID": {
			in: in{
				urlPath: "localhost:8080/v1/book/-1",
				param:   "-1",
			},
			expected: out{
				envelopeName: "error",
				statusCode:   http.StatusBadRequest,
				body:         "invalid id parameter",
			},
		},
		"Non-Existent_ID": {
			in: in{
				urlPath: "localhost:8080/v1/book/2000",
				param:   "2000",
			},
			expected: out{
				envelopeName: "error",
				statusCode:   http.StatusNotFound,
				body:         "the requested resource could not be found",
			},
		},
		"Decimal_ID": {
			in: in{
				urlPath: "localhost:8080/v1/book/1.2",
				param:   "1.2",
			},
			expected: out{
				envelopeName: "error",
				statusCode:   http.StatusBadRequest,
				body:         "invalid id parameter",
			},
		},
		"Empty_ID": {
			in: in{
				urlPath: "localhost:8080/v1/book/",
			},
			expected: out{
				envelopeName: "error",
				statusCode:   http.StatusBadRequest,
				body:         "invalid id parameter",
			},
		},
	}

	for scenario, tt := range tc {
		t.Run(scenario, func(t *testing.T) {
			reqBodyJSON, err := json.Marshal(tt.in.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			tt.expected.body, err = app.prettyJSON(envelope{tt.expected.envelopeName: tt.expected.body})
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, tt.in.urlPath, strings.NewReader(string(reqBodyJSON)))
			ctx := r.Context()
			ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
				{"id", tt.in.param},
			})
			r = r.WithContext(ctx)

			app.handleDeleteBook(w, r)

			body, err := ioutil.ReadAll(w.Result().Body)
			if err != nil {
				t.Fatal(err)
			}

			if tt.expected.statusCode != w.Result().StatusCode {
				t.Errorf("Status Code -> want: %d; got: %d", tt.expected.statusCode, w.Result().StatusCode)
			}

			if !reflect.DeepEqual(tt.expected.body, body) {
				t.Errorf("body -> want: \n%q; got: \n%q", tt.expected.body, body)
			}
		})
	}

}
