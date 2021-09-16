#!/bin/bash

# shell directory path
shdir="$( cd "$(dirname "$0")" ; pwd -P)"

# iclude utils scripts
source $shdir/scripts/utils.sh
export COMPOSE_IGNORE_ORPHANS=True

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

function blockchain_chaincode_package {
    CHAINCODE_NAME=$1
    VERSION=$2

    command "docker exec -it \
    cli.peer0.management.pusan.ac.kr \
    peer lifecycle chaincode package $CHAINCODE_DIR/$CHAINCODE_NAME-$VERSION.tar.gz --path $CHAINCODE_DIR/$CHAINCODE_NAME --lang golang --label $CHAINCODE_NAME-$VERSION"
}

function blockchain_chaincode_install {
    PEER_NAME=$1
    CHAINCODE_NAME=$2
    VERSION=$3

    command "docker exec -it \
    cli.$PEER_NAME \
    peer lifecycle chaincode install $CHAINCODE_DIR/$CHAINCODE_NAME-$VERSION.tar.gz"
}

function blockchain_chaincode_getpackageid {
    command "docker exec -it \
    cli.$1 \
    peer lifecycle chaincode queryinstalled"

    PACKAGE_ID=$(sed -n "/$2-$3/{s/^Package ID: //; s/, Label:.*$//; p;}" $bdir/log.txt)
}

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

function blockchain_chaincode {
    VERSION=1.0
    SEQUENCE=1

    for CHAINCODE_NAME in ${CHAINCODES[@]}
    do
        blockchain_chaincode_package $CHAINCODE_NAME $VERSION
        for PEER_NAME in ${PEERS[@]}
        do
            if [[ $PEER_NAME == 'peer0.trader.pusan.ac.kr' && $CHAINCODE_NAME == 'ai-model' ]]
            then
                continue
            fi
                blockchain_chaincode_install $PEER_NAME $CHAINCODE_NAME $VERSION
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
                blockchain_chaincode_approveformyorg $PEER_NAME $CHANNEL $CHAINCODE_NAME $VERSION $SEQUENCE
                sleep 1s
                blockchain_chaincode_checkcommitreadiness $PEER_NAME $CHANNEL $CHAINCODE_NAME $VERSION $SEQUENCE
        done
        blockchain_chaincode_commit 'peer0.management.pusan.ac.kr' $CHANNEL $CHAINCODE_NAME $VERSION $SEQUENCE
        blockchain_chaincode_querycommitted 'peer0.management.pusan.ac.kr' $CHANNEL
    done

    for CHANNEL in ${CHANNELS[@]}
    do
        blockchain_chaincode_init $CHANNEL
    done
}

function blockchain_upgrade {
    blockchain_chaincode_upgrade $@
}

function blockchain_test {
    blockchain_test_$1
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

function blockchain_chaincode_init {
    # peer=$1
    # channel=$2
    # chaincode=$3
    peer=peer0.management.pusan.ac.kr
    channel=$1
    chaincode=$1
    fcn_call='{"function":"InitLedger","Args":[]}'

    command "docker exec -it \
    cli.$peer \
    peer chaincode invoke  \
    --channelID $channel \
    --name $chaincode \
    -c $fcn_call \
    $GLOBAL_FLAGS"
}

function blockchain_chaincode_invoke {
    # peer, channel=$1
    # chaincode=$2
    peer=peer0.management.pusan.ac.kr
    channel=$1
    chaincode=$1

    command "docker exec -it \
    cli.$peer \
    peer chaincode invoke  \
    --channelID $channel \
    --name $chaincode \
    -c $2 \
    $GLOBAL_FLAGS"
}

function blockchain_chaincode_query {
    # peer, channel=$1
    # chaincode=$2
    peer=peer0.management.pusan.ac.kr
    channel=$1
    chaincode=$1

    command "docker exec -it \
    cli.$peer \
    peer chaincode query  \
    --channelID $channel \
    --name $chaincode \
    -c $2"
}

function blockchain_chaincode_upgrade {
    CHANNEL=$1
    CHAINCODE_NAME=$1
    VERSION=$2
    SEQUENCE=$3

    command "rm -rf $bdir/asset/chaincodes/${CHAINCODE_NAME}"
    command "cp -rf $sdir/asset/chaincodes/${CHAINCODE_NAME} $bdir/asset/chaincodes/${CHAINCODE_NAME}"

    blockchain_chaincode_package $CHAINCODE_NAME $VERSION
    for PEER_NAME in ${PEERS[@]}
    do
        if [[ $PEER_NAME == 'peer0.trader.pusan.ac.kr' && $CHAINCODE_NAME == 'ai-model' ]]
        then
            continue
        fi
        blockchain_chaincode_install $PEER_NAME $CHAINCODE_NAME $VERSION
    done

    for PEER_NAME in ${PEERS[@]}
    do
        if [[ $PEER_NAME == 'peer0.trader.pusan.ac.kr' && $CHAINCODE_NAME == 'ai-model' ]]
        then
            continue
        fi
        blockchain_chaincode_approveformyorg $PEER_NAME $CHANNEL $CHAINCODE_NAME $VERSION $SEQUENCE
        sleep 1s
        blockchain_chaincode_checkcommitreadiness $PEER_NAME $CHANNEL $CHAINCODE_NAME $VERSION $SEQUENCE
    done
    blockchain_chaincode_commit 'peer0.management.pusan.ac.kr' $CHANNEL $CHAINCODE_NAME $VERSION $SEQUENCE
}

function file_upload {
    channel=$1
    temp=$2
    if [ $channel = "data" ]; then
        contents=`cat $temp`
        echo $contents | tr ' ' '!' > up.txt
    else
        contents=`hexdump -v -e '/1 "%02X!"' $temp`
        echo $contents > up.txt
    fi
    FILECONTENTS=`cat up.txt`
    rm up.txt
}

function file_download {
    # peer, channel=$1
    # downloaded file=$2
    # chaincode=$3
    peer=peer0.management.pusan.ac.kr
    channel=$1
    chaincode=$1
    file=$2
    
    command "docker exec -it \
    cli.$peer \
    peer chaincode invoke  \
    --channelID $channel \
    --name $chaincode \
    -c $3 \
    $GLOBAL_FLAGS" > down.txt
    
    con=$(head -2 down.txt | tail -1)
    echo = ${con:144:-3} 1> down.txt
    contents=$(cat down.txt)
    if [ $channel = "data" ]; then
        echo $contents | tr '!' '\n' | tr -d '= ' 1> download/$channel/$file.csv
    else
        down=`echo $contents | tr -d '!' | tr -d '= '`
        echo -n $down | xxd -r -p  1> download/$channel/$file.h5
    fi
    rm down.txt
}

function blockchain_init {
    date=$(date '+%Y-%m-%d-%H-%M-%S')
    price=1000

    # file upload
    file_upload data upload/data/iris.csv
    blockchain_chaincode_invoke data '{"function":"PutCommonData","Args":["yohan","iris","1.0","iris_classfication","R.A.Fisher","'$FILECONTENTS'","'$date'"]}'
    
    sleep 3s
    file_upload data upload/data/wine.csv
    blockchain_chaincode_invoke data '{"function":"PutCommonData","Args":["hyoeun","wine","1.2","wine_classfication","PARVUS","'$FILECONTENTS'","'$date'"]}'
    
    sleep 3s
    file_upload data upload/data/cancer.csv
    blockchain_chaincode_invoke data '{"function":"PutCommonData","Args":["yohan","cancer","2.0","cancer_classfication","L.Mangasarian.","'$FILECONTENTS'","'$date'"]}'

    file_upload ai-model upload/ai-model/test_model.h5
    blockchain_chaincode_invoke ai-model '{"function":"PutAIModel","Args":["hyoeun","test_model","1.0","Python","'$price'","CCC","test_input","'$FILECONTENTS'","'$date'"]}'
    sleep 2s

    file_upload ai-model upload/ai-model/model_test.h5
    blockchain_chaincode_invoke ai-model '{"function":"PutAIModel","Args":["yohan","model_test","2.0","Python","'$price'","AAA","input_test","'$FILECONTENTS'","'$date'"]}'
    sleep 2s


    blockchain_chaincode_invoke trade '{"function":"Transfer","Args":["bank","hyoeun","300000","'$date'","transfer"]}'
    blockchain_chaincode_invoke trade '{"function":"Transfer","Args":["bank","yohan","300000","'$date'","transfer"]}'
    sleep 2s

    blockchain_chaincode_invoke trade '{"function":"BuyModel","Args":["yohan","A_hyoeun_test_model_1.0","'$price'","'$date'"]}'
}

function blockchain_test_trade {
    date=$(date '+%Y-%m-%d-%H-%M-%S')
    price=100
    ################################################### trade chaincode ####################################################

    blockchain_chaincode_query trade '{"function":"GetCurrentMeow","Args":["hyoeun"]}'
    blockchain_chaincode_query trade '{"function":"GetCurrentMeow","Args":["yohan"]}'

    # NOTE meow is lacking error
    blockchain_chaincode_invoke trade '{"function":"Transfer","Args":["hyoeun","gydms","30","'$date'","transfer"]}'

    blockchain_chaincode_invoke trade '{"function":"Transfer","Args":["bank","hyoeun","300000","'$date'","transfer"]}'
    sleep 2s

    blockchain_chaincode_query trade '{"function":"GetCurrentMeow","Args":["hyoeun"]}'
    blockchain_chaincode_query trade '{"function":"GetCurrentMeow","Args":["yohan"]}'

    blockchain_chaincode_query trade '{"function":"GetQueryHistory","Args":["hyoeun"]}'

    # blockchain_chaincode_invoke trade '{"function":"GetModel","Args":["A_hyoeun_test_model_1.0"]}'
}

function blockchain_test_data {
    date=$(date '+%Y-%m-%d-%H-%M-%S')
    price=100
   
    ################################################## data chaincode ####################################################
    blockchain_chaincode_query data '{"function":"GetAllCommonDataInfo","Args":[]}'
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["DC"]}'
    
    # file upload
    file_upload data upload/data/iris.csv
    blockchain_chaincode_invoke data '{"function":"PutCommonData","Args":["yohan","iris","1.0","iris_classfication","R.A.Fisher","'$FILECONTENTS'","'$date'"]}'
    
    sleep 3s
    file_upload data upload/data/wine.csv
    blockchain_chaincode_invoke data '{"function":"PutCommonData","Args":["hyoeun","wine","1.2","wine_classfication","PARVUS","'$FILECONTENTS'","'$date'"]}'
    
    sleep 3s
    file_upload data upload/data/cancer.csv
    blockchain_chaincode_invoke data '{"function":"PutCommonData","Args":["yohan","cancer","2.0","cancer_classfication","L.Mangasarian.","'$FILECONTENTS'","'$date'"]}'
    
    sleep 3s
    # get datainfo
    blockchain_chaincode_query data '{"function":"GetAllCommonDataInfo","Args":[]}'
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["DC"]}'

    blockchain_chaincode_query data '{"function":"GetCommonDataInfo","Args":["yohan","cancer","2.0"]}'

    #file download
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["user1"]}'
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["user2"]}'
    download_file="iris"
    file_download data $download_file '{"function":"GetCommonDataContents","Args":["yohan","'$download_file'","1.0","user1"]}'
    download_file="wine"
    file_download data $download_file '{"function":"GetCommonDataContents","Args":["hyoeun","'$download_file'","1.2","user1"]}'
    download_file="cancer"
    file_download data $download_file '{"function":"GetCommonDataContents","Args":["yohan","'$download_file'","2.0","user2"]}'
    sleep 3s

    # count
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["DC"]}'
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["user1"]}'
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["user2"]}'

    # NOTE data is exist error
    blockchain_chaincode_invoke data '{"function":"PutCommonData","Args":["yohan","iris","1.0","iris_classfication","R.A.Fisher","aaaaa","'$date'"]}'
}

function blockchain_test_ai {
    date=$(date '+%Y-%m-%d-%H-%M-%S')
    price=100
    
    #################################################### ai-model chaincode ####################################################
    blockchain_chaincode_query ai-model '{"function":"GetAllAIModelInfo","Args":[]}'
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["AC"]}'
    
    file_upload ai-model upload/ai-model/test_model.h5
    blockchain_chaincode_invoke ai-model '{"function":"PutAIModel","Args":["hyoeun","test_model","1.0","Python","'$price'","CCC","test_input","'$FILECONTENTS'","'$date'"]}'
    file_upload ai-model upload/ai-model/model_test.h5
    blockchain_chaincode_invoke ai-model '{"function":"PutAIModel","Args":["yohan","model_test","2.0","Python","'$price'","AAA","input_test","'$FILECONTENTS'","'$date'"]}'
    sleep 3s

    # get ai model info
    blockchain_chaincode_query ai-model '{"function":"GetAllAIModelInfo","Args":[]}'
    blockchain_chaincode_query ai-model '{"function":"GetAIModelInfo","Args":["hyoeun","test_model","1.0"]}'

    # file download
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["user1"]}'
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["user2"]}'
    download_file="test_model"
    file_download ai-model $download_file '{"function":"GetAIModelContents","Args":["hyoeun","'$download_file'","1.0","user1"]}'
    file_download ai-model $download_file '{"function":"GetAIModelContents","Args":["hyoeun","'$download_file'","1.0","user2"]}'
    download_file="model_test"
    file_download ai-model $download_file '{"function":"GetAIModelContents","Args":["yohan","'$download_file'","2.0","user1"]}'
    sleep 3s

    # count
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["AC"]}'
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["user1"]}'
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["user2"]}'

    # NOTE ai-model is exist error
    blockchain_chaincode_invoke ai-model '{"function":"PutAIModel","Args":["hyoeun","test_model","1.0","C","'$price'","CCC","iris_learning","aaaaa","'$date'"]}'
    
}

function blockchain_check { 
    ################################################## data chaincode ####################################################
    # get datainfo
    blockchain_chaincode_query data '{"function":"GetAllCommonDataInfo","Args":[]}'
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["DC"]}'

    blockchain_chaincode_query data '{"function":"GetCommonDataInfo","Args":["yohan","cancer","2.0"]}'
    #file download
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["user1"]}'
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["user2"]}'
    download_file="iris"
    file_download data $download_file '{"function":"GetCommonDataContents","Args":["yohan","'$download_file'","1.0","hyoeun"]}'
    download_file="wine"
    file_download data $download_file '{"function":"GetCommonDataContents","Args":["hyoeun","'$download_file'","1.2","user1"]}'
    download_file="cancer"
    file_download data $download_file '{"function":"GetCommonDataContents","Args":["yohan","'$download_file'","2.0","user2"]}'

    # count
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["DC"]}'
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["user1"]}'
    blockchain_chaincode_query data '{"function":"GetCommonDataCount","Args":["user2"]}'


    #################################################### ai-model chaincode ####################################################
     # get ai model info
    blockchain_chaincode_query ai-model '{"function":"GetAllAIModelInfo","Args":[]}'
    blockchain_chaincode_query ai-model '{"function":"GetAIModelInfo","Args":["hyoeun","test_model","1.0"]}'

    # file download
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["user1"]}'
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["user2"]}'
    download_file="test_model"
    file_download ai-model $download_file '{"function":"GetAIModelContents","Args":["hyoeun","'$download_file'","1.0","hyoeun"]}'
    file_download ai-model $download_file '{"function":"GetAIModelContents","Args":["hyoeun","'$download_file'","1.0","user2"]}'
    download_file="model_test"
    file_download ai-model $download_file '{"function":"GetAIModelContents","Args":["yohan","'$download_file'","2.0","user1"]}'

    # count
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["AC"]}'
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["user1"]}'
    blockchain_chaincode_query ai-model '{"function":"GetAIModelCount","Args":["user2"]}'
}


function main {
    case $1 in
        all | clean | build | up | down | channel | chaincode | test | upgrade | check | init)
            cmd=blockchain_$1
            shift
            $cmd $@
            ;;
        *)
            bockchain_usage
            exit
            ;;
    esac
}

main $@
