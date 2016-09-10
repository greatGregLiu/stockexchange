package integration_test

import (
	"net/http"
	"os"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	var (
		args       []string
		session    *gexec.Session
		sessionErr error
	)

	JustBeforeEach(func() {
		session, sessionErr = runner.Start(args...)
	})

	AfterEach(func() {
		if session != nil {
			session.Kill()
		}
	})

	It("is listenting on HTTP port 9292", func() {
		Expect(sessionErr).NotTo(HaveOccurred())
		_, err := http.Get("http://127.0.0.1:9292")
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Environment variable PORT", func() {
		BeforeEach(func() {
			os.Setenv("PORT", "8899")
		})

		AfterEach(func() {
			os.Clearenv()
		})

		It("is listenting on that HTTP port", func() {
			Expect(sessionErr).NotTo(HaveOccurred())
			resp, err := http.Get("http://127.0.0.1:8899")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("when the port is not integer", func() {
			BeforeEach(func() {
				os.Setenv("PORT", "wrong_port")
			})

			It("is listenting on the default HTTP port", func() {
				Expect(sessionErr).NotTo(HaveOccurred())
				_, err := http.Get("http://127.0.0.1:9292")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when command line argument is provided", func() {
			BeforeEach(func() {
				args = []string{"--addr=127.0.0.1:8080"}
			})

			It("does not overried the PORT", func() {
				Expect(sessionErr).NotTo(HaveOccurred())
				_, err := http.Get("http://127.0.0.1:8080")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("when the addr is provided", func() {
		BeforeEach(func() {
			args = []string{"--addr=127.0.0.1:8080"}
		})

		It("is listenting on that HTTP port", func() {
			Expect(sessionErr).NotTo(HaveOccurred())
			_, err := http.Get("http://127.0.0.1:8080")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when the addr is wrong", func() {
		BeforeEach(func() {
			args = []string{"--addr=wrong_host_and_port"}
		})

		It("returns an error", func() {
			Expect(session).To(BeNil())
			Expect(sessionErr.Error()).To(ContainSubstring("The provided wrong_host_and_port addr is not correct"))
		})
	})
})