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
			PlaylistId: mockPlaylist.ID,
			ExternalId: "cameFromSomewhere",
			VoteUp: 0,
			VoteDown: 0,
			Score: 0,
		}
		invalidSong = models.Song{
			Name: "",
			PlaylistId: 0,
			ExternalId: "",
		}
	})

	var assertSongValidationBehavior = func(s *models.Song, success bool) {
		validity := "invalid"
		if success {
			validity = "valid"
		}
		It("should be " + validity, func() {
			resp, state := s.Validate(mockAccount.ID)
			Expect(state).To(Equal(success), "%s %+v", resp["message"], s)
		})
	}

	Describe("Validating Song data", func() {
		Context("With correct data", func() {
			assertSongValidationBehavior(&song, true)
		})

		Context("With an empty name", func() {
			wrongName := models.Song{
				Name: invalidSong.Name,
				PlaylistId: song.PlaylistId,
				ExternalId: song.ExternalId,
			}
			assertSongValidationBehavior(&wrongName, false)
		})

		Context("With a playlistId equal to 0", func() {
			wrongId := models.Song{
				Name: song.Name,
				PlaylistId: invalidSong.PlaylistId,
				ExternalId: song.ExternalId,
			}
			assertSongValidationBehavior(&wrongId, false)
		})

		Context("With an empty externalId", func() {
			wrongExternal := models.Song{
				Name: song.Name,
				PlaylistId: song.PlaylistId,
				ExternalId: invalidSong.ExternalId,
			}
			assertSongValidationBehavior(&wrongExternal, false)

		})

	})

	Describe("Creating a Song", func() {
		Context("With invalid data", func() {
			It("should fail", func() {
				resp := invalidSong.Create(mockAccount.ID)
				Expect(resp["status"]).To(BeFalse(), resp["message"])
			})
		})
		
		Context("With valid data", func() {
			AfterEach(func() {
				if song.ID > 0 {
					db := models.GetDB()
					db.Delete(&song)
				}
			})

			It("should attribute an id and return the song", func() {
				resp := song.Create(mockAccount.ID)

				Expect(resp["status"]).To(BeTrue(), "%s %+v", resp["message"], song)
				Expect(resp["song"]).NotTo(BeNil())
				Expect(song.ID).To(BeNumerically(">", 0))
			})

		})
	})

	Describe("Fetching all songs from a playlist", func() {
		Context("With a playlist ID", func() {
			It("should return a list", func() {
				var s interface{} = models.GetSongs(mockPlaylist.ID)
				songs, ok := s.([]*models.Song)

				Expect(ok).To(BeTrue())
				Expect(songs).ToNot(BeNil())
			})
		})
	})
})
