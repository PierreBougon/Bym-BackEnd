package models_test

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Songs", func() {
	var (
		song models.Song
		invalidSong models.Song
	)

	BeforeEach(func () {
		song = models.Song{
			Name: "Music",
			PlaylistId: 1,
			ExternalId: "cameFromSomewhere",
		}
		invalidSong = models.Song{
			Name: "",
			PlaylistId: 0,
			ExternalId: "",
		}
	})

	Describe("Validating Song data", func() {
		Context("With an empty name", func() {
			wrongName := models.Song{
				Name: invalidSong.Name,
				PlaylistId: song.PlaylistId,
				ExternalId: song.ExternalId,
			}
			It("should be invalid", func() {
				_, state := wrongName.Validate(mockAccount.ID)
				Expect(state).To(Equal(false))
			})
		})

		Context("With a playlistId equal to 0", func() {
			wrongId := models.Song{
				Name: song.Name,
				PlaylistId: invalidSong.PlaylistId,
				ExternalId: song.ExternalId,
			}
			It("should be invalid", func() {
				_, state := wrongId.Validate(mockAccount.ID)
				Expect(state).To(Equal(false))
			})
		})

		Context("With an empty externalId", func() {
			wrongExternal := models.Song{
				Name: song.Name,
				PlaylistId: song.PlaylistId,
				ExternalId: invalidSong.ExternalId,
			}
			It("should be invalid", func() {
				_, state := wrongExternal.Validate(mockAccount.ID)
				Expect(state).To(Equal(false))
			})

		})

		Context("With correct data", func() {
			It("should be valid", func() {
				_, state := song.Validate(mockAccount.ID)
				Expect(state).To(Equal(true))
			})
		})
	})

	Describe("Creating a Song", func() {
		Context("With invalid data", func() {
			It("should fail", func() {
				resp := invalidSong.Create(mockAccount.ID)
				Expect(resp["status"]).To(BeFalse())
			})
		})
		
		Context("With valid data", func() {
			It("should attribute an id and return the song", func() {
				resp := song.Create(mockAccount.ID)

				Expect(resp["status"]).To(BeTrue())
				Expect(resp["song"]).NotTo(BeNil())
				Expect(song.ID).To(BeNumerically(">", 0))
			})

			AfterEach(func() {
				if song.ID > 0 {
					db := models.GetDB()
					db.Delete(&song)
				}
			})
		})
	})

	Describe("Fetching all songs from a playlist", func() {
		Context("With a playlist ID", func() {
			It("should return a list", func() {
				var s interface{} = models.GetSongs(1)
				songs, ok := s.([]*models.Song)

				Expect(ok).To(BeTrue())
				Expect(songs).ToNot(BeNil())
			})
		})
	})
})
