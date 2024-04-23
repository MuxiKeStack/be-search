package repository

import (
	"context"
	searchv1 "github.com/MuxiKeStack/be-api/gen/proto/search/v1"
	"github.com/MuxiKeStack/be-search/domain"
	"github.com/MuxiKeStack/be-search/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type historyRepository struct {
	dao dao.HistoryDAO
}

func NewHistoryRepository(dao dao.HistoryDAO) HistoryRepository {
	return &historyRepository{dao: dao}
}

func (repo *historyRepository) Create(ctx context.Context, history domain.SearchHistory) error {
	return repo.dao.Insert(ctx, repo.toEntity(history))
}

func (repo *historyRepository) GetUserSearchHistories(ctx context.Context, uid int64, location searchv1.SearchLocation,
	offset int64, limit int64) ([]domain.SearchHistory, error) {
	histories, err := repo.dao.GetUserSearchHistories(ctx, uid, int32(location), offset, limit)
	return slice.Map(histories, func(idx int, src dao.SearchHistory) domain.SearchHistory {
		return repo.toDomain(src)
	}), err
}

func (repo *historyRepository) HideAllUserSearchHistories(ctx context.Context, uid int64, location searchv1.SearchLocation) error {
	return repo.dao.HideAllUserSearchHistories(ctx, uid, int32(location))
}

func (repo *historyRepository) HideUserSearchHistoriesByIds(ctx context.Context, uid int64, location searchv1.SearchLocation,
	historyIds []int64) error {
	return repo.dao.HideUserSearchHistoriesByIds(ctx, uid, int32(location), historyIds)
}

func (repo *historyRepository) toEntity(history domain.SearchHistory) dao.SearchHistory {
	return dao.SearchHistory{
		Id:       history.Id,
		Uid:      history.Uid,
		Location: int32(history.SearchLocation),
		Keyword:  history.Keyword,
		Status:   int32(history.Status),
	}
}

func (repo *historyRepository) toDomain(history dao.SearchHistory) domain.SearchHistory {
	return domain.SearchHistory{
		Id:             history.Id,
		Uid:            history.Uid,
		SearchLocation: searchv1.SearchLocation(history.Location),
		Keyword:        history.Keyword,
		Status:         searchv1.VisibilityStatus(history.Status),
	}
}
