package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type OCRService struct {
	aiService *AIService
}

func NewOCRService(ai *AIService) *OCRService {
	return &OCRService{
		aiService: ai,
	}
}

type ImageExtractResult struct {
	Title        string       `json:"title"`
	BaseServings int          `json:"base_servings"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
}

func (s *OCRService) ExtractFromImage(ctx context.Context, fileBytes []byte, contentType string) (*ImageExtractResult, error) {
	prompt := fmt.Sprintf(`이 요리 레시피 이미지에서 제목, 기준 인분 수, 재료 목록, 조리 순서를 추출하세요.
결과는 반드시 아래 JSON 형식으로만 응답하세요. 마크다운 코드블록은 포함하지 마세요.

형식:
{
  "title": "레시피 제목",
  "base_servings": 2,
  "ingredients": [
    {
      "name": "재료명", 
      "amount_num": 100, 
      "unit": "g", 
      "note": "용량 표현",
      "scaling_type": "linear",
      "scaling_factor": 1.0
    }
  ],
  "instructions": ["1단계 설명", "2단계 설명", ...],
  "error": ""
}

%s

재료나 조리법을 찾을 수 없으면 빈 배열 []을 사용하세요.`, CommonPromptInstructions)

	text, err := s.aiService.GenerateFromImage(ctx, prompt, fileBytes, contentType)
	if err != nil {
		return nil, fmt.Errorf("Gemini 분석 실패: %w", err)
	}

	var result struct {
		ImageExtractResult
		Error string `json:"error"`
	}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("데이터 파싱 실패: %w (응답: %s)", err, text)
	}

	if result.Error != "" {
		return nil, errors.New(result.Error)
	}

	return &result.ImageExtractResult, nil
}
