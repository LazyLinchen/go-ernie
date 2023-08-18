package go_ernie

const (
	ErnieBot      = "erniebot"
	ErnieBotTurBo = "erniebot-turbo"
	BLOOMZ7B      = "BLOOMZ-7B"
)

var enableChatCompletionModels = map[string]string{
	ErnieBot:      "/wenxinworkshop/chat/completions",
	ErnieBotTurBo: "/wenxinworkshop/chat/eb-instant",
	BLOOMZ7B:      "/wenxinworkshop/chat/bloomz_7b1",
}

func isSupportedChatCompletionModel(model string) bool {
	_, ok := enableChatCompletionModels[model]
	return ok
}

func chatCompletionUri(model string) string {
	return enableChatCompletionModels[model]
}
