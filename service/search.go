package service

import (
	"context"
	searchv1 "github.com/MuxiKeStack/be-api/gen/proto/search/v1"
	"github.com/MuxiKeStack/be-search/domain"
	"github.com/MuxiKeStack/be-search/pkg/logger"
	"github.com/MuxiKeStack/be-search/repository"
	"strings"
	"time"
)

type SearchService interface {
	SearchCourse(ctx context.Context, keyword string, uid int64, location searchv1.SearchLocation) ([]domain.Course, error)
	GetUserSearchHistories(ctx context.Context, uid int64, location searchv1.SearchLocation, offset int64,
		limit int64) ([]domain.SearchHistory, error)
	HideUserSearchHistories(ctx context.Context, uid int64, location searchv1.SearchLocation, removeAll bool, historyIds []int64) error
}

type searchService struct {
	courseRepo  repository.CourseRepository
	historyRepo repository.HistoryRepository
	l           logger.Logger
}

func NewSearchService(courseRepo repository.CourseRepository, historyRepo repository.HistoryRepository, l logger.Logger) SearchService {
	return &searchService{courseRepo: courseRepo, historyRepo: historyRepo, l: l}
}

func (s *searchService) SearchCourse(ctx context.Context, keyword string, uid int64, location searchv1.SearchLocation) ([]domain.Course, error) {
	var (
		courses []domain.Course
		err     error
	)
	// 处理表达式
	if strings.HasPrefix(keyword, "fav:") {
		keywords := strings.Split(strings.TrimPrefix(keyword, "fav:"), " ")
		// 从我的收藏中search
		courses, err = s.courseRepo.Search(ctx, keywords, uid)
	} else {
		keywords := strings.Split(keyword, " ")
		// 从全部中search
		courses, err = s.courseRepo.Search(ctx, keywords, 0)
	}
	if err != nil {
		return nil, err
	}
	// 搜索成功后添加一个历史记录
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := s.historyRepo.Create(ctx, domain.SearchHistory{
			Uid:            uid,
			SearchLocation: location,
			Keyword:        keyword,
			Status:         searchv1.VisibilityStatus_Visible,
		})
		if er != nil {
			s.l.Error("添加历史记录失败",
				logger.Error(err),
				logger.String("keyword", keyword),
				logger.Int64("uid", uid))
		}
	}()
	return courses, nil
}

func (s *searchService) GetUserSearchHistories(ctx context.Context, uid int64, location searchv1.SearchLocation, offset int64, limit int64) ([]domain.SearchHistory, error) {
	return s.historyRepo.GetUserSearchHistories(ctx, uid, location, offset, limit)
}

func (s *searchService) HideUserSearchHistories(ctx context.Context, uid int64, location searchv1.SearchLocation, removeAll bool, historyIds []int64) error {
	if removeAll {
		return s.historyRepo.HideAllUserSearchHistories(ctx, uid, location)
	}
	return s.historyRepo.HideUserSearchHistoriesByIds(ctx, uid, location, historyIds)
}
