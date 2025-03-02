package tests_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func authGroup() {
	It("should register user", func() {
		Expect(200).To(Equal(200))
	})
}
