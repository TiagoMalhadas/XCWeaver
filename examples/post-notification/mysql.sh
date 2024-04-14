#!/bin/bash

docker network create mysql --subnet=174.17.0.0/16

docker run -d --net=mysql --hostname management1 -v ${PWD}/cluster_mysql/ndb_mgmd:/var/lib/mysql -v ${PWD}/cluster_mysql/my.cnf:/etc/my.cnf -v ${PWD}/cluster_mysql/mysql-cluster.cnf:/etc/mysql-cluster.cnf   --name=management1 --ip=174.17.0.2 mysql/mysql-cluster:8.0.32 ndb_mgmd --ndb-nodeid=1 --reload --initial

docker run -d --net=mysql -v ${PWD}/cluster_mysql/ndb1:/var/lib/mysql  -v ${PWD}/cluster_mysql/mysql-cluster.cnf:/etc/mysql-cluster.cnf --name=ndb1 --ip=174.17.0.3 mysql/mysql-cluster ndbd --ndb-nodeid=2 --connect-string 174.17.0.2

docker run -d --net=mysql -v ${PWD}/cluster_mysql/ndb2:/var/lib/mysql  -v ${PWD}/cluster_mysql/mysql-cluster.cnf:/etc/mysql-cluster.cnf --name=ndb2 --ip=174.17.0.4 mysql/mysql-cluster ndbd --ndb-nodeid=3 --connect-string 174.17.0.2

docker run -d -v ${PWD}/cluster_mysql/mysqld1:/var/lib/mysql  -v ${PWD}/cluster_mysql/mysql-cluster.cnf:/etc/mysql-cluster.cnf --net=mysql --name=mysql1 --ip=174.17.0.10 -e MYSQL_ROOT_PASSWORD=password mysql/mysql-cluster mysqld --ndb-nodeid=4 --ndb-connectstring 174.17.0.2

docker run -d -v ${PWD}/cluster_mysql/mysqld2:/var/lib/mysql  -v ${PWD}/cluster_mysql/mysql-cluster.cnf:/etc/mysql-cluster.cnf --net=mysql --name=mysql2 --ip=174.17.0.11 -e MYSQL_ROOT_PASSWORD=password mysql/mysql-cluster mysqld --ndb-nodeid=5 --ndb-connectstring 174.17.0.2

#docker run --name mysql-1 -e MYSQL_ROOT_PASSWORD=password -p 3306:3306 -d mysql:latest