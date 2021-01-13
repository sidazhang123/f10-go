make -f "./Makefile" build
sudo docker rm -f accumulator-srv
sudo docker rmi accumulator-srv
sudo rm -rf /home/ubuntu/f10-go-volumes/accumulator-srv/logs/*
sudo docker build -t accumulator-srv:latest .
sudo docker run --name=accumulator-srv --volume="/home/ubuntu/f10-go-volumes/accumulator-srv/logs:/logs" -d --net=host --log-opt max-size=10m --log-opt max-file=3 accumulator-srv
sudo docker logs -f accumulator-srv
wait
EOF
