* i have picked the docker-compose configuration file from 
    https://hub.docker.com/r/bitnami/kafka

docker-compose up -d

go run main.go -- for data receiver 

(
    if you get this error 
    %4|1731830689.740|TERMINATE|rdkafka#producer-1| [thrd:app]: Producer terminating with 1 message (4 bytes) still in queue or transit: use flush() to wait for outstanding message delivery
    then  you are connected to kafka client
)