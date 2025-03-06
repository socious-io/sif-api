package tests_test

import (
	"bytes"
	"encoding/json"
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
}
