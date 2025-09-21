package minimal

import (
	"regexp"
	"time"

	"github.com/dreamerjackson/crawler/limiter"
	"github.com/dreamerjackson/crawler/spider"
	"golang.org/x/time/rate"
)

// ExampleTask 是一个最小可运行的任务示例：
// - 访问 https://baidu.com/
// - 解析 <title> 文本并输出到存储。
var ExampleTask = &spider.Task{
	Options: spider.Options{
		Name:     "example_baidu_home",
		Reload:   true,
		WaitTime: 1,
		MaxDepth: 1,
		// 内置限速器，保证 Request.Fetch() 不会因为 nil 限速器而 panic。
		Limit: limiter.Multi(
			rate.NewLimiter(limiter.Per(1, 1*time.Second), 1),
		),
	},
	Rule: spider.RuleTree{
		Root: func() ([]*spider.Request, error) {
			return []*spider.Request{{
				Priority: 1,
				URL:      "https://baidu.com/",
				Method:   "GET",
				RuleName: "parse_title",
			}}, nil
		},
		Trunk: map[string]*spider.Rule{
			"parse_title": {
				ItemFields: []string{"title"},
				ParseFunc: func(ctx *spider.Context) (spider.ParseResult, error) {
					re := regexp.MustCompile(`(?i)<title>([^<]+)</title>`)
					m := re.FindSubmatch(ctx.Body)
					title := ""
					if len(m) >= 2 {
						title = string(m[1])
					}
					data := ctx.Output(map[string]interface{}{"title": title})
					return spider.ParseResult{Items: []interface{}{data}}, nil
				},
			},
		},
	},
}