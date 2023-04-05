package paging

import proto_paging "github.com/flamefatex/grpc-gateway-example/proto/gen/go/common/paging"

const DefaultLimit = 10

// Normalize 标准化
func Normalize(p *proto_paging.Paging) *proto_paging.Paging {
	if p == nil {
		p = &proto_paging.Paging{
			Page:     1,
			PageSize: 10,
		}
	}
	return p
}

// OffsetLimit 获取sql的offset limit
func OffsetLimit(p *proto_paging.Paging) (offset, limit int) {
	// 初始化默认返回数量
	limit = DefaultLimit

	if p != nil {
		offset = int((p.Page - 1) * p.PageSize)
		limit = int(p.PageSize)
	}
	return
}

// WithTotal 填充总数
func WithTotal(p *proto_paging.Paging, total int64) *proto_paging.Paging {
	if p != nil && p.PageSize != 0 {
		p.TotalCount = int32(total)
		p.TotalPage = p.TotalCount / p.PageSize
	}

	return p
}
