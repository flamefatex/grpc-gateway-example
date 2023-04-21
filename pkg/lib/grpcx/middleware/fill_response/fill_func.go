package fill_response

import "context"

// FillResponseFunc 注入ctx方法
type FillResponseFunc func(ctx context.Context, resp interface{}, err error) (newResp interface{})
