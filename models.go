package go_ernie

var ChatModelEndpoint = map[string]string{
	"ERNIE-4.0-Turbo-8K":           "/chat/ernie-4.0-turbo-8k",
	"ERNIE-4.0-8K-Latest":          "/chat/ernie-4.0-8k-latest",
	"ERNIE-4.0-8K-0613":            "/chat/ernie-4.0-8k-0613",
	"ERNIE-3.5-8K-0613":            "/chat/ernie-3.5-8k-0613",
	"ERNIE-Bot-turbo":              "/chat/eb-instant",
	"ERNIE-Lite-8K-0922":           "/chat/eb-instant",
	"ERNIE-Lite-8K":                "/chat/ernie-lite-8k",
	"ERNIE-Lite-8K-0308":           "/chat/ernie-lite-8k",
	"ERNIE-3.5-8K":                 "/chat/completions",
	"ERNIE-Bot":                    "/chat/completions",
	"ERNIE-4.0-8K":                 "/chat/completions_pro",
	"ERNIE-4.0-8K-Preview":         "/chat/ernie-4.0-8k-preview",
	"ERNIE-4.0-8K-Preview-0518":    "/chat/completions_adv_pro",
	"ERNIE-4.0-8K-0329":            "/chat/ernie-4.0-8k-0329",
	"ERNIE-4.0-8K-0104":            "/chat/ernie-4.0-8k-0104",
	"ERNIE-Bot-4":                  "/chat/completions_pro",
	"ERNIE-Bot-8k":                 "/chat/ernie_bot_8k",
	"ERNIE-3.5-128K":               "/chat/ernie-3.5-128k",
	"ERNIE-3.5-8K-preview":         "/chat/ernie-3.5-8k-preview",
	"ERNIE-3.5-8K-0329":            "/chat/ernie-3.5-8k-0329",
	"ERNIE-3.5-4K-0205":            "/chat/ernie-3.5-4k-0205",
	"ERNIE-3.5-8K-0205":            "/chat/ernie-3.5-8k-0205",
	"ERNIE-3.5-8K-1222":            "/chat/ernie-3.5-8k-1222",
	"ERNIE Speed":                  "/chat/ernie_speed",
	"ERNIE-Speed":                  "/chat/ernie_speed",
	"ERNIE-Speed-8K":               "/chat/ernie_speed",
	"ERNIE-Speed-128K":             "/chat/ernie-speed-128k",
	"ERNIE Speed-AppBuilder":       "/chat/ai_apaas",
	"ERNIE-Tiny-8K":                "/chat/ernie-tiny-8k",
	"ERNIE-Function-8K":            "/chat/ernie-func-8k",
	"ERNIE-Character-8K":           "/chat/ernie-char-8k",
	"ERNIE-Character-Fiction-8K":   "/chat/ernie-char-fiction-8k",
	"ERNIE-Bot-turbo-AI":           "/chat/ai_apaas",
	"EB-turbo-AppBuilder":          "/chat/ai_apaas",
	"BLOOMZ-7B":                    "/chat/bloomz_7b1",
	"Llama-2-7b-chat":              "/chat/llama_2_7b",
	"Llama-2-13b-chat":             "/chat/llama_2_13b",
	"Llama-2-70b-chat":             "/chat/llama_2_70b",
	"Qianfan-Chinese-Llama-2-7B":   "/chat/qianfan_chinese_llama_2_7b",
	"Qianfan-Chinese-Llama-2-13B":  "/chat/qianfan_chinese_llama_2_13b",
	"Qianfan-Chinese-Llama-2-70B":  "/chat/qianfan_chinese_llama_2_70b",
	"Meta-Llama-3-8B":              "/chat/llama_3_8b",
	"Meta-Llama-3-70B":             "/chat/llama_3_70b",
	"Qianfan-BLOOMZ-7B-compressed": "/chat/qianfan_bloomz_7b_compressed",
	"ChatGLM2-6B-32K":              "/chat/chatglm2_6b_32k",
	"AquilaChat-7B":                "/chat/aquilachat_7b",
	"XuanYuan-70B-Chat-4bit":       "/chat/xuanyuan_70b_chat",
	"ChatLaw":                      "/chat/chatlaw",
	"Yi-34B-Chat":                  "/chat/yi_34b_chat",
	"Mixtral-8x7B-Instruct":        "/chat/mixtral_8x7b_instruct",
	"Gemma-7B-it":                  "/chat/gemma_7b_it",
}

func isSupportedChatCompletionModel(model string) bool {
	_, ok := ChatModelEndpoint[model]
	return ok
}

func chatCompletionUri(model string) string {
	return ChatModelEndpoint[model]
}

var CompletionModelEndpoint = map[string]string{
	"SQLCoder-7B":           "/completions/sqlcoder_7b",
	"CodeLlama-7b-Instruct": "/completions/codellama_7b_instruct",
}

func isSupportedCompletionModel(model string) bool {
	_, ok := CompletionModelEndpoint[model]
	return ok
}

func completionUri(model string) string {
	return CompletionModelEndpoint[model]
}
