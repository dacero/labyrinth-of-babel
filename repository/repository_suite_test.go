package repository_test

import (
	"testing"
	"log"
	"os"
	"database/sql"
	"strings"
	"io/ioutil"

	"github.com/dacero/labyrinth-of-babel/models"
	"github.com/dacero/labyrinth-of-babel/repository"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func resetDB() {
	log.Print("Initializing db... ")
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	db, err := sql.Open("mysql", "root:"+password+"@tcp(mysql:3306)/")
	if err != nil {
		log.Fatal("Error when opening DB: ", err)
	}
	defer db.Close()

	file, err := ioutil.ReadFile("./db/test.sql")
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
	RunSpecs(t, "Repository Suite")
}

var _ = Describe("Repository", func() {
	var (
		cell    	models.Cell
		newCell		models.Cell
		newCellId	string
		lobRepo 	repository.LobRepository
		err     	error
		cellId  = "72aed05b-cb2d-4cad-bf70-05d8ae02a7bc"
	)

	BeforeEach(func() {
		lobRepo = repository.NewLobRepository()
	})

	Describe("Retrieving a cell", func() {
		Context("that exists", func() {

			BeforeEach(func() {
				cell, err = lobRepo.GetCell(cellId)
			})

			It("should return that cell and no error", func() {
				Expect(err).To(BeNil())
				Expect(cell.Id).To(Equal(cellId))
			})

			It("should contain all sources", func() {
				log.Print(cell.Sources)
				Expect(len(cell.Sources)).To(Equal(2))
			})

			It("should contain all links", func() {
				Expect(len(cell.Links)).To(Equal(2))
			})
		})

		Context("that does NOT exist", func() {
			BeforeEach(func() {
				cell, err = lobRepo.GetCell("Inexistent cell")
			})
			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
	
	Describe("When I search for cells", func() {
		Context("given I provide a term only used in one", func() {
			It("should return only one cell", func() {
				term := "shorter"
				cells := lobRepo.SearchCells(term)
				Expect(len(cells)).To(Equal(1))
			})
		})
		Context("given I provide a term only used in all", func() {
			It("should return three cells", func() {
				term := "idea"
				cells := lobRepo.SearchCells(term)
				Expect(len(cells)).To(Equal(3))
				log.Print(cells)
			})
		})
	})
	
	Describe("Creating a new cell", func() {
		Context("with proper information", func() {
			BeforeEach(func() {
				newCell = models.Cell{Title: "This is a new cell",
									Body: "This is the body of the new cell",
									Room: "This is a room",
									Sources: []models.Source{ models.Source{Source:"Confucius"} } }
				newCellId, err = lobRepo.NewCell(newCell)
			})

			It("should return a new cell id and no error", func() {
				Expect(err).To(BeNil())
				Expect(len(newCellId)).To(Equal(len("72aed05b-cb2d-4cad-bf70-05d8ae02a7bc")))
			})
			It("should insert a cell in the repository, including sources", func() {
				cell, err = lobRepo.GetCell(newCellId)
				Expect(err).To(BeNil())
				Expect(cell.Body).To(Equal(newCell.Body))
				sources := cell.Sources
				Expect(len(sources)).To(Equal(len(newCell.Sources)))
				Expect(sources[0]).To(Equal(newCell.Sources[0]))
			})
		})
		Context("with a new room", func() {
			BeforeEach(func() {
				newCell = models.Cell{Title: "This is a new cell",
									Body: "This is the body of the new cell",
									Room: "Habitacion 2",
									Sources: []models.Source{ models.Source{Source:"Confucius"} } }
				newCellId, err = lobRepo.NewCell(newCell)
			})
			It("should return a new cell id and no error", func() {
				Expect(err).To(BeNil())
				Expect(len(newCellId)).To(Equal(len("72aed05b-cb2d-4cad-bf70-05d8ae02a7bc")))
			})
		})
		Context("with an empty room", func() {
			BeforeEach(func() {
				newCell = models.Cell{Title: "This is a new cell",
									Body: "This is the body of the new cell",
									Room: "",
									Sources: []models.Source{ models.Source{Source:"Confucius"} } }
				newCellId, err = lobRepo.NewCell(newCell)
			})
			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})
		Context("with an empty body", func() {
			BeforeEach(func() {
				newCell = models.Cell{Title: "This is a new cell",
									Body: "",
									Room: "This is a room",
									Sources: []models.Source{ models.Source{Source:"Confucius"} } }
				newCellId, err = lobRepo.NewCell(newCell)
			})
			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})
	})
	
	Describe("Searching a source", func() {
		Context("with existing terms", func() {
			It("should return an array of elements", func() {
				foundSources := lobRepo.SearchSources("Confu")
				Expect(len(foundSources)).To(Equal(1))
			})
		})
		Context("with inexisting terms", func() {
			It("should return an empty array", func() {
				foundSources := lobRepo.SearchSources("dshfksjfh")
				Expect(len(foundSources)).To(Equal(0))
			})
		})
	})
	
	Describe("Searching a room", func() {
		Context("with existing terms", func() {
			It("should return an array of elements", func() {
				foundRooms := lobRepo.SearchRooms("Habita")
				Expect(len(foundRooms)).To(Equal(1))
			})
		})
		Context("with inexisting terms", func() {
			It("should return an empty array", func() {
				foundRooms := lobRepo.SearchRooms("dshfksjfh")
				Expect(len(foundRooms)).To(Equal(0))
			})
		})
	})
	
	Describe("Updating a cell", func() {
		Context("with good information", func() {
			It("should return no error", func() {
				updateCell := models.Cell{Id: cellId,
					Title: "This is being updated",
					Body: "The new cell being updated",
					Room: "This is a room"}
				update, err := lobRepo.UpdateCell(updateCell)
				Expect(err).To(BeNil())
				Expect(update).To(Equal(int64(1)))
			})
		})
		Context("with a new room", func() {
			It("should create the room and return no error", func() {
				updateCell := models.Cell{Id: cellId,
					Title: "This is being updated",
					Body: "The new cell being updated",
					Room: "Should create this room"}
				update, err := lobRepo.UpdateCell(updateCell)
				Expect(err).To(BeNil())
				Expect(update).To(Equal(int64(1)))
			})
		})
		Context("with an empty room", func() {
			It("should return error", func() {
				updateCell := models.Cell{Id: cellId,
					Title: "This is being updated",
					Body: "The new cell being updated",
					Room: " "} //using empty spaces to test trimming
				update, err := lobRepo.UpdateCell(updateCell)
				Expect(err).ToNot(BeNil())
				Expect(update).To(Equal(int64(0)))
			})
		})
		Context("with an empty body", func() {
			It("should return error", func() {
				updateCell := models.Cell{Id: cellId,
					Title: "This is being updated",
					Body: "    ",  //using empty spaces to test trimming
					Room: "This is a room"}
				update, err := lobRepo.UpdateCell(updateCell)
				Expect(err).ToNot(BeNil())
				Expect(update).To(Equal(int64(0)))
			})
		})
	})
	
	Describe("When I add sources to a cell", func() {
		var newSource models.Source
		Context("given I add a propoer new source", func() {
			BeforeEach(func() {
				newSource = models.Source{Source:"New Source"}
				cell, err = lobRepo.AddSourceToCell(cellId, newSource)
			})
			It("should return no error", func() {
				Expect(err).To(BeNil())
			})
			It("should add the new source to the cell", func() {
				Expect(len(cell.Sources)).To(Equal(3))
			})
		})
		Context("given I add a duplicated source", func() {
			BeforeEach(func() {
				newSource = models.Source{Source:"Confucius"}
				cell, err = lobRepo.AddSourceToCell(cellId, newSource)
			})
			It("should return no error", func() {
				Expect(err).To(BeNil())
			})
			It("should NOT add the new source to the cell", func() {
				Expect(len(cell.Sources)).To(Equal(3))
			})
		})
		Context("given I add an empty source", func() {
			BeforeEach(func() {
				newSource = models.Source{Source:"    "} //spaces to test trimming
				cell, err = lobRepo.AddSourceToCell(cellId, newSource)
			})
			It("should return no error", func() {
				Expect(err).To(BeNil())
			})
			It("should NOT add the new source to the cell", func() {
				Expect(len(cell.Sources)).To(Equal(3))
			})
		})
	})
	Describe("When I remove a source from a cell", func() {
		Context("given I provide a proper source", func() {
			BeforeEach(func() {
				source := models.Source{Source:"Confucius"}
				cell, err = lobRepo.RemoveSourceFromCell(cellId, source)
			})
			It("should return no error", func() {
				Expect(err).To(BeNil())
			})
			It("should remove the source from the cell", func() {
				Expect(len(cell.Sources)).To(Equal(2))
			})
		})
		Context("given I provide a source not linked to the cell", func() {
			BeforeEach(func() {
				source := models.Source{Source:"A different source"}
				cell, err = lobRepo.RemoveSourceFromCell(cellId, source)
			})
			It("should return no error", func() {
				Expect(err).To(BeNil())
			})
			It("should return the cell unchanged", func() {
				Expect(len(cell.Sources)).To(Equal(2))
			})
		})
	})
	
	Describe("When I link 2 cells", func() {
		var cellA string
		var cellB string
		Context("given they exist and are not linked", func() {
			BeforeEach(func() {
				cellA = "417ecfe7-d2b4-4e43-afd4-dbf5f431d97d"
				cellB = "df38bd04-0ec4-41bf-9e53-d0eeb95a4939"
				err = lobRepo.LinkCells(cellA, cellB)
			})
			It("should return no error", func() {
				Expect(err).To(BeNil())
			})
			It("should link the 2 cells", func() {
				areLinked, err := lobRepo.CheckLink(cellA, cellB)
				Expect(err).To(BeNil())
				Expect(areLinked).To(Equal(true))
			})
		})
		Context("given they are already linked", func() {
			BeforeEach(func() {
				cellA = "417ecfe7-d2b4-4e43-afd4-dbf5f431d97d"
				cellB = "72aed05b-cb2d-4cad-bf70-05d8ae02a7bc"
				err = lobRepo.LinkCells(cellA, cellB)
			})
			It("should return error", func() {
				Expect(err).ToNot(BeNil())
			})
		})
	})
	
	Describe("When I unlink 2 cells", func() {
		var cellA string
		var cellB string
		Context("given they exist and are linked", func() {
			BeforeEach(func() {
				cellA = "417ecfe7-d2b4-4e43-afd4-dbf5f431d97d"
				cellB = "72aed05b-cb2d-4cad-bf70-05d8ae02a7bc"
				err = lobRepo.UnlinkCells(cellA, cellB)
			})
			It("should return no error", func() {
				Expect(err).To(BeNil())
			})
			It("should link the 2 cells", func() {
				areLinked, err := lobRepo.CheckLink(cellA, cellB)
				Expect(err).To(BeNil())
				Expect(areLinked).To(Equal(false))
			})
		})
	})
})
