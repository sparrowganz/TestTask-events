build:
	@echo Building Application image
	docker-compose build

start:
	@echo starting web application
	docker-compose up

start-daemon:
	@echo starting web application as daemon
	docker-compose up -d

stop:
	@echo stoping web application
	docker-compose down

stress-test:
	@echo starting stress test
	ab -p body.txt -T application/json -c 300 -n 5000 http://localhost:8080/api/events