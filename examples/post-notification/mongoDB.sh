#!/bin/bash

# Start the primary MongoDB container (master)
#sudo docker run -d --name mongo-1 --net rabbits -p 27017:27017 mongo mongod --replSet rs0
#sudo docker run -d --name mongo-1 --net rabbits -p 27017:27017 --cap-add=NET_ADMIN mongo:4.4 mongod --replSet rs0

# Start the secondary MongoDB container (replica)
#sudo docker run -d --name mongo-2 --net rabbits -p 27018:27017 mongo mongod --replSet rs0
#sudo docker run -d --name mongo-2 --net rabbits -p 27018:27017 --cap-add=NET_ADMIN mongo:4.4 mongod --replSet rs0

# Wait for MongoDB containers to initialize
#sleep 10

# Connect to the primary MongoDB container and initiate the replica set
#sudo docker exec -it mongo-1 mongosh --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-1:27017'}, {_id: 1, host: 'mongo-2:27017'}]})"
sudo docker exec -it mongo-1 mongo --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-1:27017'}, {_id: 1, host: 'mongo-2:27017'}]})"

# Print status of the replica set
echo "Replica set status:"
#sudo docker exec -it mongo-1 mongosh --eval "rs.status()"
sudo docker exec -it mongo-1 mongo --eval "rs.status()"

#ms="ms"

#get mongo-2 ip address
#IP_ADDRESS=$(sudo docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mongo-2)

#add a delay on the mongoDB replication process
#sudo docker exec -it mongo-1 apt-get update
#sudo docker exec -it mongo-1 apt-get install iproute2
#echo "delete old rules"
#sudo docker exec -it mongo-1 tc qdisc del dev eth0 root
#sudo docker exec -it mongo-1 tc qdisc add dev eth0 root handle 1: prio
#sudo docker exec -it mongo-1 tc qdisc add dev eth0 parent 1:3 handle 30: netem delay $1$ms
#sudo docker exec -it mongo-1 tc filter add dev eth0 protocol ip parent 1:0 prio 3 u32 match ip dst $IP_ADDRESS flowid 1:3
#echo "done"