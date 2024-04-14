#!/bin/bash

sudo docker run -d --rm --net rabbits -v ${PWD}/cluster_rabbitmq/config/rabbit-1/:/config/ -e RABBITMQ_CONFIG_FILE=/config/rabbitmq -e RABBITMQ_ERLANG_COOKIE=WIWVHCDTCIUAWANLMQAW --hostname rabbit-1 --name rabbit-1 -p 8081:15672 rabbitmq:3.8-management

sudo docker run -d --rm --net rabbits -v ${PWD}/cluster_rabbitmq/config/rabbit-2/:/config/ -e RABBITMQ_CONFIG_FILE=/config/rabbitmq -e RABBITMQ_ERLANG_COOKIE=WIWVHCDTCIUAWANLMQAW --hostname rabbit-2 --name rabbit-2 -p 8082:15672 rabbitmq:3.8-management