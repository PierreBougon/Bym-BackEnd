package models_test

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Songs", func() {
	var (
		invalidSong models.Song
	)

	BeforeEach(func () {
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
			assertSongValidationBehavior(&mockSong, true)
		})

		Context("With an empty name", func() {
			wrongName := models.Song{
				Name: invalidSong.Name,
				PlaylistId: mockSong.PlaylistId,
				ExternalId: mockSong.ExternalId,
			}
			assertSongValidationBehavior(&wrongName, false)
		})

		Context("With a playlistId equal to 0", func() {
			wrongId := models.Song{
				Name: mockSong.Name,
				PlaylistId: invalidSong.PlaylistId,
				ExternalId: mockSong.ExternalId,
			}
			assertSongValidationBehavior(&wrongId, false)
		})

		Context("With an empty externalId", func() {
			wrongExternal := models.Song{
				Name: mockSong.Name,
				PlaylistId: mockSong.PlaylistId,
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
			var (
				song models.Song
			    resp map[string]interface{}
			    cleanSong = func() {
					if song.ID > 0 {
						db := models.GetDB()
						db.Delete(&song)
					}
				}
			)

			It("should succeed", func() {
				defer cleanSong()
				song = models.Song{
					Name: "New" + mockSong.Name,
					PlaylistId: mockSong.PlaylistId,
					ExternalId: mockSong.ExternalId,
				}

				resp = song.Create(mockAccount.ID)
				Expect(resp["status"]).To(BeTrue(), "%s %+v", resp["message"], song)
			})
			It("should return the created song", func() {
				Expect(resp["song"]).NotTo(BeNil())
				res, ok := (resp["song"]).(*models.Song)

				Expect(ok).To(BeTrue(), "It did not return an instance of Song", resp)
				Expect(res.ID).To(BeNumerically(">", mockSong.ID))
			})
		})
	})

	Describe("Updating a Song", func() {
		var oldName string
		newSongName := "updatedName"
		newSong := models.Song{
			Name: newSongName,
			// other attributes are not used to update
		}

		BeforeEach(func() {
			oldName = mockSong.Name
		})

		AfterEach(func() {
			if mockSong.Name == newSongName {
				models.GetDB().
					Model(&mockSong).
					Update("name", oldName)
			}
		})

		Context("Which does not belong to the user", func() {
			It("should fail with an error message", func() {
				resp := mockSong.UpdateSong(mockPlaylist.UserId + 1, mockSong.ID, &newSong)

				Expect(resp["status"]).To(BeFalse())
			})

		})

		Context("Which belongs to the user", func() {
			It("should successfully modify the Song name", func() {
				resp := mockSong.UpdateSong(mockAccount.ID, mockSong.ID, &newSong)

				Expect(resp["status"]).To(BeTrue())

				mock := models.Song{}
				err := models.GetDB().Table("songs").Where("id = ?", mockSong.ID).Find(&mock).Error
				if err != nil {
					Fail(err.Error())
				}
				Expect(mock.Name).To(Equal(newSongName))
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

	Describe("Fetching songs ranking from a playlist", func() {
		Context("With an invalid (non-existing or deleted) playlist ID", func() {
			It("should return an empty list of type Ranking", func() {
				var s interface{} = models.GetSongsRanking(0)
				songs, ok := s.([]*models.Ranking)

				Expect(ok).To(BeTrue())
				Expect(songs).To(BeEmpty())
			})
		})

		Context("With a valid playlist ID", func() {
			It("should return a list of type Ranking", func() {
				var s interface{} = models.GetSongsRanking(mockPlaylist.ID)
				songs, ok := s.([]*models.Ranking)

				Expect(ok).To(BeTrue())
				Expect(songs).ToNot(BeNil())
			})
		})
	})

	Describe("Fetching a Song ranking data", func() {
		Context("With an invalid Song id", func() {
			It("should return nil", func() {
				ranking := models.GetSongRankingById(0)
				Expect(ranking).To(BeNil())
			})
		})

		Context("With a valid Song id", func() {
			It("should return a *Ranking", func() {
				ranking := models.GetSongRankingById(mockSong.ID)
				Expect(*ranking).ToNot(Equal(models.Ranking{}))
			})
		})
	})
})
