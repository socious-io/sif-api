package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func projectsGroup() {
	It("should create project", func() {
		for i, data := range projectsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			projectsData[i]["id"] = body["id"]
		}
	})

	It("should get list of projects", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/projects?page=1&limit=10", nil)
		req.Header.Set("Authorization", usersAuths[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(body["results"]).To(HaveLen(1))
	})
	It("should get single project", func() {
		for _, data := range projectsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%s", data["id"]), nil)
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(body["id"]).To(Equal(data["id"]))
		}
	})

	It("should vote to project", func() {
		for _, data := range projectsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", fmt.Sprintf("/projects/%s/votes", data["id"]), nil)
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusCreated))
		}
	})
}
