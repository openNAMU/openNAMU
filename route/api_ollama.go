package route

import (
	"bytes"
	"context"
	stdjson "encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"opennamu/route/tool"
)

const ollama_default_url = "http://127.0.0.1:11434"
const ollama_timeout = 10 * time.Minute
const ollama_idle_timeout = 90 * time.Second

type ollama_request_data struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollama_response_data struct {
	Response string `json:"response"`
	Error    string `json:"error"`
	Done     bool   `json:"done"`
}

func Api_ollama(model string, prompt string) (string, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		model = "gemma4:e4b"
	}

	request_data, err := stdjson.Marshal(ollama_request_data{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		ollama_default_url+"/api/generate",
		bytes.NewReader(request_data),
	)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: ollama_timeout}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama status: %s", response.Status)
	}

	var response_data ollama_response_data
	if err := stdjson.NewDecoder(io.LimitReader(response.Body, 4<<20)).Decode(&response_data); err != nil {
		return "", err
	}
	if response_data.Error != "" {
		return "", fmt.Errorf("ollama: %s", response_data.Error)
	}

	return response_data.Response, nil
}

type ollama_stream_data struct {
	response ollama_response_data
	err      error
}

func Api_ollama_stream(model string, prompt string) (string, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		model = "gemma4:e4b"
	}

	request_data, err := stdjson.Marshal(ollama_request_data{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	})
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), ollama_timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ollama_default_url+"/api/generate",
		bytes.NewReader(request_data),
	)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama status: %s", response.Status)
	}

	chunk_channel := make(chan ollama_stream_data, 1)
	go func() {
		decoder := stdjson.NewDecoder(io.LimitReader(response.Body, 4<<20))
		for {
			response_data := ollama_response_data{}
			if decode_err := decoder.Decode(&response_data); decode_err != nil {
				chunk_channel <- ollama_stream_data{err: decode_err}
				return
			}
			chunk_channel <- ollama_stream_data{response: response_data}
			if response_data.Done {
				return
			}
		}
	}()

	idle_timer := time.NewTimer(ollama_idle_timeout)
	defer idle_timer.Stop()
	reset_idle_timer := func() {
		if !idle_timer.Stop() {
			select {
			case <-idle_timer.C:
			default:
			}
		}
		idle_timer.Reset(ollama_idle_timeout)
	}

	answer := strings.Builder{}
	for {
		select {
		case chunk := <-chunk_channel:
			if chunk.err != nil {
				return "", chunk.err
			}
			reset_idle_timer()
			if chunk.response.Error != "" {
				return "", fmt.Errorf("ollama: %s", chunk.response.Error)
			}
			answer.WriteString(chunk.response.Response)
			if chunk.response.Done {
				return answer.String(), nil
			}
		case <-idle_timer.C:
			cancel()
			return "", fmt.Errorf("ollama idle timeout")
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func Api_ollama_stream_post(config tool.Config, question string, model string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "ai_use", config.IP) {
		return map[string]any{"response": "require auth"}
	}

	question = strings.TrimSpace(question)
	if question == "" {
		return map[string]any{"response": "error", "data": "question"}
	}
	if tool.Get_len(question) > 1000 {
		question = tool.Get_slice(question, 0, 1000)
	}

	model = strings.TrimSpace(model)
	if model == "" {
		model = "gemma4:e4b"
	}

	context_data, source_list := ollama_document_context(db, config, question)
	prompt := "너는 위키 문서 검색을 돕는 AI다. 아래 참고 문서에 있는 내용만 근거로 답변하고, 근거가 없으면 모른다고 답변해라. 참고 문서: " + context_data + " 질문: " + question
	answer, err := Api_ollama_stream(model, prompt)
	if err != nil {
		return map[string]any{"response": "error", "data": "ollama", "source": source_list}
	}

	return map[string]any{"response": "ok", "data": answer, "source": source_list}
}
