package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Ingredient struct {
	Name          string  `json:"name"`
	AmountNum     float64 `json:"amount_num"`     // 숫자 (계산용)
	Unit          string  `json:"unit"`           // 표준 단위 (ml, g, 개 등)
	Note          string  `json:"note"`           // 원본 표현이나 조리 상태 (예: "1스푼", "다져서", "약간")
	ScalingType   string  `json:"scaling_type"`   // linear (정비례) or culinary (조리식)
	ScalingFactor float64 `json:"scaling_factor"` // 추가 인분당 증가 비율 (linear는 1.0, culinary는 0.5~0.9)
}

type AIService struct{}

func NewAIService() *AIService {
	return &AIService{}
}

const CommonPromptInstructions = `[필수 지시사항]
1. 요리 관련 확인: 만약 입력된 정보가 요리 레시피와 관련이 없거나, 재료 정보를 전혀 찾을 수 없는 경우 응답 JSON에 "error": "레시피 정보를 찾을 수 없습니다" 필드를 포함하고 다른 필드는 비우거나 기본값을 사용하세요.
2. 원본 인분 수 파악: 영상이나 이미지에서 이 레시피가 몇 인분 기준(예: 1인분, 2인분 등)인지 파악하여 응답의 'base_servings' 필드에 정수로 넣으세요. 명시되지 않은 경우 1을 기본값으로 사용하세요.
3. 재료 추출: 재료의 용량은 영상이나 이미지에 나온 '원본 그대로' 추출하세요. 1인분으로 환산하지 마세요.
4. 스케일링 규칙 설정 (향후 계산용):
   - 각 재료마다 'scaling_type'과 'scaling_factor'를 설정하세요.
   - 고기, 채소, 생선 등 주재료: scaling_type="linear", scaling_factor=1.0
   - 소금, 간장, 설탕, 고추장 등 양념류와 물, 육수 등 액체류: scaling_type="culinary", scaling_factor=0.5~0.8 (요리 상식에 따라 추가 인분 시 적게 늘어나는 비율)
5. 표준 단위 환산:
   - 'amount_num'에는 원본 용량을 'ml' 또는 'g'으로 환산한 숫자를 넣으세요.
   - 환산 규칙: 1큰술(1T)=15ml, 1작은술(1t)=5ml, 종이컵 1컵=180ml, 일반컵(1C)=200ml.
   - '개', '줄', '봉지' 등 셀 수 있는 단위는 'unit'에 넣고 숫자는 'amount_num'에 넣으세요.
6. 원본 표현 보존: 'note' 필드에는 원본의 표현(예: "2스푼", "반 컵")을 그대로 적으세요.
재료나 조리법을 찾을 수 없으면 빈 배열 []을 사용하세요.
`

func (s *AIService) getClient(ctx context.Context) (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY가 설정되지 않았습니다")
	}
	return genai.NewClient(ctx, option.WithAPIKey(apiKey))
}

func (s *AIService) GenerateFromText(ctx context.Context, prompt string) (string, error) {
	client, err := s.getClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.ResponseMIMEType = "application/json"

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini 응답이 비어있습니다")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return "", fmt.Errorf("Gemini 응답이 텍스트 형식이 아닙니다")
	}

	return string(text), nil
}

func (s *AIService) GenerateFromImage(ctx context.Context, prompt string, fileBytes []byte, contentType string) (string, error) {
	client, err := s.getClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.ImageData(strings.Split(contentType, "/")[1], fileBytes),
	)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini 응답을 받을 수 없습니다")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return "", fmt.Errorf("Gemini 응답 형식이 텍스트가 아닙니다")
	}

	return string(text), nil
}
