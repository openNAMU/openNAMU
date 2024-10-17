package route

import (
    "context"
    "database/sql"
    "log"
    "opennamu/route/tool"

    "github.com/google/generative-ai-go/genai"
    jsoniter "github.com/json-iterator/go"
    "google.golang.org/api/option"
)

func Api_func_llm(call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    db := tool.DB_connect()
    defer db.Close()

    var api_key string

    stmt, err := db.Prepare(tool.DB_change("select data from user_set where name = 'llm_api_key' and id = ?"))
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    err = stmt.QueryRow(other_set["ip"]).Scan(api_key)
    if err != nil {
        if err == sql.ErrNoRows {
            api_key = ""
        } else {
            log.Fatal(err)
        }
    }

    ctx := context.Background()

    client, err := genai.NewClient(ctx, option.WithAPIKey(api_key))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    model := client.GenerativeModel("gemini-pro")
    resp, err := model.GenerateContent(ctx, genai.Text(other_set["prompt"]))
    if err != nil {
        log.Fatal(err)
    }

    text := resp.Candidates[0].Content.Parts[0]

    json_data, _ := json.Marshal(map[string]genai.Part{"data": text})
    return string(json_data)
}
