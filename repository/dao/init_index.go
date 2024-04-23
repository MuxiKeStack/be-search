package dao

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/olivere/elastic/v7"
	"time"
)

var (
	//go:embed course_index.json
	courseIndex string
)

// InitES 创建索引
func InitES(client *elastic.Client) error {
	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return tryCreateIndex(ctx, client, CourseIndexName, courseIndex)
}

func tryCreateIndex(ctx context.Context,
	client *elastic.Client,
	idxName, idxCfg string,
) error {
	ok, err := client.IndexExists(idxName).Do(ctx)
	if err != nil {
		return fmt.Errorf("检测 %s 是否存在失败 %w", idxName, err)
	}
	if ok {
		return nil
	}
	_, err = client.CreateIndex(idxName).
		Body(idxCfg).Do(ctx)
	if err != nil {
		return fmt.Errorf("初始化 %s 失败 %w", idxName, err)
	}
	return err
}
