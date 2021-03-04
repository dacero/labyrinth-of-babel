package repository_test

import (
	"testing"

	"github.com/dacero/labyrinth-of-babel/models"
	"github.com/dacero/labyrinth-of-babel/repository"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRepository(t *testing.T) {
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
	})
})
