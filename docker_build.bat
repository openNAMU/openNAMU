docker build . -t opennamu/opennamu:dev
docker build -t opennamu/opennamu:dev-ko -f Dockerfile.ko .
docker login
docker push opennamu/opennamu:dev
docker push opennamu/opennamu:dev-ko