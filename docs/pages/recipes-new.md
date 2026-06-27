# `/recipes/new` — 새 레시피 등록

## 역할
세 가지 방식으로 레시피를 등록하는 페이지.
- **YouTube**: 영상 URL 입력 → AI가 설명글에서 재료 자동 추출
- **Instagram**: URL 입력 (추출 버튼 없음, 저장만 가능 — 추출 미구현)
- **직접 입력**: 사진 업로드 → Gemini OCR로 레시피 추출

## 사용자 액션

| 액션 | 트리거 | API 호출 | 성공 | 실패 |
|------|--------|----------|------|------|
| YouTube 탭 선택 | `click` | 없음 | 폼 초기화, URL 인풋 노출 | — |
| Instagram 탭 선택 | `click` | 없음 | 폼 초기화, URL 인풋 노출 | — |
| 직접 입력 탭 선택 | `click` | 없음 | 폼 초기화, 사진 업로드 영역 노출 | — |
| 재료 가져오기 버튼 (YouTube) | `click` | `POST /api/youtube/extract` | 재료 목록 채워짐, 토스트: "N인분 기준 레시피 정보를 가져왔습니다!" | 토스트 에러 |
| 사진 업로드 (직접 입력) | `change` | `POST /api/recipes/extract-image` | 재료+조리법 채워짐, 토스트: "사진에서 레시피 정보를 추출했습니다!" | 토스트 에러 |
| 인분 + 버튼 | `click` | 없음 | 인분 수 +1, 재료 용량 재계산 | — |
| 인분 - 버튼 | `click` | 없음 | 인분 수 -1 (최소 1) | — |
| + 재료 추가 | `click` | 없음 | 빈 재료 행 추가 | — |
| 재료 ✕ 버튼 | `click` | 없음 | 해당 행 삭제 | — |
| + 단계 추가 (직접 입력) | `click` | 없음 | 빈 조리법 단계 추가 | — |
| 단계 ✕ 버튼 (직접 입력) | `click` | 없음 | 해당 단계 삭제 (마지막 단계면 내용만 초기화) | — |
| 레시피 저장하기 (폼 submit) | `submit` | `POST /api/recipes` | 토스트: "레시피가 성공적으로 등록되었습니다!", `/recipes` 로 이동 | 토스트 에러 |
| 레시피 목록으로 돌아가기 | `click` (a 태그) | 없음 | `/recipes` 로 이동 | — |

## 에러 케이스

| 상황 | 에러 메시지 | 표시 방식 |
|------|------------|-----------|
| YouTube 추출 실패 (네트워크·서버) | 서버 응답 텍스트 또는 `"재료 추출 중 오류가 발생했습니다"` | 토스트 (error) |
| YouTube 영상 설명글에 재료 없음 | `"영상 설명글에서 재료를 찾지 못했습니다."` | 토스트 (error) |
| 사진 분석 실패 (네트워크·서버) | 서버 응답 텍스트 또는 `"사진 분석 중 오류가 발생했습니다"` | 토스트 (error) |
| 이미지 압축 실패 | `"압축 실패"` | 토스트 (error) |
| 레시피 저장 실패 (네트워크·서버) | 서버 응답 텍스트 또는 `"등록 중 오류가 발생했습니다. 다시 시도해주세요."` | 토스트 (error) |

## 백엔드 에러 (서버측)

| API | 상황 | HTTP | 서버 메시지 |
|-----|------|------|------------|
| `POST /api/youtube/extract` | URL 없음 또는 파싱 실패 | 400 | `"URL이 필요합니다"` |
| `POST /api/youtube/extract` | AI 추출 실패 | 500 | 에러 문자열 |
| `POST /api/recipes/extract-image` | 멀티파트 파싱 실패 (10MB 초과 포함) | 400 | `"파일 파싱 실패"` |
| `POST /api/recipes/extract-image` | image 필드 없음 | 400 | `"이미지 파일이 필요합니다"` |
| `POST /api/recipes/extract-image` | 파일 읽기 실패 | 500 | `"파일 읽기 실패"` |
| `POST /api/recipes/extract-image` | AI 분석 실패 | 500 | 에러 문자열 |
| `POST /api/recipes` | 요청 바디 파싱 실패 | 400 | `"잘못된 요청 본문"` |
| `POST /api/recipes` | DB 저장 실패 | 500 | 에러 문자열 |

## 상태

| 상태 | 초기값 | 설명 |
|------|--------|------|
| `platform` | `'youtube'` | 선택된 탭 |
| `url` | `''` | 링크 인풋 값 |
| `title` | `''` | 레시피 이름 |
| `ingredients` | `[]` | 현재 인분에 맞게 계산된 재료 목록 |
| `originalIngredients` | `[]` | AI 추출 원본 재료 (인분 계산 기준) |
| `baseServings` | `1` | AI 추출 기준 인분 수 |
| `currentServings` | `1` | 현재 선택 인분 수 |
| `isExtracting` | `false` | 추출 중 여부 (버튼/업로드 비활성화) |
| `isSubmitting` | `false` | 저장 중 여부 (저장 버튼 비활성화) |

## 저장 버튼 비활성화 조건
`isSubmitting === true` OR `title === ''`

## 이미지 전처리
업로드된 이미지는 API 전송 전 클라이언트에서 최대 1200px로 리사이즈 후 JPEG(quality 0.7)로 압축.

## 연관 API
- `POST /api/youtube/extract` → `YouTubeHandler.Extract`
- `POST /api/recipes/extract-image` → `OCRHandler.ExtractImage`
- `POST /api/recipes` → `RecipeHandler.Create`

## 연관 컴포넌트
- `$lib/toast.svelte` — `toast.add(message, type)`
- `$lib/api/recipes.ts` — `extractFromYouTube()`, `extractFromImage()`, `createRecipe()`

## 미구현 항목
- Instagram 탭 → URL 저장 시 실제 추출 기능 없음
