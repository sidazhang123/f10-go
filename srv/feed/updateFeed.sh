make -f "./Makefile" build
sudo docker rm -f feed-srv
sudo docker rmi feed-srv
sudo rm -rf /home/ubuntu/f10-go-volumes/feed-srv/logs/*
sudo docker build -t feed-srv:latest .
sudo docker run --name=feed-srv --volume="/home/ubuntu/f10-go-volumes/feed-srv/logs:/logs" \
--volume="/home/ubuntu/ad:/ad" \
-d --net=host --log-opt max-size=10m --log-opt max-file=3 feed-srv
sudo docker logs -f feed-srv
wait
EOF
