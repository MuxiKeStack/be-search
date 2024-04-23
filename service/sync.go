package service

import (
	"context"
	"github.com/MuxiKeStack/be-search/domain"
	"github.com/MuxiKeStack/be-search/repository"
)

type SyncService interface {
	InputCourse(ctx context.Context, course domain.Course) error
	InputCourseCompositeScore(ctx context.Context, courseId int64, score float64) error
}

type syncService struct {
	courseRepo repository.CourseRepository
}

func NewSyncService(courseRepo repository.CourseRepository) SyncService {
	return &syncService{courseRepo: courseRepo}
}

func (s *syncService) InputCourseCompositeScore(ctx context.Context, courseId int64, score float64) error {
	return s.courseRepo.InputCourseCompositeScore(ctx, courseId, score)
}

func (s *syncService) InputCourse(ctx context.Context, course domain.Course) error {
	return s.courseRepo.InputCourse(ctx, course)
}
