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

	Describe("Validating Song info", func() {
		AssertFailedValidationBehavior := func(s models.Song) {
			It("should be invalid", func() {
				_, state := s.Validate()
				Expect(state).To(BeFalse())
			})
		}

		Context("With an empty name", func() {
			song.Name = ""
			AssertFailedValidationBehavior(song)
		})

		Context("With a songId equal to 0", func() {
			song.PlaylistId = 0
			AssertFailedValidationBehavior(song)
		})

		Context("With an empty externalId", func() {
			song.ExternalId = ""
			AssertFailedValidationBehavior(song)
		})

		Context("With correct data", func() {
			It("should be valid", func() {
				_, state := song.Validate()
				Expect(state).To(BeTrue())
			})
		})
	})

	Describe("Creating a Song", func() {
		Context("With invalid data", func() {
			It("should fail", func() {
				resp := invalidSong.Create()
				Expect(resp["status"]).To(BeFalse())
			})
		})
		
		Context("With valid data", func() {
			It("should attribute an id and return the song", func() {
				resp := song.Create()

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
				songs := models.GetSongs(1)
				Expect(songs).ToNot(BeNil())
			})
		})
	})
})
