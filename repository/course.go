package repository

import (
	"context"
	"github.com/MuxiKeStack/be-search/domain"
	"github.com/MuxiKeStack/be-search/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type courseRepository struct {
	dao dao.CourseDAO
}

func NewCourseRepository(dao dao.CourseDAO) CourseRepository {
	return &courseRepository{dao: dao}
}

func (repo *courseRepository) Search(ctx context.Context, keywords []string, uid int64) ([]domain.Course, error) {
	courses, err := repo.dao.Search(ctx, keywords, uid)
	return slice.Map(courses, func(idx int, src dao.Course) domain.Course {
		return repo.toDomain(src)
	}), err
}

func (repo *courseRepository) InputCourseCompositeScore(ctx context.Context, courseId int64, score float64) error {
	return repo.dao.InputCourseCompositeScore(ctx, courseId, score)
}

func (repo *courseRepository) InputCourse(ctx context.Context, course domain.Course) error {
	return repo.dao.InputCourse(ctx, dao.Course{
		Id:      course.Id,
		Name:    course.Name,
		Teacher: course.Teacher,
	})
}

func (repo *courseRepository) toDomain(course dao.Course) domain.Course {
	return domain.Course{
		Id:             course.Id,
		Name:           course.Name,
		Teacher:        course.Teacher,
		CompositeScore: course.CompositeScore,
	}
}
