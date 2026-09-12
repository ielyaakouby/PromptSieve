.PHONY: dev-backend dev-frontend build-api build-cli build-frontend test docker-build
.PHONY: deploy destroy docker backend frontend all build
.PHONY: deploy-backend deploy-frontend deploy-all
.PHONY: destroy-backend destroy-frontend destroy-all
.PHONY: docker-build-backend docker-build-frontend docker-build-all

CMD := $(firstword $(MAKECMDGOALS))
ARG1 := $(word 2,$(MAKECMDGOALS))
ARG2 := $(word 3,$(MAKECMDGOALS))

ifeq ($(CMD),deploy)
deploy:
	$(MAKE) --no-print-directory deploy-$(or $(ARG1),all)
backend frontend all:
	@:
endif

ifeq ($(CMD),destroy)
destroy:
	$(MAKE) --no-print-directory destroy-$(or $(ARG1),all)
backend frontend all:
	@:
endif

ifeq ($(CMD),docker)
docker:
	@if [ "$(ARG1)" != "build" ]; then echo "Usage: make docker build [backend|frontend|all]"; exit 1; fi
	$(MAKE) --no-print-directory docker-build-$(or $(ARG2),all)
build backend frontend all:
	@:
endif

deploy-backend:
	@if docker container inspect promptsieve-backend >/dev/null 2>&1; then \
		docker start promptsieve-backend; \
	else \
		docker run -d --name promptsieve-backend -p 8080:8080 promptsieve-backend:latest; \
	fi
	@echo "PromptSieve API: http://127.0.0.1:8080"

deploy-frontend:
	@if docker container inspect promptsieve-frontend >/dev/null 2>&1; then \
		docker start promptsieve-frontend; \
	else \
		docker run -d --name promptsieve-frontend -p 3005:3000 \
			-e NEXT_PUBLIC_API_URL=http://localhost:8080 \
			promptsieve-frontend:latest; \
	fi
	@echo "PromptSieve UI: http://localhost:3005"

deploy-all: deploy-backend deploy-frontend

destroy-backend:
	-docker stop promptsieve-backend
	-docker rm promptsieve-backend

destroy-frontend:
	-docker stop promptsieve-frontend
	-docker rm promptsieve-frontend

destroy-all: destroy-frontend destroy-backend

docker-build-backend:
	docker build -t promptsieve-backend:latest -f backend/Dockerfile backend

docker-build-frontend:
	docker build -t promptsieve-frontend:latest -f frontend/Dockerfile frontend

docker-build-all: docker-build-backend docker-build-frontend

docker-build: docker-build-all

dev-backend:
	cd backend && go run ./cmd/promptsieve-server

dev-frontend:
	cd frontend && npm run dev

build-api:
	cd backend && go build -o ../bin/promptsieve-server ./cmd/promptsieve-server

build-cli:
	cd backend && go build -o ../bin/promptsieve ./cmd/promptsieve

build-frontend:
	cd frontend && npm run build

test:
	cd backend && go test -race ./...
