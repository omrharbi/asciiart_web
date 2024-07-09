
docker rm -f $(docker ps -aq)
docker rmi -f $(docker images -q)
docker build --tag ascii_art_web .
docker run -d -p 9090:8080 --name ascii_art_web_docker ascii_art_web
