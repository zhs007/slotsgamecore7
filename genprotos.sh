docker stop slotsgamecore7genprotos
docker rm slotsgamecore7genprotos
docker build -f ./Dockerfile.genprotos -t slotsgamecore7genprotos .
docker run -v $PWD/sgc7pb:/src/app/sgc7pb --name slotsgamecore7genprotos -d slotsgamecore7genprotos