package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/PierreBougon/Bym-BackEnd/models"
)

var _ = Describe("Playlists", func() {
	var (
		invalidPlaylist models.Playlist
	)

	BeforeEach(func () {
		invalidPlaylist = models.Playlist{
			Name: "Mu",
			UserId: 0,
			Songs: make([]models.Song, 0),
		}
	})

	Describe("Validating Playlist data", func() {
		Context("With an name shorter than 3 characters", func() {
			wrongName := models.Playlist{
				Name: invalidPlaylist.Name,
				UserId: mockPlaylist.UserId,
				Songs: mockPlaylist.Songs,
			}
			AssertValidationBehavior(&wrongName, false)
		})

		Context("With an userId lesser than 1", func() {
			wrongUser := models.Playlist{
				Name: mockPlaylist.Name,
				UserId: invalidPlaylist.UserId,
				Songs: mockPlaylist.Songs,
			}
			AssertValidationBehavior(&wrongUser, false)
		})

		Context("With correct data", func() {
			AssertValidationBehavior(&mockPlaylist, true)
		})
	})

	Describe("Creating a Playlist", func() {
		Context("With invalid data", func() {
			It("should fail", func() {
				resp := invalidPlaylist.Create(mockAccount.ID)
				Expect(resp["status"]).To(BeFalse())
			})
		})

	})

	Describe("Fetching a playlist", func() {
		Context("With a wrong playlist ID", func() {
			It("should return nothing", func() {
				playlist := models.GetPlaylistById(0)
				Expect(playlist).To(BeNil())
			})
		})

		Context("With a valid playlist ID", func() {
			It("should return a playlist", func() {
				playlist := models.GetPlaylistById(mockPlaylist.ID)
				Expect(playlist).ToNot(BeNil())
			})
		})
	})

	Describe("Fetching all playlists from a user", func() {
		Context("With a correct ID", func() {
			It("should return a list", func() {
				var s interface{} = models.GetPlaylists(mockAccount.ID)
				playlists, ok := s.([]*models.Playlist)

				Expect(ok).To(BeTrue())
				Expect(playlists).ToNot(BeNil())
			})
		})
	})
})
