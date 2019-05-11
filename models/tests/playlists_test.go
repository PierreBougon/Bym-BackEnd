package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/PierreBougon/Bym-BackEnd/models"
)

var _ = Describe("Playlists", func() {
	var (
		invalidPlaylist models.Playlist

		playlistRecordExists = func(playlist models.Playlist) bool {
			return !models.GetDB().First(&models.Playlist{}, playlist.ID).RecordNotFound()
		}

		cleanPlaylist = func(playlist models.Playlist) {
			if playlistRecordExists(playlist) {
				models.GetDB().
					Delete(&playlist)
			}
		}
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
			It("should be invalid", func() {
				AssertValidationBehavior(&models.Playlist{
					Name: invalidPlaylist.Name,
					UserId: mockPlaylist.UserId,
					Songs: mockPlaylist.Songs,
				}, false)
			})
		})

		Context("With an userId lesser than 1", func() {
			It("should be invalid", func() {
				AssertValidationBehavior(&models.Playlist{
					Name: mockPlaylist.Name,
					UserId: invalidPlaylist.UserId,
					Songs: mockPlaylist.Songs,
				}, false)
			})
		})

		Context("With correct data", func() {

			It("should be valid", func() {
				AssertValidationBehavior(&mockPlaylist, true)
			})
		})
	})

	Describe("Creating a Playlist", func() {
		var playlist models.Playlist

		AfterEach(func() { cleanPlaylist(playlist) })

		Context("With invalid data", func() {
			It("should fail", func() {
				resp := invalidPlaylist.Create(mockAccount.ID)
				Expect(resp["status"]).To(BeFalse())
			})
		})

		Context("With valid data", func() {
			It("should succeed and return the created playlist", func() {
				playlist = models.Playlist{
					Name: "New" + mockPlaylist.Name,
				}

				resp := playlist.Create(mockAccount.ID)
				Expect(resp["status"]).To(BeTrue(), "%s %+v", resp["message"], playlist)

				Expect(resp["playlist"]).NotTo(BeNil())
				res, ok := (resp["playlist"]).(*models.Playlist)

				Expect(ok).To(BeTrue(), "It did not return an instance of Playlist", resp)
				Expect(playlistRecordExists(*res)).To(BeTrue(), "The return record does not exist in the database")
			})
		})

	})

	Describe("Deleting a Playlist", func() {
		var (
			playlistToDelete models.Playlist

			shouldRecordStillExist = func(shouldExist bool) {
				stillExist := playlistRecordExists(playlistToDelete)
				Expect(stillExist).To(Equal(shouldExist), "%+v", playlistToDelete)
			}
		)

		BeforeEach(func() {
			playlistToDelete = models.Playlist{
				Name: "ToDelete" + mockPlaylist.Name,
				UserId: mockAccount.ID,
			}
			err := models.GetDB().Create(&playlistToDelete).Error
			if err != nil {
				panic("database error : " + err.Error())
			}
		})

		AfterEach(func() { cleanPlaylist(playlistToDelete) })

		Context("Which does not belong to the user", func() {
			It("should fail with an error message", func() {
				resp := mockPlaylist.DeletePlaylist(mockPlaylist.UserId + 1, playlistToDelete.ID)

				Expect(resp["status"]).To(BeFalse(), "playlist was deleted without ownership")

				shouldRecordStillExist(true)
			})
		})

		Context("Which belong to the user", func() {
			It("should succeed", func() {
				resp := mockPlaylist.DeletePlaylist(mockPlaylist.UserId, playlistToDelete.ID)

				Expect(resp["status"]).To(BeTrue(), "playlist was not deleted")

				shouldRecordStillExist(false)
			})
		})
	})

	Describe("Updating a Playlist", func() {
		var oldName string
		newPlaylistName := "updatedName"
		newPlaylist := models.Playlist{
			Name: newPlaylistName,
			// other attributes are not used to update
		}

		BeforeEach(func() {
			oldName = mockPlaylist.Name
		})

		AfterEach(func() {
			models.GetDB().First(&mockPlaylist, mockPlaylist.ID)
			if mockPlaylist.Name == newPlaylistName {
				models.GetDB().
					Model(&mockPlaylist).
					Update("name", oldName)
				mockPlaylist.Name = oldName
			}
		})

		Context("Which does not belong to the user", func() {
			It("should fail with an error message", func() {
				resp := mockPlaylist.UpdatePlaylist(mockPlaylist.UserId + 1, mockPlaylist.ID, &newPlaylist)

				Expect(resp["status"]).To(BeFalse())
			})

		})

		Context("Which belongs to the user", func() {
			It("should successfully modify the Playlist name", func() {
				resp := mockPlaylist.UpdatePlaylist(mockAccount.ID, mockPlaylist.ID, &newPlaylist)

				Expect(resp["status"]).To(BeTrue())

				mock := models.Song{}
				err := models.GetDB().Table("playlists").
					Where("id = ?", mockPlaylist.ID).Find(&mock).Error
				if err != nil {
					Fail(err.Error())
				}
				Expect(mock.Name).To(Equal(newPlaylistName))
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
				var s interface{} = models.GetPlaylistsByUser(mockAccount.ID)
				playlists, ok := s.([]*models.Playlist)

				Expect(ok).To(BeTrue())
				Expect(playlists).ToNot(BeNil())
			})
		})
	})
})
