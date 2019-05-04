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

func LoadMockAccount() {
	ret := models.GetUser(6)
	mockAccountPassword := "123456"
	if ret == nil {
		mockAccount = models.Account{
			Email: "test@gmail.com",
			Password: string([]byte(mockAccountPassword)),
			TokenVersion: 0,
		}
		mockAccount.Create()
	} else {
		mockAccount = *ret
	}
	mockAccount.Password = string([]byte(mockAccountPassword))
}

func loadMockPlaylist() {
	ret := models.GetPlaylistById(15)
	if ret == nil {
		mockPlaylist = models.Playlist{
			Name: "MockTest",
			UserId: mockAccount.ID,
			Songs: make([]models.Song, 0),
		}
		mockPlaylist.Create(mockAccount.ID)
	} else {
		mockPlaylist = *ret
	}
}

var _ = BeforeSuite(func() {
	LoadMockAccount()
	loadMockPlaylist()
})

var AssertValidationBehavior = func(t models.Table, success bool) {
	validity := "invalid"
	if success {
		validity = "valid"
	}
	It("should be " + validity, func() {
		resp, state := t.Validate()
		Expect(state).To(Equal(success), resp["message"])
	})
}