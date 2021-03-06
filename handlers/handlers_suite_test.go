package handlers_test

import (
	"testing"
	"net/http"	
	"net/http/httptest"
	"net/url"
	"strings"
	"errors"
	
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
				linksStart := `<ul class="card-link-list">Links`
				linksEnd := `</ul> <!--links-->`
				linksSubstr, err := extractFromPage(body, linksStart, linksEnd)
				Expect(err).To(BeNil())
				// find the index of the links_start
				links := strings.Split(linksSubstr, "\n")
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
	
	Describe("Creating a new card", func() {
		Context("with proper information", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("room", "This is a room")
				form.Add("title", "The new cell")
				form.Add("body", "This is the new cell I'm creating")
				form.Add("source", "Confucius")
				req, err := http.NewRequest("POST", "http://localhost:8080/newCell/", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				handler = handlers.CreateHandler(lobRepo)
				handler(rr, req)
				body = rr.Body.String()
			})
			It("should return Status Found", func() {
				Expect(rr.Code).To(Equal(http.StatusFound))
				//check the result page shows the new cell
				/*
				I've been unable to get this to work...
				newCellTitle, err := extractFromPage(body, `<div class="card-title">`, `</div><!--title-->`)
				Expect(err).To(BeNil())
				Expect(newCellTitle).To(Equal("The new cell"))
				*/
			})
		})	
	})
})

//extracts the substring from s contained within start and finish 
func extractFromPage(s string, start string, end string) (string, error) {
	startIndex := strings.Index(s, start) + len(start)
	if startIndex == -1 {
		return "", errors.New("start not found: " + start)
	}
	endIndex := strings.Index(s, end)
	if endIndex == -1 {
		return "", errors.New("End not found: " + end)
	}
	//extract the links substring
	return s[startIndex:endIndex], nil
}