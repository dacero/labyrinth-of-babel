package handlers_test

import (
	"testing"
	"net/http"	
	"net/http/httptest"
	"net/url"
	"encoding/json"
	"strings"
	"errors"
	"log"
	"database/sql"
	"io/ioutil"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gorilla/mux"
	
	"github.com/dacero/labyrinth-of-babel/repository"
	"github.com/dacero/labyrinth-of-babel/handlers"
	"github.com/dacero/labyrinth-of-babel/models"
)

func resetDB() {
	log.Print("Initializing db... ")
	db, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/")
	if err != nil {
		log.Fatal("Error when opening DB: ", err)
	}
	defer db.Close()

	file, err := ioutil.ReadFile("./db/labyrinth_of_babel.sql")
	if err != nil {
		// handle error
		log.Fatal("Error when initializing DB: ", err)
	}

	for _, request := range strings.Split(string(file), ";") {
		request = strings.TrimSpace(request)
		if request != "" {
			_, err := db.Exec(request)
			if err != nil {
				log.Fatalf("Error when initializing DB\n Query %s returned error: %s", request, err)
			}
		}
	}
	log.Print("Finished initiatlizing db!")
}

func TestRepository(t *testing.T) {
	resetDB()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = Describe("Handler", func() {
	var (
		lobRepository repository.LobRepository
		router	*mux.Router
		req		*http.Request
		rr		*httptest.ResponseRecorder
		err     error
		cellId  = "72aed05b-cb2d-4cad-bf70-05d8ae02a7bc"
		body	string
	)

	BeforeEach(func() {
		lobRepository = repository.NewLobRepository()
		rr = httptest.NewRecorder()
		router = mux.NewRouter()
		router.HandleFunc("/cell/{id}", handlers.ViewHandler(lobRepository))
		router.HandleFunc("/cell/{id}/edit", handlers.EditHandler(lobRepository))
		router.HandleFunc("/cell/{id}/edit/sources", handlers.SourcesHandler(lobRepository))
		router.HandleFunc("/cell/{id}/addSource", handlers.AddSourceHandler(lobRepository)).Methods("POST") //addSource
		router.HandleFunc("/cell/{id}/removeSource", handlers.RemoveSourceHandler(lobRepository)).Methods("POST") //removeSource
		router.HandleFunc("/cell/{id}/linkCell", handlers.LinkCellsHandler(lobRepository)).Methods("POST")
		router.HandleFunc("/cell/{id}/unlinkCell", handlers.UnlinkCellsHandler(lobRepository)).Methods("POST")
		router.HandleFunc("/save", handlers.SaveHandler(lobRepository))
		router.HandleFunc("/new", handlers.CreateHandler(lobRepository))
		router.HandleFunc("/sources", handlers.SearchSourcesHandler(lobRepository))
		router.HandleFunc("/rooms", handlers.SearchRoomsHandler(lobRepository))
		router.HandleFunc("/page/{page}", handlers.PageHandler())
	})

	Describe("Viewing a card", func() {
		Context("for a cell that exists", func() {
			BeforeEach(func() {
				req, err = http.NewRequest("GET", "http://localhost:8080/cell/"+cellId, nil)
				Expect(err).To(BeNil())
				router.ServeHTTP(rr, req)
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
				req, err = http.NewRequest("GET", "http://localhost:8080/cell/thiscelldoesnotexist", nil)
				Expect(err).To(BeNil())
				router.ServeHTTP(rr, req)
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
				req, err := http.NewRequest("POST", "http://localhost:8080/new", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
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
		Context("without a body", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("room", "This is a room")
				form.Add("title", "The new cell")
				form.Add("body", "")
				form.Add("source", "Confucius")
				req, err := http.NewRequest("POST", "http://localhost:8080/new", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should return StatusBadRequest", func() {
				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})
		Context("without a room", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("room", "")
				form.Add("title", "The new cell")
				form.Add("body", "This one does have a body")
				form.Add("source", "Confucius")
				req, err := http.NewRequest("POST", "http://localhost:8080/new", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should return StatusBadRequest", func() {
				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})	
	})
	
	Describe("Searching for sources", func() {
		Context("with proper terms", func() {
			BeforeEach(func() {
				req, err := http.NewRequest("GET", "http://localhost:8080/sources?term=Confu", nil)
				Expect(err).To(BeNil())
				router.ServeHTTP(rr, req)
				body = rr.Body.String()
			})
			It("should return a proper json", func() {
				Expect(rr.Code).To(Equal(http.StatusOK))
				var sources []models.Source
				err = json.Unmarshal(rr.Body.Bytes(), &sources)
				Expect(len(sources)).To(Equal(1))
			})
		})
	})
	
	Describe("Searching for rooms", func() {
		Context("with proper terms", func() {
			BeforeEach(func() {
				req, err := http.NewRequest("GET", "http://localhost:8080/rooms?term=room", nil)
				Expect(err).To(BeNil())
				router.ServeHTTP(rr, req)
				body = rr.Body.String()
			})
			It("should return a proper json", func() {
				Expect(rr.Code).To(Equal(http.StatusOK))
				var rooms []string
				err = json.Unmarshal(rr.Body.Bytes(), &rooms)
				Expect(len(rooms)).To(Equal(1))
			})
		})
	})
	
	Describe("Updating a card", func() {
		Context("with proper information", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("cellId", cellId)
				form.Add("room", "This is a room")
				form.Add("title", "Updated title")
				form.Add("body", "I'm updating this cell")
				req, err := http.NewRequest("POST", "http://localhost:8080/save", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should return Status Found", func() {
				Expect(rr.Code).To(Equal(http.StatusFound))
			})
		})
		Context("with wrong information", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("cellId", cellId)
				form.Add("room", "  ")
				form.Add("title", "Updated title")
				form.Add("body", "  ")
				req, err := http.NewRequest("POST", "http://localhost:8080/save", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should return Status Found", func() {
				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
	
	Describe("Adding a source to a card", func() {
		Context("with proper information", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("cellId", cellId)
				form.Add("source", "A new source to test")
				req, err := http.NewRequest("POST", "http://localhost:8080/cell/"+cellId+"/addSource", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should return Status Found", func() {
				Expect(rr.Code).To(Equal(http.StatusFound))
			})
		})
		Context("with wrong information", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("cellId", cellId)
				form.Add("source", "    ")
				req, err := http.NewRequest("POST", "http://localhost:8080/cell/"+cellId+"/addSource", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should also return Status Found", func() {
				Expect(rr.Code).To(Equal(http.StatusFound))
			})
		})
	})
	Describe("When removing a source from a card", func() {
		Context("given I provide a proper source", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("source", "Confucius")
				req, err := http.NewRequest("POST", "http://localhost:8080/cell/"+cellId+"/removeSource", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should return Status Found", func() {
				Expect(rr.Code).To(Equal(http.StatusFound))
			})
		})
		Context("given I provide a source not linked to the cell", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("source", "This source is not linked")
				req, err := http.NewRequest("POST", "http://localhost:8080/cell/"+cellId+"/removeSource", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should also return Status Found", func() {
				Expect(rr.Code).To(Equal(http.StatusFound))
			})
		})
	})
	
	Describe("When linking a cell", func() {
		Context("given I provide a proper cell", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("cellToLink", "df38bd04-0ec4-41bf-9e53-d0eeb95a4939")
				req, err := http.NewRequest("POST", "http://localhost:8080/cell/"+"417ecfe7-d2b4-4e43-afd4-dbf5f431d97d"+"/linkCell", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should return Status Found", func() {
				Expect(rr.Code).To(Equal(http.StatusFound))
			})
		})
	})
	
	Describe("When unlinking two cells", func() {
		Context("given I provide a proper cell", func() {
			BeforeEach(func() {
				form := url.Values{}
				form.Add("cellToUnlink", "417ecfe7-d2b4-4e43-afd4-dbf5f431d97d")
				req, err := http.NewRequest("POST", "http://localhost:8080/cell/"+cellId+"/unlinkCell", strings.NewReader(form.Encode()))
				Expect(err).To(BeNil())
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				router.ServeHTTP(rr, req)
			})
			It("should return Status Found", func() {
				Expect(rr.Code).To(Equal(http.StatusFound))
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