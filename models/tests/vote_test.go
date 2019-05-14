package models_test

import (
	"fmt"
	"github.com/PierreBougon/Bym-BackEnd/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vote", func() {
	Describe("Fetching a personal vote on a song", func() {
		Context("on which the user did not vote", func() {
			It("should fail", func() {
				res := models.GetPersonalVoteBySongId(0, mockAccount.ID)
				Expect(res).To(BeNil())
			})
		})

		Context("on which the user did vote", func() {
			It("should succeed and return said Vote", func() {
				var t interface{}

				t = models.GetPersonalVoteBySongId(mockSong.ID, mockAccount.ID)

				res, ok := t.(*models.Vote)
				Expect(ok).To(BeTrue())
				Expect(res).ToNot(Equal(models.Vote{}))
			})
		})
	})

	Describe("Fetching all votes on a song", func() {
		Context("on a non-existing song", func() {
			It("should fail", func() {
				res := models.GetVotesBySongId(0)
				Expect(res).To(BeEmpty())
			})
		})

		Context("on an existing song", func() {
			It("should return an array of Vote", func() {
				var t interface{}

				t = models.GetVotesBySongId(mockSong.ID)
				res, _ := t.([]*models.Vote)
				// Expect(ok).To(BeTrue())
				for _, vote := range res {
					fmt.Println(vote)
					Expect(vote).ToNot(Equal(models.Vote{}))
				}
			})
		})
	})



})
