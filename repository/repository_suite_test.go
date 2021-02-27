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

var _ = Describe("Reposiroty", func() {
	var (
		cell    models.Cell
		lobRepo repository.LobRepository
		err     error
		cellId  = "417ecfe7-d2b4-4e43-afd4-dbf5f431d97d"
	)

	BeforeEach(func() {
		lobRepo = repository.NewLobRepository()
	})

	Describe("Retrieving a cell", func() {
		Context("that exists", func() {
			It("should return that cell and no error", func() {
				cell, err = lobRepo.GetCell(cellId)
				Expect(err).To(BeNil())
				Expect(cell.Id).To(Equal(cellId))
			})
		})
	})
})
