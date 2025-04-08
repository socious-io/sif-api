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

	It("should comment on project", func() {
		for _, data := range projectsData {
			for i, cm := range commentsData {
				w := httptest.NewRecorder()
				reqBody, _ := json.Marshal(cm)
				req, _ := http.NewRequest("POST", fmt.Sprintf("/projects/%s/comments", data["id"]), bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", usersAuths[0])
				router.ServeHTTP(w, req)
				body := decodeBody(w.Body)
				Expect(w.Code).To(Equal(http.StatusCreated))
				Expect(body["content"]).To(Equal(cm["content"]))
				commentsData[i]["id"] = body["id"]
			}
		}
	})

	It("should update comment", func() {

		for i, cm := range commentsData {
			w := httptest.NewRecorder()
			content := "updated comment"
			reqBody, _ := json.Marshal(map[string]interface{}{
				"content": content,
			})
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/projects/comments/%s", cm["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(body["content"]).To(Equal(content))
			commentsData[i]["content"] = content
		}

	})

	It("should like comment", func() {
		for _, cm := range commentsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", fmt.Sprintf("/projects/comments/%s/likes", cm["id"]), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusCreated))
		}
	})

	It("should react to comment", func() {
		for _, cm := range commentsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(map[string]interface{}{
				"reaction": "👍",
			})
			req, _ := http.NewRequest("POST", fmt.Sprintf("/projects/comments/%s/reactions", cm["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusCreated))
		}
	})

	It("should get project comments", func() {
		for _, data := range projectsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%s/comments", data["id"]), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(len(body)).To(BeNumerically(">", 0))
		}
	})

	It("should get single comment", func() {
		for _, cm := range commentsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/comments/%s", cm["id"]), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(body["content"]).To(Equal(cm["content"]))
		}
	})

	It("should unlike comment", func() {
		for _, cm := range commentsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/projects/comments/%s/likes", cm["id"]), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
		}
	})

	It("should delete reaction", func() {
		for _, cm := range commentsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/projects/comments/%s/reactions", cm["id"]), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
		}
	})

	It("should delete comment", func() {
		for _, cm := range commentsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/projects/comments/%s", cm["id"]), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", usersAuths[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
		}
	})

}
