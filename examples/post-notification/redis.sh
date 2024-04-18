#!/bin/bash

sudo docker run -d --rm --name redis-0 --net rabbits -p 6379:6379 -v ${PWD}/clustering/redis-0:/etc/redis/ --cap-add=NET_ADMIN redis:6.0-alpine redis-server /etc/redis/redis.conf

sudo docker run -d --rm --name redis-1 --net rabbits -p 6380:6379 -v ${PWD}/clustering/redis-1:/etc/redis/ --cap-add=NET_ADMIN redis:6.0-alpine redis-server /etc/redis/redis.conf

# Wait for redis containers to initialize
sleep 10

ms="ms"

#get redis-1 ip address
IP_ADDRESS=$(sudo docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-1)

#add a delay on the redis replication process
sudo docker exec -it redis-0 apk add iproute2
sudo docker exec -it redis-0 tc qdisc del dev eth0 root
sudo docker exec -it redis-0 tc qdisc add dev eth0 root handle 1: prio
sudo docker exec -it redis-0 tc qdisc add dev eth0 parent 1:3 handle 30: netem delay $1$ms
sudo docker exec -it redis-0 tc filter add dev eth0 protocol ip parent 1:0 prio 3 u32 match ip dst $IP_ADDRESS flowid 1:3