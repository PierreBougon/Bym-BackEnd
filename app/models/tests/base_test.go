package models_test

import (
	"github.com/PierreBougon/Bym-BackEnd/app/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Base", func() {
	Describe("Connexion test", func() {
		Context("When using the db handle", func() {
			It("should be connected to a database", func() {
				db := models.GetDB()
				Expect(db).NotTo(BeNil())
			})
		})
	})

})
