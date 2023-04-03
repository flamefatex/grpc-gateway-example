package fill_response

import (
	"context"
	"reflect"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/statusx"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/tracing/opentracing"
	proto_status "github.com/flamefatex/grpc-gateway-example/proto/gen/go/common/status"
)

// DefaultFillFunc 默认填充
func DefaultFillFunc(ctx context.Context, resp interface{}, err error) (newResp interface{}) {
	newResp = resp
	// 若响应正常，则填充request_id 和 status
	if err == nil {
		if resp == nil { // 接口动态类型 & 接口动态值都为空
			return
		}

		// 获取resp实例地址的反射值对象
		valueOfResp := reflect.ValueOf(resp)
		if valueOfResp.Kind() != reflect.Ptr {
			return
		}
		if valueOfResp.IsNil() { // 接口动态值为空，需要重新初始化一个，否则会panic报错
			t := reflect.TypeOf(resp)
			if t.Kind() == reflect.Ptr { //指针类型获取真正type需要调用Elem
				t = t.Elem()
			}
			newResp = reflect.New(t).Interface() // 调用反射创建对象
			valueOfResp = reflect.ValueOf(newResp)
		}

		// 取出resp实例地址的元素
		valueOfResp = valueOfResp.Elem()
		// 获取RequestId字段的值
		vRequestId := valueOfResp.FieldByName("RequestId")
		// 尝试设置RequestId的值
		if vRequestId.IsValid() && vRequestId.Kind() == reflect.String {
			vRequestId.SetString(opentracing.GetTraceIdFromCtx(ctx))
		}
		vStatus := valueOfResp.FieldByName("Status")
		if vStatus.IsValid() {
			// 验证类型是否符合预期
			_, ok := vStatus.Interface().(*proto_status.ResponseStatus)
			if ok {
				vStatus.Set(reflect.ValueOf(statusx.OK(ctx)))
			}
		}
	}
	return
}
