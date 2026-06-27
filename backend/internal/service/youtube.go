package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type YouTubeService struct {
	aiService *AIService
}

func NewYouTubeService(ai *AIService) *YouTubeService {
	return &YouTubeService{
		aiService: ai,
	}
}

type YouTubeExtractResponse struct {
	Title        string       `json:"title"`
	BaseServings int          `json:"base_servings"`
	Ingredients  []Ingredient `json:"ingredients"`
	ImageURL     string       `json:"image_url,omitempty"`
}

func (s *YouTubeService) ExtractRecipeData(ctx context.Context, rawURL string) (*YouTubeExtractResponse, error) {
	videoID, err := s.extractVideoID(rawURL)
	if err != nil {
		return nil, err
	}

	title, description, err := s.fetchVideoInfo(ctx, videoID)
	if err != nil {
		return nil, err
	}

	ingredients, baseServings, err := s.extractIngredients(ctx, title, description)
	if err != nil {
		return nil, err
	}

	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", videoID)

	return &YouTubeExtractResponse{
		Title:        title,
		BaseServings: baseServings,
		Ingredients:  ingredients,
		ImageURL:     thumbnailURL,
	}, nil
}

func (s *YouTubeService) extractVideoID(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("유효하지 않은 URL")
	}
	switch u.Host {
	case "youtu.be":
		id := strings.TrimPrefix(u.Path, "/")
		if idx := strings.Index(id, "?"); idx != -1 {
			id = id[:idx]
		}
		if id == "" {
			return "", fmt.Errorf("비디오 ID를 찾을 수 없습니다")
		}
		return id, nil
	case "youtube.com", "www.youtube.com", "m.youtube.com":
		if v := u.Query().Get("v"); v != "" {
			return v, nil
		}
		parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
		if len(parts) >= 2 && (parts[0] == "shorts" || parts[0] == "live") {
			return parts[1], nil
		}
	}
	return "", fmt.Errorf("YouTube URL을 인식할 수 없습니다")
}

type ytSnippet struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ytItem struct {
	Snippet ytSnippet `json:"snippet"`
}

type ytAPIResponse struct {
	Items []ytItem `json:"items"`
}

func (s *YouTubeService) fetchVideoInfo(ctx context.Context, videoID string) (title, description string, err error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return "", "", fmt.Errorf("YOUTUBE_API_KEY가 설정되지 않았습니다")
	}
	endpoint := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?id=%s&part=snippet&key=%s",
		url.QueryEscape(videoID), apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var yt ytAPIResponse
	if err = json.NewDecoder(resp.Body).Decode(&yt); err != nil {
		return "", "", err
	}
	if len(yt.Items) == 0 {
		return "", "", fmt.Errorf("비디오를 찾을 수 없습니다")
	}
	return yt.Items[0].Snippet.Title, yt.Items[0].Snippet.Description, nil
}

func (s *YouTubeService) extractIngredients(ctx context.Context, title, description string) ([]Ingredient, int, error) {
	prompt := fmt.Sprintf(`아래 유튜브 요리 영상 정보에서 재료 목록과 기준 인분 수를 추출하세요.

%s

[응답 형식]
{
  "base_servings": 2,
  "ingredients": [
    {
      "name": "진간장",
      "amount_num": 30,
      "unit": "ml",
      "note": "2스푼",
      "scaling_type": "culinary",
      "scaling_factor": 0.6
    }
  ],
  "error": ""
}

제목: %s
설명: %s`, CommonPromptInstructions, title, description)

	text, err := s.aiService.GenerateFromText(ctx, prompt)
	if err != nil {
		return nil, 0, fmt.Errorf("Gemini API 호출 실패: %w", err)
	}

	var result struct {
		BaseServings int          `json:"base_servings"`
		Ingredients  []Ingredient `json:"ingredients"`
		Error        string       `json:"error"`
	}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, 0, fmt.Errorf("재료 JSON 파싱 실패: %w (응답: %s)", err, text)
	}

	if result.Error != "" {
		return nil, 0, fmt.Errorf(result.Error)
	}

	return result.Ingredients, result.BaseServings, nil
}
