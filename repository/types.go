package repository

import (
	"context"
	searchv1 "github.com/MuxiKeStack/be-api/gen/proto/search/v1"
	"github.com/MuxiKeStack/be-search/domain"
)

type CourseRepository interface {
	InputCourse(ctx context.Context, course domain.Course) error
	InputCourseCompositeScore(ctx context.Context, courseId int64, score float64) error
	// Search uid使用来限制范围的,搜某个用户或者搜索全部
	Search(ctx context.Context, keywords []string, uid int64) ([]domain.Course, error)
	DelCourse(ctx context.Context, courseId int64) error
}

type HistoryRepository interface {
	Create(ctx context.Context, history domain.SearchHistory) error
	GetUserSearchHistories(ctx context.Context, uid int64, location searchv1.SearchLocation, offset int64, limit int64) ([]domain.SearchHistory, error)
	HideAllUserSearchHistories(ctx context.Context, uid int64, location searchv1.SearchLocation) error
	HideUserSearchHistoriesByIds(ctx context.Context, uid int64, location searchv1.SearchLocation, historyIds []int64) error
}
