package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/MuxiKeStack/be-search/domain"
	"github.com/MuxiKeStack/be-search/pkg/canalx"
	"github.com/MuxiKeStack/be-search/pkg/logger"
	"github.com/MuxiKeStack/be-search/pkg/saramax"
	"github.com/MuxiKeStack/be-search/service"
	"strconv"
	"time"
)

type CourseBinlogData struct {
	// course表的字段
	Id      string `json:"id"`
	Name    string `json:"name"`
	Teacher string `json:"teacher"`
	// composite_score表的字段
	CourseId string `json:"course_id"`
	Score    string `json:"score"`
}

type MySQLBinlogConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.SyncService
}

func NewMySQLBinlogConsumer(client sarama.Client, l logger.Logger, svc service.SyncService) *MySQLBinlogConsumer {
	return &MySQLBinlogConsumer{client: client, l: l, svc: svc}
}

func (r *MySQLBinlogConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_courses",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"kstack_binlog"},
			saramax.NewHandler[canalx.Message[CourseBinlogData]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer) Consume(msg *sarama.ConsumerMessage, val canalx.Message[CourseBinlogData]) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	const (
		UPDATE = "UPDATE"
		INSERT = "INSERT"
		DELETE = "DELETE"
	)
	table := val.Table
	switch val.Type {
	case UPDATE, INSERT:
		if table == "courses" {
			for _, data := range val.Data {
				id, err := strconv.ParseInt(data.Id, 10, 64)
				if err != nil {
					return err
				}
				err = r.svc.InputCourse(ctx, domain.Course{
					Id:      id,
					Name:    data.Name,
					Teacher: data.Teacher,
				})
				if err != nil {
					return err
				}
			}
		} else if table == "composite_scores" {
			for _, data := range val.Data {
				courseId, err := strconv.ParseInt(data.CourseId, 10, 64)
				if err != nil {
					return err
				}
				score, err := strconv.ParseFloat(data.Score, 64)
				if err != nil {
					return err
				}
				err = r.svc.InputCourseCompositeScore(ctx, courseId, score)
				if err != nil {
					return err
				}
			}
		}
	case DELETE:
		if table == "courses" {
			for _, data := range val.Data {
				courseId, err := strconv.ParseInt(data.Id, 10, 64)
				if err != nil {
					return err
				}
				err = r.svc.DelCourse(ctx, courseId)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
