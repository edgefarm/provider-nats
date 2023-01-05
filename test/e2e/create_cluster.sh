#!/usr/bin/env bash
set -e

# setting up colors
BLU='\033[0;34m'
YLW='\033[0;33m'
GRN='\033[0;32m'
RED='\033[0;31m'
NOC='\033[0m' # No Color
echo_info(){
    printf "\n${BLU}%s${NOC}" "$1"
}
echo_step(){
    printf "\n${BLU}> %s${NOC}\n" "$1"
}

echo_success(){
    printf "\n${GRN}%s${NOC}\n" "$1"
}
echo_warn(){
    printf "\n${YLW}%s${NOC}" "$1"
}
echo_error(){
    printf "\n${RED}%s${NOC}" "$1"
    exit 1
}

echo_step "Create cluster"
( cd config/kind ; ./create.sh )
echo_step "Deploy main nats server"
( cd config/main-nats ; ./create.sh )
echo_step "Provision nats credentials"
( cd config/nsc ; ./create.sh )
echo_step "Deploy leaf-nats servers"
( cd config/leaf-nats ; ./create.sh )
