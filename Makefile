# ===========================
# DEV CONFIG
# ===========================
CONTAINER=ddd-arch-app

GO_RUN=go run
MIGRATE_CMD=$(GO_RUN) cmd/migrate/main.go
SEED_CMD=$(GO_RUN) cmd/seed/main.go

DOCKER_EXEC=docker exec $(CONTAINER)

MIGRATION_DIR=internal/infrastructure/database/postgres/migrations
NAME?=new_migration

# ===========================
# MIGRATION FILE (DEV)
# ===========================

migrate-create:
	$(DOCKER_EXEC) migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)


# ===========================
# MIGRATION COMMANDS (DEV)
# ===========================

migrate-up:
	$(DOCKER_EXEC) $(MIGRATE_CMD) -action up

migrate-down:
	$(DOCKER_EXEC) $(MIGRATE_CMD) -action down

migrate-fresh:
	$(DOCKER_EXEC) $(MIGRATE_CMD) -action fresh

# Contoh custom path
# migrate-up-path:
#	$(DOCKER_EXEC) $(MIGRATE_CMD) -action up -path "/custom/path/migrations"

# ===========================
# SEED COMMANDS (DEV)
# ===========================

seed-all:
	$(DOCKER_EXEC) $(SEED_CMD) all

seed-roles:
	$(DOCKER_EXEC) $(SEED_CMD) roles

seed-admin:
	$(DOCKER_EXEC) $(SEED_CMD) admin

seed-help:
	$(DOCKER_EXEC) $(SEED_CMD)
