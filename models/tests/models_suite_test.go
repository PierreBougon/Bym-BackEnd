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
	mockSong models.Song
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

func loadMockAccount() {
	mockAccountPassword := "123456"
	mockAccountEmail := "test@gmail.com"
	err := models.GetDB().
		Table("accounts").
		Where("email = ?", mockAccountEmail).
		First(&mockAccount).Error
	if err != nil {
		mockAccount = models.Account{
			Email: mockAccountEmail,
			Password: string([]byte(mockAccountPassword)),
			TokenVersion: 0,
		}
		mockAccount.Create()
	}
	mockAccount.Password = string([]byte(mockAccountPassword))
}

func loadMockPlaylist() {
	mockPlaylistName := "MockTest"
	err := models.GetDB().
		Table("playlists").
		Where("name = ? and user_id = ?", mockPlaylistName, mockAccount.ID).
		First(&mockPlaylist).Error
	if err != nil {
		mockPlaylist = models.Playlist{
			Name: mockPlaylistName,
			UserId: mockAccount.ID,
			Songs: make([]models.Song, 0),
		}
		mockPlaylist.Create(mockAccount.ID)
	}
}

func loadMockSong() {
	mockSongName := "MockSong"
	err := models.GetDB().
		Table("songs").
		Where("name = ? and playlist_id = ?", mockSongName, mockPlaylist.ID).
		First(&mockSong).Error
	if err != nil {
		mockSong = models.Song{
			Name: mockSongName,
			PlaylistId: mockPlaylist.ID,
			ExternalId: "the id is a lie",
			VoteDown: 42,
			VoteUp: 43,
			Score: 100,
		}
		mockSong.Create(mockAccount.ID)
	}
}

func loadAllMockModels() {
	loadMockAccount()
	loadMockPlaylist()
	loadMockSong()
}

var _ = BeforeSuite(func() {
	loadAllMockModels()
})

var AssertValidationBehavior = func(t models.Table, success bool) {
		resp, state := t.Validate()
		Expect(state).To(Equal(success), "%s : %+v", resp["message"], t)
}