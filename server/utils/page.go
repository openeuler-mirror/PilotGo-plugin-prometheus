package utils

import (
	"fmt"
	"reflect"

	"gitee.com/openeuler/PilotGo/sdk/response"
)

// 结构体分页查询方法
func DataPaging(p *response.PaginationQ, list interface{}, total int) (interface{}, error) {
	data := make([]interface{}, 0)
	if reflect.TypeOf(list).Kind() == reflect.Slice {
		s := reflect.ValueOf(list)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			data = append(data, ele.Interface())
		}
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if total == 0 {
		p.TotalSize = 0
	}
	num := p.PageSize * (p.Page - 1)
	if num > total {
		return nil, fmt.Errorf("页码超出")
	}
	if p.PageSize*p.Page > total {
		return data[num:], nil
	} else {
		if p.PageSize*p.Page < num {
			return nil, fmt.Errorf("读取错误")
		}
		if p.PageSize*p.Page == 0 {
			return data, nil
		} else {
			return data[num : p.Page*p.PageSize], nil
		}
	}
}
