## Installation
```
docker pull opennamu/opennamu
```

## Start
```
docker run -p 3000:3000 -v data:/app/data --name opennamu opennamu/opennamu
docker run -p <host-port>:3000 -v <host-data_directory>:/app/data --name <docker-containername> opennamu/opennamu
```