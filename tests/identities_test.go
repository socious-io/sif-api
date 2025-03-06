package tests_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func identitiesGroup() {

	It("should get identities", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/identities", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", usersAuths[0])
		router.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
		body := decodeBody(w.Body)
		identities := body["identities"].([]interface{})
		Expect(len(identities)).NotTo(Equal(0))
	})
}
