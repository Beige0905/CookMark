# 냉장고(Pantry) 기능 컨텍스트

## 현재 구조 요약

- `my_pantry` 테이블: 재료명, 수량, 단위, 유통기한, 카테고리
- 사용자별 재고 관리

## 채택 안 한 접근법

_아직 기록 없음_

## 발견한 제약사항·함정

### pantry ↔ 레시피 재료 연결은 이름 문자열 매칭뿐
`my_pantry`와 `recipes.ingredients`(JSONB) 사이에 FK 없음. 매칭은 항상 substring 비교(`strings.Contains`). 동음이의어나 부분 일치로 과매칭 가능. 현재는 사용자가 UI에서 최종 선택하는 방식으로 허용.

### 프론트엔드 PantryItem 타입에 amount/unit 없음
DB 스키마(`my_pantry`)에는 `amount_num`, `unit`, `expiry_date`, `category` 컬럼이 있지만 프론트 타입(`PantryItem`)과 레포 쿼리는 `id`, `name`, `created_at`만 사용. 수량 추적 기능은 DB만 준비된 상태, UI 미구현.

## 트레이드오프 & 맥락

_아직 기록 없음_
