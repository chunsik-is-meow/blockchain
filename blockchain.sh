#!/bin/bash

# shell directory path
shdir="$( cd "$(dirname "$0")" ; pwd -P)"

# iclude const, utils scripts
source $shdir/scripts/const.sh
source $shdir/scripts/utils.sh
export COMPOSE_IGNORE_ORPHANS=True

function bockchain_usage {
    cecho "RED" "not implemantation now"
}

function blockchain_all {
    blockchain_build
    blockchain_up
    blockchain_channel
}

function blockchain_clean {
    blockchain_down
    rm -rf $bdir
}

function blockchain_build {
    # isExist build
    # -> print err / exit

    # else
    mkdir -p $bdir &&  \
    cp -rf $sdir/asset $bdir/ && \
    mkdir -p $bdir/asset/artifacts/block && \
    mkdir -p $bdir/asset/artifacts/tx >> /dev/null 2>&1

    blockchain_build_cryptogen
    blockchain_build_configtxgen
}

function blockchain_build_cryptogen {
    command "docker run -it --rm \
    -v $bdir:/workdir \
    --workdir /workdir \
    hyperledger/fabric-tools:$fv \
    cryptogen generate --config=./asset/tool/crypto-config.yaml --output=./asset/artifacts/crypto-config"
}

function blockchain_build_configtxgen {
    blockchain_build_configtxgen_genesis_block
    blockchain_build_configtxgen_channel_tx
}

function blockchain_build_configtxgen_genesis_block {
    command "docker run -it --rm \
    -v $bdir/asset/artifacts/block:/workdir/block \
    -v $bdir/asset/artifacts/crypto-config:/workdir/crypto-config \
    -v $bdir/asset/tool/system-channel-configtx.yaml:/workdir/configtx.yaml \
    --workdir /workdir \
    hyperledger/fabric-tools:$fv \
    configtxgen -profile system-channelProfile -channelID system-channel -outputBlock /workdir/block/system-channel.block -configPath /workdir"
}

function blockchain_build_configtxgen_channel_tx {
    command "docker run -it --rm \
    -v $bdir/asset/artifacts/tx:/workdir/tx \
    -v $bdir/asset/artifacts/crypto-config:/workdir/crypto-config \
    -v $bdir/asset/tool/dna-configtx.yaml:/workdir/configtx.yaml \
    --workdir /workdir \
    hyperledger/fabric-tools:$fv \
    configtxgen -profile dnaProfile -channelID dna -outputCreateChannelTx /workdir/tx/dna.tx -configPath /workdir"
}

function blockchain_up {
    docker network create pnu >> /dev/null 2>&1
    for TARGET in "$bdir/asset/docker"/*
    do
        command "docker-compose -f $TARGET up -d"
    done

    cecho "INFO" "Waiting 10s for blockchain network stable"
    sleep 10s
}

function blockchain_down {
    for TARGET in "$bdir/asset/docker"/*
    do
        command "docker-compose -f $TARGET down"
    done
    docker network rm pnu >> /dev/null 2>&1
    # chaincode clean
    yes | docker volume prune
}

function blockchain_channel {
    blockchain_channel_create
    blockchain_channel_join
}

function blockchain_channel_create {
    command "docker exec -it \
    cli.peer0.management.pusan.ac.kr \
    peer channel create -c $CHANNEL_NAME -f /etc/hyperledger/fabric/tx/$CHANNEL_NAME.tx --outputBlock /etc/hyperledger/fabric/block/$CHANNEL_NAME.block $GLOBAL_FLAGS"
}

function blockchain_channel_join {
    for ORG in ${ORGANIZATIONS[@]}
    do
        command "docker exec -it \
        cli.peer0.$ORG.pusan.ac.kr \
        peer channel join -b /etc/hyperledger/fabric/block/$CHANNEL_NAME.block"
    done
}

function main {
    case $1 in
        all | clean | build | up | down | channel)
            cmd=blockchain_$1
            $cmd
            ;;
        *)
            bockchain_usage
            exit
            ;;
    esac
}

main $@