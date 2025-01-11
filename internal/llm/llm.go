package llm

import (
	"context"
	"math/rand/v2"

	"cloud.google.com/go/translate"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/text/language"
)

type LLM interface {
	Translate(ctx context.Context, sourceLanguage string, destLanguage, contents []string) ([][]string, error)
}

func NewLLM() LLM {
	client, err := translate.NewClient(context.Background())
	if err != nil {
		logc.Errorf(context.Background(), "llm new client err %v", err)
		//	return nil
	}

	return &llmClient{client: client}
}

type llmClient struct {
	client *translate.Client
}

func (l *llmClient) Translate(ctx context.Context, sourceLanguage string, destLanguage, contents []string) ([][]string, error) {
	// TODO(kun.li): use really llm
	var res [][]string
	res = append(res, append([]string{sourceLanguage}, contents...))
	if rand.IntN(3) == 0 {
		// mock
		logc.Infof(ctx, "mock llm work")
		for _, destLang := range destLanguage {
			_ = destLang
			var destContents []string
			destContents = append(destContents, destLang)
			for _, c := range contents {
				destContents = append(destContents, c)
			}
			res = append(res, destContents)
		}
	} else {
		// really but not work
		for _, destLang := range destLanguage {
			logc.Infof(ctx, "really llm not work")
			parseLang, err := language.Parse(destLang)
			if err != nil {
				logc.Errorf(ctx, "llm language parse %v err %v", destLang, err)
				return nil, err
			}

			resp, err := l.client.Translate(ctx, contents, parseLang, nil)
			if err != nil {
				logc.Errorf(ctx, "llm translate err %v", err)
				return nil, err
			}
			// api 第一条是被翻译语言,所以结果会多一条
			if len(resp) != len(contents)+1 {
				logc.Errorf(ctx, "llm translate result count err %v", err)
				return nil, err
			}

			var destContents []string
			for _, t := range resp[1:] {
				destContents = append(destContents, t.Text)
			}
			res = append(res, append([]string{destLang}, destContents...))
		}
	}

	return res, nil
}
