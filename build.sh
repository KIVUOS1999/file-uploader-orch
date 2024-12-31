set -e

echo "Building file-uploader-orch binary"

go build -o file-uploader-orch

echo "Login and pusing to docker-hub"

docker login
docker build -t kivuos1999/file-uploader-orch .
docker push kivuos1999/file-uploader-orch

echo "cleanup"
rm file-uploader-orch

echo "build succeed"