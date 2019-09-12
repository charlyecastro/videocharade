#! /bin/bash
cd summary
sh build.sh
cd ..
cd db
sh build.sh
cd ..
cd messaging
sh build.sh
cd ..
cd charades
sh build.sh
cd ..
cd gateway 
sh build.sh

docker push charlyecastro/pagesummaryapi
docker push charlyecastro/mysqlstore
docker push charlyecastro/summary
docker push charlyecastro/messaging
docker push charlyecastro/charades

cd ~/.ssh
export TLSCERT=/etc/letsencrypt/live/api.videocharade.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.videocharade.me/privkey.pem
export REDISADDR="redisServer:6379"
export MYSQL_ADDR="mysqlstore:3306"
export MYSQL_DATABASE="userDB"
export SESSIONKEY="SigningKey"
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
export MESSAGESADDR="messaging:80"
export SUMMARYADDR="summary:80"
export MONGOADDR="mongoserver:27017"
export RABBITADDR="rabbitmq:5672"
export RABBITNAME="Messaging"
export CHARADESADDR="info441-charades:80"
export CHARADESADDR="info441-charades:80"

ssh -i 441finalproject.pem ec2-user@35.168.189.93 'bash -s' << EOF
sudo service docker start

docker pull charlyecastro/pagesummaryapi
docker pull charlyecastro/mysqlstore
docker pull charlyecastro/summary
docker pull charlyecastro/messaging
docker pull charlyecastro/charades

docker rm -f pagesummaryapi
docker rm -f summary
docker rm -f messaging
docker rm -f redisServer
docker rm -f mysqlstore
docker rm -f mongoserver
docker rm -f rabbitmq
docker rm -f info441-charades

docker network rm messagesnetwork

docker network create messagesnetwork

docker run -d --name mongoserver --network messagesnetwork mongo

docker run -d  --name rabbitmq  --network messagesnetwork rabbitmq:3-management

docker run -d --name mysqlstore -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD --network messagesnetwork charlyecastro/mysqlstore

sleep 20

docker run -d --name redisServer --network messagesnetwork redis

docker run -d --name summary --network messagesnetwork charlyecastro/summary

docker run -d --name messaging -e MESSAGESADDR=$MESSAGESADDR -e MONGOADDR=$MONGOADDR -e RABBITADDR=$RABBITADDR --network messagesnetwork charlyecastro/messaging

docker run -d --name info441-charades --network messagesnetwork  -e CHARADESADDR=$CHARADESADDR -e RABBITADDR=$RABBITADDR -e MQNAME=$RABBITNAME -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD  -e MYSQL_ADDR=$MYSQL_ADDR -e MYSQL_DATABASE=$MYSQL_DATABASE charlyecastro/charades

docker run -d --name pagesummaryapi -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e REDISADDR=$REDISADDR -e SESSIONKEY=$SESSIONKEY -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD  -e MYSQL_ADDR=$MYSQL_ADDR -e MYSQL_DATABASE=$MYSQL_DATABASE -e SESSIONKEY=$SESSIONKEY -e SUMMARYADDR=$SUMMARYADDR -e MESSAGESADDR=$MESSAGESADDR -e RABBITADDR=$RABBITADDR -e CHARADESADDR=$CHARADESADDR --network messagesnetwork charlyecastro/pagesummaryapi
exit

EOF
