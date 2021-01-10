#!/bin/sh

case "$1" in
up)
  docker network create cal_network
  docker-compose -f deployments/docker-compose.yaml up -d --build
  ;;

down)
  docker-compose -f deployments/docker-compose.yaml down
  docker network rm cal_network
  docker volume rm deployments_dbdata
  ;;

tests)
  	docker network create cal_network
  	docker-compose -f deployments/docker-compose.yaml up -d --build
  	sleep 10
  	docker-compose -f deployments/docker-compose-tests.yaml up --build
  	rc=$?
  	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose-tests.yaml down
  	docker network rm cal_network
  	docker volume rm deployments_dbdata
  	 exit $rc
  	 ;;
*)
  echo Usage ./deploy.sh up|down|tests
  ;;
esac