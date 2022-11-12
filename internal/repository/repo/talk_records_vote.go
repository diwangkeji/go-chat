package repo

import (
	"context"

	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
)

type TalkRecordsVote struct {
	*Base
	cache *cache.TalkVote
}

type VoteStatistics struct {
	Count   int            `json:"count"`
	Options map[string]int `json:"options"`
}

func NewTalkRecordsVote(base *Base, cache *cache.TalkVote) *TalkRecordsVote {
	return &TalkRecordsVote{Base: base, cache: cache}
}

func (repo *TalkRecordsVote) GetVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
	// 读取缓存
	if uids, err := repo.cache.GetVoteAnswerUser(ctx, vid); err == nil {
		return uids, nil
	}

	uids, err := repo.SetVoteAnswerUser(ctx, vid)
	if err != nil {
		return nil, err
	}

	return uids, nil
}

func (repo *TalkRecordsVote) SetVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
	uids := make([]int, 0)

	err := repo.Db.WithContext(ctx).Table("talk_records_vote_answer").Where("vote_id = ?", vid).Pluck("user_id", &uids).Error

	if err != nil {
		return nil, err
	}

	_ = repo.cache.SetVoteAnswerUser(ctx, vid, uids)

	return uids, nil
}

func (repo *TalkRecordsVote) GetVoteStatistics(ctx context.Context, vid int) (*VoteStatistics, error) {
	value, err := repo.cache.GetVoteStatistics(ctx, vid)
	if err != nil {
		return repo.SetVoteStatistics(ctx, vid)
	}

	statistic := &VoteStatistics{}

	_ = jsonutil.Decode(value, statistic)

	return statistic, nil
}

func (repo *TalkRecordsVote) SetVoteStatistics(ctx context.Context, vid int) (*VoteStatistics, error) {
	var (
		err          error
		vote         *model.TalkRecordsVote
		answerOption map[string]interface{}
		options      = make([]string, 0)
	)

	tx := repo.Db.WithContext(ctx)

	if err = tx.Table("talk_records_vote").First(&vote, vid).Error; err != nil {
		return nil, err
	}

	_ = jsonutil.Decode(vote.AnswerOption, &answerOption)

	err = tx.Table("talk_records_vote_answer").Where("vote_id = ?", vid).Pluck("option", &options).Error
	if err != nil {
		return nil, err
	}

	opts := make(map[string]int)

	for option := range answerOption {
		opts[option] = 0
	}

	for _, option := range options {
		opts[option] += 1
	}

	statistic := &VoteStatistics{
		Options: opts,
		Count:   len(options),
	}

	_ = repo.cache.SetVoteStatistics(ctx, vid, jsonutil.Encode(statistic))

	return statistic, nil
}