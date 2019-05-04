package models_test

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	mockAccount	models.Account
	mockPlaylist models.Playlist
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

func loadMockAccount() {
	ret := models.GetUser(6)
	if ret == nil {
		mockAccount = models.Account{
			Email: "test@gmail.com",
			Password: "123456",
			TokenVersion: 0,
		}
		mockAccount.Create()
	} else {
		mockAccount = *ret
	}
}

func loadMockPlaylist() {
	ret := models.GetPlaylistById(6)
	if ret == nil {
		mockPlaylist = models.Playlist{
			Name: "test",
			UserId: mockAccount.ID,
		}
		mockPlaylist.Create(mockAccount.ID)
	} else {
		mockPlaylist = *ret
	}
}

var _ = BeforeSuite(func() {
	loadMockAccount()
	loadMockPlaylist()
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