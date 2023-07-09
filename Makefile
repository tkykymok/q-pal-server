init:
	docker compose up -d

remake:
	@make destroy
	@make init

destroy:
	docker compose down --rmi all --volumes --remove-orphans

remake-db:
	rm -rf ./docker/mysql/mysql_data
	docker compose rm -fsv db
	docker compose up db -d

exe-go:
	docker compose exec go sh

exe-db:
	docker compose exec db sh

create-models:
	docker compose exec -w /app/pkg go sqlboiler mysql models -p models --no-tests --wipe --add-global-variants
