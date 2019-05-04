package models_test

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	mockAccount	models.Account
	created		bool
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
	ret := models.GetUser(6)
	created = false
	if ret == nil {
		mockAccount := models.Account{
			Email: "test@gmail.com",
			Password: "123456",
		}
		mockAccount.Create()
		created = true
	} else {
		mockAccount = *ret
	}
})

var _ = AfterSuite(func() {
	if created {
		models.GetDB().Delete(mockAccount)
	}
})

var AssertValidationBehavior = func(t models.Table, success bool) {
	validity := "invalid"
	if success {
		validity = "valid"
	}
	It("should be " + validity, func() {
		_, state := t.Validate()
		Expect(state).To(Equal(success))
	})
}