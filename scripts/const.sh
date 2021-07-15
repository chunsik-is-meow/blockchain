# src directory path
sdir=$shdir/src

# build directory path 😊
bdir=$shdir/build

# fabric-version
fv=2.2.2

CHANNEL_NAME=dna
ORDERER_ADDR=orderer0.pusan.ac.kr:7050
GLOBAL_FLAGS="-o $ORDERER_ADDR --tls --cafile /etc/hyperledger/fabric/orderer-tls/tlsca.pusan.ac.kr-cert.pem"
ORGANIZATIONS=(management verification-01 verification-02 trader)