package models_test

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	mockAccount	models.Account
	mockAccountVoter models.Account
	mockPlaylist models.Playlist
	mockSong models.Song
	mockVote models.Vote
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
	loadAllMockModels()
})

var AssertValidationBehavior = func(t models.Table, success bool) {
	resp, state := t.Validate()
	Expect(state).To(Equal(success), "%s : %+v", resp["message"], t)
}

func loadMockAccount(account *models.Account, email string) {
	mockAccountPassword := "123456"
	err := models.GetDB().
		Table("accounts").
		Where("email = ?", email).
		First(account).Error
	if err != nil {
		account = &models.Account{
			Email: email,
			Password: string([]byte(mockAccountPassword)),
			TokenVersion: 0,
		}
		account.Create()
	}
	account.Password = string([]byte(mockAccountPassword))
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
			Status: "None",
		}
		mockSong.Create(mockAccount.ID)
	}
}

func loadMockVote() {
	err := models.GetDB().
		Table("votes").
		Where("user_id = ? and song_id = ?", mockAccount.ID, mockSong.ID).
		First(&mockVote).Error
	if err != nil {
		mockVote = models.Vote{
			UpVote: true,
			DownVote: false,
			UserId: mockAccount.ID,
			SongId: mockSong.ID,
		}
		models.GetDB().Save(&mockVote)
	}
}

func loadAllMockModels() {
	loadMockAccount(&mockAccount, "test@gmail.com")
	loadMockAccount(&mockAccountVoter, "VoterTest@gmail.com")
	loadMockPlaylist()
	loadMockSong()
	loadMockVote()
}
