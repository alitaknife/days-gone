package utils

import (
	"errors"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"strconv"
)

var rule = []string{
	"id@required|integer#id不能为空值|请输入一个整数",
}

func ValidId(r *ghttp.Request) (int, error) {
	strId, ok := r.Get("id").(string)
	if ok {
		err := gvalid.CheckMap(r.Context(), g.Map{"id": strId}, rule)
		if err != nil {
			if v, ok := err.(gvalid.Error); ok {
				return -1, errors.New(gerror.Current(v).Error())
			}
			return -1, errors.New("please enter legal parameters")
		}
		idInt, _ := strconv.Atoi(strId)
		return idInt, err
	} else {
		return -1, errors.New("please enter legal parameters")
	}
}
