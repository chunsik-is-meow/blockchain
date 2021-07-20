#!/bin/bash

# shell directory path
shdir="$( cd "$(dirname "$0")" ; pwd -P)"

# iclude utils scripts
source $shdir/scripts/utils.sh
export COMPOSE_IGNORE_ORPHANS=True

function bockchain_usage {
    cecho "RED" "not implemantation now"
}

function blockchain_all {
    blockchain_build
    blockchain_up
    blockchain_channel
    blockchain_chaincode
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
    for CHANNEL_NAME in ${CHANNELS[@]}
    do
        blockchain_build_configtxgen_channel_tx $CHANNEL_NAME
    done
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
    -v $bdir/asset/tool/$1-configtx.yaml:/workdir/configtx.yaml \
    --workdir /workdir \
    hyperledger/fabric-tools:$fv \
    configtxgen -profile $1Profile -channelID $1 -outputCreateChannelTx /workdir/tx/$1.tx -configPath /workdir"
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
    docker image rm -f $(docker images -aq --filter reference='pnu-peer*') 2>/dev/null || true
    yes | docker volume prune
}

function blockchain_channel {
    for CHANNEL_NAME in ${CHANNELS[@]}
    do
        blockchain_channel_create $CHANNEL_NAME
        sleep 3s
        for ORG_NAME in ${ORGANIZATIONS[@]}
        do
            if [[ $ORG_NAME == 'trader' && $CHANNEL_NAME == 'ai-model' ]]
            then
                continue
            fi
            blockchain_channel_join $ORG_NAME $CHANNEL_NAME
        done
    done
}

function blockchain_channel_create {
    command "docker exec -it \
    cli.peer0.management.pusan.ac.kr \
    peer channel create -c $1 -f /etc/hyperledger/fabric/tx/$1.tx --outputBlock /etc/hyperledger/fabric/block/$1.block $GLOBAL_FLAGS"
}

function blockchain_channel_join {
    command "docker exec -it \
    cli.peer0.$1.pusan.ac.kr \
    peer channel join -b /etc/hyperledger/fabric/block/$2.block"
}

function blockchain_chaincode {
    for CHAINCODE_NAME in ${CHAINCODES[@]}
    do
        blockchain_chaincode_package $CHAINCODE_NAME
        for PEER_NAME in ${PEERS[@]}
        do
            if [[ $PEER_NAME == 'peer0.trader.pusan.ac.kr' && $CHAINCODE_NAME == 'ai-model' ]]
            then
                continue
            fi
                blockchain_chaincode_install $PEER_NAME $CHAINCODE_NAME
        done
    done

    for CHANNEL in ${CHANNELS[@]}
    do
        CHAINCODE_NAME=$CHANNEL
        for PEER_NAME in ${PEERS[@]}
        do
            if [[ $PEER_NAME == 'peer0.trader.pusan.ac.kr' && $CHAINCODE_NAME == 'ai-model' ]]
            then
                continue
            fi
                blockchain_chaincode_approveformyorg $PEER_NAME $CHANNEL $CHAINCODE_NAME $VERSION 1
                sleep 1s
                blockchain_chaincode_checkcommitreadiness $PEER_NAME $CHANNEL $CHAINCODE_NAME $VERSION 1
        done
        blockchain_chaincode_commit 'peer0.management.pusan.ac.kr' $CHANNEL $CHAINCODE_NAME $VERSION 1
        blockchain_chaincode_querycommitted 'peer0.management.pusan.ac.kr' $CHANNEL
    done
}

function blockchain_chaincode_package {
    command "docker exec -it \
    cli.peer0.management.pusan.ac.kr \
    peer lifecycle chaincode package $CHAINCODE_DIR/$1-$VERSION.tar.gz --path $CHAINCODE_DIR/$1 --lang golang --label $1-$VERSION"
}

function blockchain_chaincode_install {
    command "docker exec -it \
    cli.$1 \
    peer lifecycle chaincode install $CHAINCODE_DIR/$2-$VERSION.tar.gz"
}

function blockchain_chaincode_approveformyorg {
    peer=$1
    channel=$2
    chaincode=$3
    version=$4
    sequence=$5
    blockchain_chaincode_getpackageid $peer $chaincode $version

    command "docker exec -it \
    cli.$peer \
    peer lifecycle chaincode approveformyorg \
    --channelID $channel \
    --name $chaincode \
    --version $version \
    --package-id $PACKAGE_ID \
    --sequence $sequence \
    $GLOBAL_FLAGS"
}

function blockchain_chaincode_checkcommitreadiness {
    peer=$1
    channel=$2
    chaincode=$3
    version=$4
    sequence=$5

    command "docker exec -it \
    cli.$peer \
    peer lifecycle chaincode checkcommitreadiness  \
    --channelID $channel \
    --name $chaincode \
    --version $version \
    --sequence $sequence \
    $GLOBAL_FLAGS"
}

function blockchain_chaincode_commit {
    peer=$1
    channel=$2
    chaincode=$3
    version=$4
    sequence=$5

    command "docker exec -it \
    cli.$peer \
    peer lifecycle chaincode commit  \
    --channelID $channel \
    --name $chaincode \
    --version $version \
    --sequence $sequence \
    $GLOBAL_FLAGS"
}

function blockchain_chaincode_querycommitted {
    peer=$1
    channel=$2

    command "docker exec -it \
    cli.$peer \
    peer lifecycle chaincode querycommitted  \
    --channelID $channel \
    $GLOBAL_FLAGS"
}

function blockchain_chaincode_getpackageid {
    command "docker exec -it \
    cli.$1 \
    peer lifecycle chaincode queryinstalled"

    PACKAGE_ID=$(sed -n "/$2-$3/{s/^Package ID: //; s/, Label:.*$//; p;}" $bdir/log.txt)
}

function main {
    case $1 in
        all | clean | build | up | down | channel | chaincode )
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