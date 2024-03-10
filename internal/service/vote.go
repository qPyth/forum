package service

import (
	"forum/internal/models"
	s "forum/internal/storage"
)

type VoteService struct {
	PostStorage    s.Posts
	CommentStorage s.Comments
	VotesStorage   s.Votes
}

func NewVoteService(postStorage s.Posts, commentStorage s.Comments, votesStorage s.Votes) *VoteService {
	return &VoteService{PostStorage: postStorage, CommentStorage: commentStorage, VotesStorage: votesStorage}
}

func (v *VoteService) MakeVote(input VoteInput) (*VoteOutput, error) {
	vote := models.Vote{
		UserID: input.UserID,
		Action: input.Action,
		ItemID: input.ID,
	}

	if input.IsPost {
		vote.Item = "post"
		if !v.PostStorage.PostExist(input.ID) {
			return nil, models.ErrWrongVoteItem
		}
	} else if input.IsComment {
		vote.Item = "comment"
		exist, err := v.CommentStorage.CommentExist(input.ID)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, err
		}
	} else {
		return nil, models.ErrWrongVoteItem
	}
	var repoAction string
	var oppositeRepoAction string
	if vote.Action == "like" {
		repoAction = "like_count"
		oppositeRepoAction = "dislike_count"
	} else if vote.Action == "dislike" {
		repoAction = "dislike_count"
		oppositeRepoAction = "like_count"
	} else {
		return nil, models.ErrWrongAction
	}

	voteInDB, err := v.VotesStorage.IsExists(vote)
	if err != nil {
		return nil, err
	}
	if voteInDB == nil {
		if err = v.VotesStorage.Create(vote); err != nil {
			return nil, err
		}
		// TODO need switch for comments
		err = v.updateVotes(vote.ItemID, repoAction, 1, vote.Item)
		if err != nil {
			return nil, err
		}
	} else {
		if voteInDB.Action == vote.Action {
			err = v.VotesStorage.Delete(voteInDB.ID)
			if err != nil {
				return nil, err
			}
			err = v.updateVotes(vote.ItemID, repoAction, -1, vote.Item)
			if err != nil {
				return nil, err
			}
			vote.Action = "none"
		} else {
			err = v.VotesStorage.Update(voteInDB.ID, vote.Action)
			if err != nil {
				return nil, err
			}
			err = v.updateVotes(vote.ItemID, oppositeRepoAction, -1, vote.Item)
			if err != nil {
				return nil, err
			}

			err = v.updateVotes(vote.ItemID, repoAction, 1, vote.Item)
			if err != nil {
				return nil, err
			}
		}
		
	}
	//
	var output VoteOutput
	output.Action = vote.Action
	if vote.Item == "post" {
		post, err := v.PostStorage.GetPostByID(vote.ItemID)
		if err != nil {
			return nil, err
		}
		output.LikeCount = post.LikeCount
		output.DislikeCount = post.DislikeCount
	} else {
		comment, err := v.CommentStorage.GetByID(vote.ItemID)
		if err != nil {
			return nil, err
		}
		output.LikeCount = comment.LikeCount
		output.DislikeCount = comment.DislikeCount
	}

	return &output, nil
}

func (v *VoteService) updateVotes(itemID int, action string, delta int, item string) error {
	if item == "post" {
		err := v.PostStorage.UpdateVotes(itemID, action, delta)
		if err != nil {
			return err
		}
	} else {
		err := v.CommentStorage.UpdateVotes(itemID, action, delta)
		if err != nil {
			return err
		}
	}
	return nil
}
