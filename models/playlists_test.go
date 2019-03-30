package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/PierreBougon/Bym-BackEnd/models"
)

var _ = Describe("Playlists", func() {
	var (
		playlist models.Playlist
		invalidPlaylist models.Playlist
	)

	BeforeEach(func () {
		playlist = models.Playlist{
			Name: "Music",
			UserId: 1,
			Songs: make([]models.Song, 0),
		}
		invalidPlaylist = models.Playlist{
			Name: "Mu",
			UserId: 0,
			Songs: make([]models.Song, 0),
		}
	})

	Describe("Validating Playlist info", func() {
		AssertFailedValidationBehavior := func(p models.Playlist) {
			It("should be invalid", func() {
				_, state := p.Validate()
				Expect(state).To(BeFalse())
			})
		}

		Context("With an name shorter than 3 characters", func() {
			playlist.Name = "12"
			AssertFailedValidationBehavior(playlist)
		})

		Context("With an userId lesser than 1", func() {
			playlist.UserId = 0
			AssertFailedValidationBehavior(playlist)
		})

		Context("With correct data", func() {
			It("should be valid", func() {
				_, state := playlist.Validate()
				Expect(state).To(BeTrue())
			})
		})
	})

	Describe("Creating a Playlist", func() {
		created := false
		mockUser := &models.Account{}

		BeforeEach(func() {
			mockUser = models.GetUser(1)
			if mockUser == nil {
				mockUser = &models.Account{
					Email: "test@gmail.com",
					Password: "123456",
				}
				mockUser.Create()
				created = true
			}
		})

		AfterEach(func() {
			if created {
				models.GetDB().Delete(mockUser)
			}
		})

		Context("With invalid data", func() {
			It("should fail", func() {
				resp := invalidPlaylist.Create(mockUser.ID)
				Expect(resp["status"]).To(BeFalse())
			})

		})

		Context("With valid data", func() {
			It("should attribute an id and return the playlist", func() {
				resp := playlist.Create(mockUser.ID)

				Expect(resp["status"]).To(BeTrue())
				Expect(resp["playlist"]).NotTo(BeNil())
				Expect(playlist.ID).To(BeNumerically(">", 0))
			})

			AfterEach(func() {
				if playlist.ID > 0 {
					db := models.GetDB()
					db.Delete(&playlist)
				}
			})
		})
	})
/*
	Describe("Fetching all playlists from a playlist", func() {
		Context("With a playlist ID", func() {
			It("should return a list", func() {
				playlists := models.GetPlaylists(1)
				Expect(playlists).ToNot(BeNil())
			})
		})
	})*/
})
