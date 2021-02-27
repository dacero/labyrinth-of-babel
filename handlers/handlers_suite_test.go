package handlers_test

import (
	"testing"
	"net/http"	
	"net/http/httptest"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	
	"github.com/dacero/labyrinth-of-babel/repository"
	"github.com/dacero/labyrinth-of-babel/handlers"
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = Describe("Handler", func() {
	var (
		lobRepo repository.LobRepository
		req		*http.Request
		rr		*httptest.ResponseRecorder
		handler	func(w http.ResponseWriter, r *http.Request)
		err     error
		cellId  = "72aed05b-cb2d-4cad-bf70-05d8ae02a7bc"
	)

	BeforeEach(func() {
		lobRepo = repository.NewLobRepository()
		rr = httptest.NewRecorder()
	})

	Describe("Viewing a card", func() {
		Context("for a cell that exists", func() {
			BeforeEach(func() {
				req, err = http.NewRequest("GET", "http://localhost:8080/view/"+cellId, nil)
				Expect(err).To(BeNil())
				handler = handlers.ViewHandler(lobRepo)
				handler(rr, req)
			})
			It("should return Status OK", func() {
				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
		Context("for a cell that does not exist", func() {
			BeforeEach(func() {
				req, err = http.NewRequest("GET", "http://localhost:8080/view/error", nil)
				Expect(err).To(BeNil())
				handler = handlers.ViewHandler(lobRepo)
				handler(rr, req)
			})
			It("should return NOT FOUND error", func() {
				Expect(rr.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})