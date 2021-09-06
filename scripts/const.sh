# src directory path
sdir=$shdir/src

# build directory path 😊
bdir=$shdir/build

# fabric-version
fv=2.2.2

# for prod
# ORDERER_ADDR=orderer0.pusan.ac.kr:7050
# GLOBAL_FLAGS="-o $ORDERER_ADDR --tls --cafile /etc/hyperledger/fabric/orderer-tls/tlsca.pusan.ac.kr-cert.pem"
# ORGANIZATIONS=(management verification-01 verification-02 trader)
# PEERS=(peer0.management.pusan.ac.kr peer0.verification-01.pusan.ac.kr peer0.verification-02.pusan.ac.kr peer0.trader.pusan.ac.kr)
# CHANNELS=(data trade ai-model)
# CHAINCODE_DIR=/etc/hyperledger/fabric/chaincodes
# CHAINCODES=(data trade ai-model)
# VERSION=1.0

# for dev
ORDERER_ADDR=orderer0.pusan.ac.kr:7050
GLOBAL_FLAGS="-o $ORDERER_ADDR --tls --cafile /etc/hyperledger/fabric/orderer-tls/tlsca.pusan.ac.kr-cert.pem"
ORGANIZATIONS=(management verification-01 verification-02 trader)
PEERS=(peer0.management.pusan.ac.kr peer0.verification-01.pusan.ac.kr peer0.verification-02.pusan.ac.kr peer0.trader.pusan.ac.kr)
CHANNELS=(data trade ai-model)
CHAINCODE_DIR=/etc/hyperledger/fabric/chaincodes
CHAINCODES=(data trade ai-model)
VERSION=1.0