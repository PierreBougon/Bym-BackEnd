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

func failedMockCreationPanic(mock string) {
	panic("Failed to fetch and/or create a mock record : " + mock)
}

func loadMockAccount(account *models.Account, email string) {
	mockAccountPassword := "123456"
	res := models.GetDB().
		Table("accounts").
		Where("email = ?", email).
		First(account)
	err := res.Error

	if err != nil && !res.RecordNotFound() {
		failedMockCreationPanic("account : " + err.Error())
	} else {
		createdAccount := &models.Account{
			Email: email,
			Password: string([]byte(mockAccountPassword)),
			TokenVersion: 0,
		}
		createdAccount.Create()
		models.GetDB().First(account, createdAccount)
	}
	account.Password = string([]byte(mockAccountPassword))
}

func loadMockPlaylist() {
	mockPlaylist.Name = "MockTest"
	mockPlaylist.UserId = mockAccount.ID
	err := models.GetDB().
		Table("playlists").
		FirstOrCreate(&mockPlaylist).Error

	if err != nil {
		failedMockCreationPanic("playlist : " + err.Error())
	}
}

func loadMockSong() {
	mockSongName := "MockSong"
	res := models.GetDB().
		Table("songs").
		Where("name = ? and playlist_id = ?", mockSongName, mockPlaylist.ID).
		First(&mockSong)
	err := res.Error

	if err != nil && !res.RecordNotFound() {
		failedMockCreationPanic("song : " + err.Error())
	} else {
		createdSong := models.Song{
			Name: mockSongName,
			PlaylistId: mockPlaylist.ID,
			ExternalId: "the id is a lie",
			VoteDown: 42,
			VoteUp: 43,
			Score: 100,
			Status: "None",
		}
		createdSong.Create(mockAccount.ID)
		models.GetDB().First(&mockSong, createdSong)
	}
}

func loadMockVote() {
	res := models.GetDB().
		Table("votes").
		Where("user_id = ? and song_id = ?", mockAccount.ID, mockSong.ID).
		First(&mockVote)
	err := res.Error

	if err != nil && !res.RecordNotFound() {
		failedMockCreationPanic("vote : " + err.Error())
	} else {
		createdVote := models.Vote{
			UpVote: true,
			DownVote: false,
			UserId: mockAccount.ID,
			SongId: mockSong.ID,
		}
		models.GetDB().Create(&createdVote)
		models.GetDB().First(&mockVote, createdVote)
	}
}

func loadAllMockModels() {
	loadMockAccount(&mockAccount, "test@gmail.com")
	loadMockAccount(&mockAccountVoter, "VoterTest@gmail.com")
	loadMockPlaylist()
	loadMockSong()
	loadMockVote()
}
