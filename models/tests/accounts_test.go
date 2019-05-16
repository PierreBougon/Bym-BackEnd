package models_test

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Accounts", func() {
	var (
		invalidAccount models.Account

		accountRecordExists = func(account models.Account) bool {
			return !models.GetDB().First(&models.Account{}, account.ID).RecordNotFound()
		}

		cleanAccount = func(account models.Account) {
			if accountRecordExists(account) {
				models.GetDB().
					Delete(&account)
			}
		}
	)

	BeforeEach(func() {
		invalidAccount = models.Account{
			Email: "notAnEmail.fr",
			Password: "hi",
		}
	})

	Describe("Validating account data", func() {
		Context("With a wrong email", func() {
			It("should be invalid", func() {
				AssertValidationBehavior(&models.Account{
					Email: invalidAccount.Email,
					Password: mockAccount.Password,
				}, false)
			})
		})

		Context("With an already used email", func() {
			It("should be invalid", func() {
				AssertValidationBehavior(&models.Account{
					Email: mockAccount.Email,
					Password: mockAccount.Password,
				}, false)
			})
		})

		Context("With a password having less than 6 character", func() {
			It("should be invalid", func() {
				AssertValidationBehavior(&models.Account{
					Email: mockAccount.Email,
					Password: invalidAccount.Password,
				}, false)
			})
		})

		Context("With correct data", func() {
			It("should be valid", func() {
				AssertValidationBehavior(&models.Account{
					Email: "NotDuplicate" + mockAccount.Email,
					Password: mockAccount.Password,
				}, true)
			})
		})
	})

	Describe("Creating an Account", func() {
		var account models.Account

		AfterEach(func() { cleanAccount(account) })

		Context("With invalid data", func() {
			It("should fail", func() {
				resp := invalidAccount.Create()
				Expect(resp["status"]).To(BeFalse())
			})
		})

		Context("With valid data", func() {
			It("should succeed and return the created account", func() {
				account = models.Account{
					Email: "New" + mockAccount.Email,
					Password: "123456",
				}

				resp := account.Create()
				Expect(resp["status"]).To(BeTrue(), "%s %+v", resp["message"], account)

				Expect(resp["token"]).NotTo(BeEmpty())
				Expect(accountRecordExists(account)).To(BeTrue(), "The return record does not exist in the database")
			})
		})
	})

	Describe("Logging into an Account", func() {
		Context("With the a non-existing email", func() {
			It("should fail and return an error message", func() {
				resp := models.Login(invalidAccount.Email, invalidAccount.Password)

				Expect(resp["status"]).To(BeFalse(), resp["message"])
				Expect(resp["account"]).To(BeNil())
			})
		})

		Context("With the an incorrect password", func() {
			It("should fail and return an error message", func() {
				resp := models.Login(mockAccount.Email, invalidAccount.Password)

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

	PDescribe("Updating a Account", func() {
		var oldPass string
		newAccountPassword := "updatedPass"

		BeforeEach(func() {
			oldPass = mockAccount.Password
		})

		AfterEach(func() {
			models.GetDB().First(&mockAccount, mockAccount.ID)
			if mockAccount.Password == newAccountPassword {
				models.GetDB().
					Model(&mockAccount).
					Update("password", oldPass)
				mockAccount.Password = oldPass
			}
		})

		PContext("Which does not belong to the user", func() {
			It("should fail with an error message", func() {
				resp := models.UpdatePassword(mockAccount.ID, newAccountPassword)

				Expect(resp["status"]).To(BeFalse(), "Password should not have been updated")
			})
		})

		Context("Which belongs to the user", func() {
			It("should successfully modify the Account password", func() {
				resp := models.UpdatePassword(mockAccount.ID, newAccountPassword)

				Expect(resp["status"]).To(BeTrue())

				mock := models.Account{}
				err := models.GetDB().Table("accounts").
					Where("id = ?", mockAccount.ID).Find(&mock).Error
				if err != nil {
					Fail(err.Error())
				}
				Expect(mock.Password).To(Equal(newAccountPassword))
			})
		})
	})
})
