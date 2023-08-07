package presenter

type T any

type ApiResponse struct {
	Data     T        `json:"data" `
	Messages []string `json:"messages" `
}

func GetSuccessResponse(messages ...string) ApiResponse {
	data := ApiResponse{
		Messages: messages,
	}

	return data
}
