package models_test

import (
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
	var AssertValidationBehavior = func(s *models.Account, success bool) {
		validity := "invalid"
		if success {
			validity = "valid"
		}
		It("should be " + validity, func() {
			resp, state := s.Validate()
			Expect(state).To(Equal(success), "%s %+v", resp["message"], s)
		})
	}
	Describe("Validating account data", func() {
		Context("With a wrong email", func() {
			wrongEmail := models.Account{
				Email: invalidAccount.Email,
				Password: mockAccount.Password,
			}
			AssertValidationBehavior(&wrongEmail, false)
		})

		Context("With an already used email", func() {
			wrongEmail := models.Account{
				Email: mockAccount.Email,
				Password: mockAccount.Password,
			}
			AssertValidationBehavior(&wrongEmail, false)
		})

		Context("With a password having less than 6 character", func() {
			wrongPassword := models.Account{
				Email: mockAccount.Email,
				Password: invalidAccount.Password,
			}
			AssertValidationBehavior(&wrongPassword, false)
		})

		Context("With correct data", func() {
			It("should be valid", func() {
				newAccount:= models.Account{
					Email: "NEWtest@gmail.com",
					Password: mockAccount.Password,
				}
				resp, state := newAccount.Validate()
				Expect(state).To(BeTrue(), resp["message"])
			})
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

				Expect(resp["status"]).To(BeFalse(), resp["message"])
				Expect(resp["account"]).To(BeNil())
			})
		})

		Context("With the right credentials", func() {
			It("should return the Account you logged into and assign it a token", func() {
				resp := models.Login(mockAccount.Email, mockAccount.Password)


				Expect(resp["status"]).To(BeTrue(), resp["message"])
				Expect(resp["token"]).ToNot(BeNil())
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
