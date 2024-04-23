package grpc

import (
	"context"
	searchv1 "github.com/MuxiKeStack/be-api/gen/proto/search/v1"
	"github.com/MuxiKeStack/be-search/domain"
	"github.com/MuxiKeStack/be-search/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type SearchServiceServer struct {
	svc service.SearchService
	searchv1.UnimplementedSearchServiceServer
}

func NewSearchServiceServer(svc service.SearchService) *SearchServiceServer {
	return &SearchServiceServer{svc: svc}
}

func (s *SearchServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSearchServiceServer(server, s)
}

func (s *SearchServiceServer) SearchCourse(ctx context.Context, request *searchv1.SearchCourseRequest) (*searchv1.SearchCourseResponse, error) {
	courses, err := s.svc.SearchCourse(ctx, request.GetKeyword(), request.GetUid(), request.GetLocation())
	return &searchv1.SearchCourseResponse{
		Courses: slice.Map(courses, func(idx int, src domain.Course) *searchv1.Course {
			return &searchv1.Course{
				Id:             src.Id,
				Name:           src.Name,
				Teacher:        src.Teacher,
				CompositeScore: src.CompositeScore,
			}
		}),
	}, err
}

func (s *SearchServiceServer) GetUserSearchHistories(ctx context.Context, request *searchv1.GetUserHistoryRequest) (*searchv1.GetUserHistoryResponse, error) {
	hs, err := s.svc.GetUserSearchHistories(ctx, request.GetUid(), request.GetLocation(),
		request.GetOffset(), request.GetLimit())
	return &searchv1.GetUserHistoryResponse{
		Histories: slice.Map(hs, func(idx int, src domain.SearchHistory) *searchv1.SearchHistory {
			return &searchv1.SearchHistory{
				Id:      src.Id,
				Keyword: src.Keyword,
				Status:  src.Status,
			}
		}),
	}, err
}

func (s *SearchServiceServer) HideUserSearchHistories(ctx context.Context, request *searchv1.HideUserSearchHistoriesRequest) (*searchv1.HideUserSearchHistoriesResponse, error) {
	err := s.svc.HideUserSearchHistories(ctx, request.GetUid(), request.GetLocation(),
		request.GetRemoveAll(), request.GetHistoryIds())
	return &searchv1.HideUserSearchHistoriesResponse{}, err
}
