package domain

import searchv1 "github.com/MuxiKeStack/be-api/gen/proto/search/v1"

type SearchHistory struct {
	Id             int64
	Uid            int64
	SearchLocation searchv1.SearchLocation
	Keyword        string
	Status         searchv1.VisibilityStatus
}
