# blockchain start
./test.sh
# blockchain API test
./blockchain.sh test

## Pre-requirement
```
sudo apt-get update
```

## Error: The required parameter 'sequence' is empty. Rerun the command with --sequence flag
docker kill $(docker ps -a -q)
docker system prune
docker volume prune
## if the same error occur...
in docker desktop: image clean up
