up-email-ms:
	@echo "up email-ms in progress.."
	docker compose -f ../email-ms/docker-compose.yaml up -d --build
	@echo "up email-ms done"

up-transaction-ms:
	@echo "up transaction-ms in progress.."
	docker compose -f ../transaction-ms/docker-compose.yaml up -d --build
	@echo "up transaction-ms done"

down-transaction-ms:
	@echo "down transaction-ms in progress.."
	docker compose -f ../transaction-ms/docker-compose.yaml down	

down-email-ms:
	@echo "down email-ms in progress.."
	docker compose -f ../email-ms/docker-compose.yaml down

up-all: down-transaction-ms down-email-ms up-email-ms up-transaction-ms
	@echo "all is up and running"