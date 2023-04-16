#! /bin/bash
#
docker-compose down
docker volume prune
docker volume rm explorer_pgdata
docker volume rm explorer_walletstore
docker network prune
