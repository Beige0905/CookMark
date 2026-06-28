# 002 — CookingLog로 RecipeNote 대체

날짜: 2026-06-27
상태: 결정됨 (ADR 001 일부 무효화)

## 배경

ADR 001에서 설계한 RecipeNote(개인 메모 + 재료 퍼센트 조정)와 CookingLog(날짜별 요리 일지)가 역할이 겹친다는 판단.

## 결정

RecipeNote를 제거하고 CookingLog로 통합한다. 재료 퍼센트 조정 기능도 함께 제거한다.

**CookingLog 스펙:**
- 한 레시피에 여러 개 누적 (날짜별 기록)
- 필드: `recipe_id`, `user_id`, `comment` (자유 텍스트), `cooked_at`

## 이유

- RecipeNote의 "조정 기능"은 UX 복잡도 대비 사용 빈도가 낮을 것으로 판단
- CookingLog의 날짜별 감상 메모가 실제 사용 패턴에 더 적합
- 두 개를 유지하면 사용자가 어디에 뭘 써야 하는지 혼란

## 영향 범위

**제거:**
- `recipe_personal_notes` DB 테이블
- `RecipeNote` 모델, `NoteService`, `NoteHandler`, `pg_note.go`
- `/api/recipes/{id}/note` GET/PUT 엔드포인트
- 레시피 상세 페이지의 편집 모드 UI (재료 조정 배너, +/- 버튼, 메모 textarea)

**추가:**
- `cooking_logs` 테이블에 `user_id` 컬럼 추가 (마이그레이션)
- CookingLog CRUD API
- 레시피 상세 페이지에 요리 기록 섹션
