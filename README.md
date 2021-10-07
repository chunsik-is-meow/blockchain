# blockchain start
./test.sh
or
rm -rf ./build
./blockchain all

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
sudo chown -R (USERNAME) ./ 

## if the same error occur...
in docker desktop: image clean up

## upload file directory is upload/*
## after API test, you can find downloaded files in download/* directory.