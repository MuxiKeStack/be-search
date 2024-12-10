package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const CourseIndexName = "course_index"

type Course struct {
	Id             int64   `json:"id"`
	Name           string  `json:"name"`
	Teacher        string  `json:"teacher"`
	CompositeScore float64 `json:"composite_score"`
}

type CourseElasticDAO struct {
	client *elastic.Client
}

func (dao *CourseElasticDAO) DelCourse(ctx context.Context, courseId int64) error {
	_, err := dao.client.Delete().
		Index(CourseIndexName).
		Id(strconv.FormatInt(courseId, 10)).
		Do(ctx)
	return err
}

func NewCourseElasticDAO(client *elastic.Client) CourseDAO {
	return &CourseElasticDAO{client: client}
}

// InputCourse input除了综合分数以外的内容
func (dao *CourseElasticDAO) InputCourse(ctx context.Context, course Course) error {
	_, err := dao.client.Index().
		Index(CourseIndexName).
		Id(strconv.FormatInt(course.Id, 10)).
		BodyJson(course).
		Do(ctx)
	return err
}

func (dao *CourseElasticDAO) InputCourseCompositeScore(ctx context.Context, courseId int64, score float64) error {
	_, err := dao.client.Update().
		Index(CourseIndexName).
		Id(strconv.FormatInt(courseId, 10)).
		Doc(map[string]any{
			"composite_score": score,
		}).
		DetectNoop(true).
		Do(ctx)
	return err
}

func (dao *CourseElasticDAO) Search(ctx context.Context, keywords []string, uid int64) ([]Course, error) {
	queryString := strings.Join(keywords, " ")
	query := elastic.NewBoolQuery().Should(
		elastic.NewMatchQuery("name", queryString).Boost(1.2),
		elastic.NewMatchQuery("teacher", queryString),
	)
	resp, err := dao.client.Search(CourseIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	courses := make([]Course, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var ele Course
		er := json.Unmarshal(hit.Source, &ele)
		if er != nil {
			return nil, er
		}
		courses = append(courses, ele)
	}
	return courses, nil
}
