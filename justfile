# 기본 명령어: 모든 서비스 실행 (백그라운드)
up:
    docker-compose up -d

# 백엔드와 DB만 실행 (프론트는 로컬에서 npm run dev 할 때 사용, 백그라운드)
backend:
    docker-compose up -d db backend

# DB만 실행 (백그라운드)
db:
    docker-compose up -d db

# 로컬 개발 환경 (DB 실행 후 backend에서 air 실행)
dev: db
    cd backend && air

# 개발용 Docker (프론트 HMR 포함, http://localhost:5173)
up-dev:
    docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

# 컨테이너 끄기
down:
    docker-compose down

# DB 로그만 확인
logs-db:
    docker-compose logs -f db

# DB 마이그레이션 실행 (예: just db-migrate 003_add_recipe_notes.sql)
db-migrate file:
    docker-compose exec -T db psql -U ${POSTGRES_USER:-recipe_user} -d ${POSTGRES_DB:-recipe_db} < db/migrations/{{file}}

# Go 백엔드 테스트 실행
test:
    cd backend && go test ./...

# 프론트엔드 타입체크
typecheck:
    cd frontend && npm run check

# 배포 전 로컬 검증 (테스트 + 타입체크)
check: test typecheck

# Railway 헬스체크
health:
    @curl -sf https://cookmark-production.up.railway.app/health \
        && echo "✓ healthy" \
        || echo "✗ unhealthy"

# 배포 완료 대기 (Railway 빌드 후 헬스체크 폴링)
wait-deploy:
    @echo "Railway 빌드 대기 중... (약 2-3분 소요)"
    @for i in $(seq 1 24); do \
        sleep 10; \
        printf "[$$(expr $i \* 10)s] "; \
        if curl -sf https://cookmark-production.up.railway.app/health > /dev/null 2>&1; then \
            echo "✓ 배포 완료"; \
            exit 0; \
        fi; \
        echo "대기 중..."; \
    done; \
    echo "⚠ 4분 초과 - Railway 대시보드 확인"

# 배포 파이프라인 (체크 → 푸시 → 완료 대기)
deploy: check
    git push origin main
    @echo "✓ 푸시 완료 - Railway 자동 배포 시작"
    just wait-deploy
