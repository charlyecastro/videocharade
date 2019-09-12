#! /bin/bash
sh ./build.sh
docker push charlyecastro/charadesapp
cd ~/.ssh
ssh -i 441finalproject.pem ec2-user@23.22.236.54 'bash -s' << EOF
docker rm -f charadesApp
docker pull charlyecastro/charadesapp
docker run -d --name charadesApp -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro charlyecastro/charadesapp
exit

EOF