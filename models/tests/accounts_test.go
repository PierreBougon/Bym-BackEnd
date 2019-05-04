package models_test

import (
	"fmt"
	"github.com/PierreBougon/Bym-BackEnd/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Accounts", func() {
	var invalidAccount models.Account

	BeforeEach(func() {
		invalidAccount = models.Account{
			Email: "notAnEmail.fr",
			Password: "hi",
		}
	})

	Describe("Validating account data", func() {
		Context("With a wrong email", func() {
			wrongName := models.Account{
				Email: invalidAccount.Email,
				Password: mockAccount.Password,
			}
			AssertValidationBehavior(&wrongName, false)
		})

		Context("With a password having less than 6 character", func() {
			wrongPassword := models.Account{
				Email: mockAccount.Email,
				Password: invalidAccount.Password,
			}
			AssertValidationBehavior(&wrongPassword, false)
		})

		Context("With correct data", func() {
			AssertValidationBehavior(&mockAccount, true)
		})
	})

	Describe("Creating an Account", func() {
		Context("With invalid data", func() {
			It("should fail", func() {
				resp := invalidAccount.Create()
				Expect(resp["status"]).To(BeFalse())
			})
		})
	})


	Describe("Logging into an Account", func() {
		Context("With the wrong credentials", func() {
			It("should fail and return an error message", func() {
				resp := models.Login(invalidAccount.Email, invalidAccount.Password)

				Expect(resp["status"]).To(BeFalse())
				Expect(resp["account"]).To(BeNil())
			})
		})

		Context("With the right credentials", func() {
			It("should return the Account you logged into and assign it a token", func() {
				resp := models.Login(mockAccount.Email, mockAccount.Password)
				fmt.Println(resp)
				Expect(resp["status"]).To(BeTrue())
				Expect(resp["account"]).ToNot(BeNil())

				account, ok := resp["account"].(models.Account)

				Expect(ok).To(BeTrue())
				Expect(account.Token).ToNot(BeNil())
			})
		})
	})

	Describe("Fetching an user account", func() {
		Context("With an unknown account id", func() {
			It("should return nothing", func() {
				account := models.GetUser(0)
				Expect(account).To(BeNil())
			})
		})

		Context("With an existing account id", func() {
			It("should return an account", func() {
				account := models.GetUser(mockAccount.ID)
				Expect(account).ToNot(BeNil())
				Expect(account.ID).To(Equal(mockAccount.ID))

			})
		})
	})
})
