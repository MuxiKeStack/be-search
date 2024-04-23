//go:build wireinject

package main

import (
	"github.com/MuxiKeStack/be-search/events"
	"github.com/MuxiKeStack/be-search/grpc"
	"github.com/MuxiKeStack/be-search/ioc"
	"github.com/MuxiKeStack/be-search/repository"
	"github.com/MuxiKeStack/be-search/repository/dao"
	"github.com/MuxiKeStack/be-search/service"
	"github.com/google/wire"
)

func InitApp() *App {
	wire.Build(
		wire.Struct(new(App), "*"),
		// consumers
		ioc.InitConsumers,
		events.NewMySQLBinlogConsumer,
		service.NewSyncService,
		repository.NewCourseRepository,
		dao.NewCourseElasticDAO,
		ioc.InitESClient,
		// grpc
		ioc.InitGRPCxKratosServer,
		grpc.NewSearchServiceServer,
		service.NewSearchService,
		repository.NewHistoryRepository,
		dao.NewSearchHistoryGORMDAO,
		// 第三方
		ioc.InitKafka,
		ioc.InitDB,
		ioc.InitEtcdClient,
		ioc.InitLogger,
	)
	return &App{}
}
