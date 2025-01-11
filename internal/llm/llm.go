package llm

import "context"

type LLM interface {
	Translate(ctx context.Context, sourceLanguage string, destLanguage, contents []string) ([][]string, error)
}

func NewLLM() LLM {
	return &llmClient{}
}

type llmClient struct {
}

func (l *llmClient) Translate(ctx context.Context, language string, destLanguage, contents []string) ([][]string, error) {
	// TODO(kun.li): use really llm
	var res [][]string
	res = append(res, append([]string{language}, contents...))
	for _, destLang := range destLanguage {
		_ = destLang
		var destContents []string
		destContents = append(destContents, destLang)
		for _, c := range contents {
			destContents = append(destContents, c)
		}
		res = append(res, destContents)
	}
	return res, nil
}
