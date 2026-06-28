# Project Wiki

대화 중 발생하는 지식 증발을 방지하기 위한 위키입니다.
`docs/decisions/`가 공식 설계 결정을 담당하고, 이 위키는 "왜"와 "맥락"을 담습니다.

## 페이지 목록

- [constraints.md](constraints.md) — 전역 제약사항 및 함정 (DB, 라이브러리 quirk 등)
- [recipes.md](recipes.md) — 레시피 기능 컨텍스트
- [fridge.md](fridge.md) — 냉장고(pantry) 기능 컨텍스트
- [cooking-log.md](cooking-log.md) — 쿠킹 로그 기능 컨텍스트
- [auth.md](auth.md) — 인증 기능 컨텍스트

## 이 위키가 만들어진 배경

Karpathy의 LLM Wiki 패턴(gist: 442a6bf)을 참고해서 설계. 핵심 문제는 대화 중 나온 트레이드오프 분석, 채택 안 한 접근법, 발견한 제약사항이 세션 끝나면 증발하는 것.

`docs/decisions/`와 역할 분리: decisions = 공식 설계 결정, wiki = "왜"와 맥락.

### 채택 안 한 접근법

**앱 기능으로 활용** (유저 요리 경험 기반 개인 위키) → 개발 도구가 더 시급해서 기각

**위키를 ~/.claude 메모리 시스템에 두기** → git 버전 관리와 팀 공유를 위해 docs/wiki/ 선택

**유형별 페이지 구조** (rejected-approaches.md, constraints.md 등으로만 구성) → 특정 기능 관련 내용을 모아보기 어려워서 기각, 기능별 페이지 선택

**기능×주제 매트릭스** (fridge/constraints.md 등) → 파일이 너무 많아져서 기각

**외부 API 호출 Stop 훅** (셸 스크립트가 Claude API 직접 호출) → 구현 복잡도 대비 실익 불확실, CLAUDE.md 지시로 대체

### Claude Code 스킬 파일 규칙

스킬 파일은 반드시 `SKILL.md`로 명명해야 하고, `name`과 `description` 프론트매터가 필요함. 스킬 디렉토리명이 아니라 프론트매터의 `name` 값이 `/커맨드명`으로 사용됨. `wiki.md`처럼 다른 이름으로 만들면 Skill 도구가 인식하지 못함.

---

## 기록 기준

다음 중 하나에 해당하면 위키에 기록합니다:
1. **채택 안 한 접근법** — 검토했지만 왜 안 썼는지
2. **제약사항·함정** — 디버깅이나 구현 중 발견한 비자명한 동작
3. **기능 컨텍스트** — 지금 구조가 된 이유, 트레이드오프
