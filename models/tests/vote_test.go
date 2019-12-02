package models_test

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	"github.com/PierreBougon/Bym-BackEnd/websocket"
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
				res, ok := t.([]*models.Vote)
				Expect(ok).To(BeTrue())
				for _, vote := range res {
					Expect(vote).ToNot(Equal(models.Vote{}))
				}
			})
		})
	})

	Describe("Updating the vote on a song", func() {
		var switchMockVoteValueToCorrespondTo = func(vote models.Vote, upVote bool) {
			models.GetDB().Find(&vote, vote.ID)
			if vote.UpVote != upVote {
				vote.UpVote = upVote
				vote.DownVote = !upVote
				models.GetDB().Save(vote)
			}
		}

		Context("which does not exist", func() {
			It("should fail with an upVote", func() {
				res := models.UpVoteSong(0, mockAccount.ID, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
				Expect(res["status"]).To(BeFalse())
			})

			It("should fail with an downVote", func() {
				res := models.DownVoteSong(0, mockAccount.ID, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
				Expect(res["status"]).To(BeFalse())
			})
		})

		Context("which exist but already has a corresponding vote of the user on", func() {
			It("should succeed with an upVote", func() {
				switchMockVoteValueToCorrespondTo(mockVote, true)
				res := models.UpVoteSong(mockSong.ID, mockAccount.ID, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
				Expect(res["status"]).To(BeTrue())
			})

			It("should succeed with an downVote", func() {
				switchMockVoteValueToCorrespondTo(mockVote, false)
				res := models.DownVoteSong(mockSong.ID, mockAccount.ID, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
				Expect(res["status"]).To(BeTrue())
			})
		})

		Context("on which there was no vote beforehand", func() {
			var vote models.Vote

			var fetchVote = func(v *models.Vote) bool {
				res := models.GetDB().
					Find(v, "user_id = ? and song_id = ?", v.UserId, v.SongId)
				if !res.RecordNotFound() {
					return true
				} else if res.Error != nil {
					panic("connection error")
				}
				return false
			}

			BeforeEach(func() {
				vote = models.Vote{
					UserId: mockAccountVoter.ID,
					SongId: mockSong.ID,
				}
			})

			AfterEach(func() {
				if fetchVote(&vote) {
					models.GetDB().Delete(&vote)
				}
			})

			It("should create a new Vote with an upVote", func() {
				res := models.UpVoteSong(mockSong.ID, mockAccountVoter.ID, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
				Expect(res["status"]).To(BeTrue())

				found := fetchVote(&vote)
				Expect(found).To(BeTrue(), "The upVote was not created")
				Expect(vote.UpVote).To(BeTrue())
			})

			It("should create a new Vote with an downVote", func() {
				res := models.DownVoteSong(mockSong.ID, mockAccountVoter.ID, websocket.NotifyPlaylistSubscribers, websocket.PlaylistNeedRefresh)
				Expect(res["status"]).To(BeTrue())

				found := fetchVote(&vote)
				Expect(found).To(BeTrue(), "The downVote was not created")
				Expect(vote.DownVote).To(BeTrue())
			})
		})
	})

})
