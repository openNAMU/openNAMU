package route

import (
	"context"
	"opennamu/route/tool"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func Api_func_llm(config tool.Config, prompt string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_key := ""
	tool.QueryRow_DB(
		db,
		"select data from user_set where name = 'llm_api_key' and id = ?",
		[]any{&api_key},
		config.IP,
	)

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(api_key))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		panic(err)
	}

	text := resp.Candidates[0].Content.Parts[0]

	return_data := make(map[string]any)
	return_data["data"] = text
	return return_data
}
