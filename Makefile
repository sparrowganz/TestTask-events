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
	docker run --rm -i -t -v `pwd`:`pwd` -w `pwd` --network="testtask-events_default" jordi/ab \
   -k -n 20000 -c 500 -p body.txt "http://app:8080/api/events"
