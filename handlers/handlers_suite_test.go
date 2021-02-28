package handlers_test

import (
	"testing"
	"net/http"	
	"net/http/httptest"
	"strings"
	
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
		body	string
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
				body = rr.Body.String()
			})
			It("should return Status OK", func() {
				Expect(rr.Code).To(Equal(http.StatusOK))
			})
			It("should contain the right number of links", func() {
				// first I need to parse the body
				links_start := `<ul class="card-link-list">Links`
				links_end := `</ul> <!--links-->`
				// find the index of the links_start
				links_start_index := strings.Index(body, links_start) + len(links_start)
				Expect(links_start_index).ToNot(Equal(-1))
				links_end_index := strings.Index(body, links_end)
				Expect(links_end_index).ToNot(Equal(-1))
				//extract the links substring
				links_substr := body[links_start_index:links_end_index]
				links := strings.Split(links_substr, "\n")
				var clean_links []string
				for _, link := range links {
					if strings.Contains(link, `<li class="card-link">`) {
						clean_links = append(clean_links, strings.Trim(link, "\t "))
					}
				}
				//there should be 2 links in that cell
				Expect(len(clean_links)).To(Equal(2))
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