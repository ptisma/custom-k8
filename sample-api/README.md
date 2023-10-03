Simple dummy API with 2 endpoints which returns JSON response
By default version is set by env file, but this is overriden by setting the env var inside the container

docker build -t ptisma/simple-api-app .
docker run -p 8080:8080 ptisma/simple-api-app
docker push ptisma/simple-api-app