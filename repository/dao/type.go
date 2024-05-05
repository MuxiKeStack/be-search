package dao

import "context"

type CourseDAO interface {
	InputCourse(ctx context.Context, course Course) error
	InputCourseCompositeScore(ctx context.Context, courseId int64, score float64) error
	Search(ctx context.Context, keywords []string, uid int64) ([]Course, error)
	DelCourse(ctx context.Context, courseId int64) error
}

type HistoryDAO interface {
	Insert(ctx context.Context, history SearchHistory) error
	GetUserSearchHistories(ctx context.Context, uid int64, location int32, offset int64, limit int64) ([]SearchHistory, error)
	HideAllUserSearchHistories(ctx context.Context, uid int64, location int32) error
	HideUserSearchHistoriesByIds(ctx context.Context, uid int64, location int32, historyIds []int64) error
}
