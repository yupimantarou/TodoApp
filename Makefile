reload:
	docker-compose -f ./docker/docker-compose.yml down
	docker-compose -f ./docker/docker-compose.yml up -d --build
	docker ps -a
down:
	docker-compose -f ./docker/docker-compose.yml down --volumes