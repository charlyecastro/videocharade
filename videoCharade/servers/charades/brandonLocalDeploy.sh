
sh build.sh
docker rm -f info441-charades
docker run -d --name info441-charades --network info441-network -e CHARADESADDR=":80" -e RABBITADDR=$MQADDR -e MQNAME=$MQNAME -e LEADERBOARDADDR=$DSN -p 8081:80 charlyecastro/charades
