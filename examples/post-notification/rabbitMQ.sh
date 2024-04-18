#!/bin/bash

sudo docker run -d --rm --net rabbits -v ${PWD}/cluster_rabbitmq/config/rabbit-1/:/config/ -e RABBITMQ_CONFIG_FILE=/config/rabbitmq -e RABBITMQ_ERLANG_COOKIE=WIWVHCDTCIUAWANLMQAW --hostname rabbit-1 --name rabbit-1 -p 15672:15672 -p 5672:5672 rabbitmq:3.8-management

sudo docker run -d --rm --net rabbits -v ${PWD}/cluster_rabbitmq/config/rabbit-2/:/config/ -e RABBITMQ_CONFIG_FILE=/config/rabbitmq -e RABBITMQ_ERLANG_COOKIE=WIWVHCDTCIUAWANLMQAW --hostname rabbit-2 --name rabbit-2 -p 15673:15672 -p 5673:5672 rabbitmq:3.8-management