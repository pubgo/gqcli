package gq

import (
	"github.com/storyicon/graphquery"
	"time"
)

// Parse 解析文件
// doc 需要解析的文档或者内容
// expr gq表达式
func Parse(doc, expr string) ParseResp {
	timeStart := time.Now().UnixNano()
	resp := graphquery.ParseFromString(doc, expr)
	timeCost := time.Now().UnixNano() - timeStart
	if len(resp.Errors) > 0 {
		return ParseResp{
			Error:    resp.Errors[0],
			TimeCost: timeCost,
		}
	}

	return ParseResp{
		TimeCost: timeCost,
		Data:     resp.Data,
	}
}
