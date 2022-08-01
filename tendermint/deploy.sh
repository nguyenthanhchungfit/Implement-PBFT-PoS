#!/bin/bash

APP_NAME=my-tendermint
echo "Build app ${APP_NAME}"

APP_BUILD=main2
rm -rf ${APP_BUILD}

ENTRY_POINT=${APP_BUILD}.go
go build ${ENTRY_POINT}

if [ ! -f "./${APP_BUILD}" ]; then
    echo "Build failed!"
    exit 1;
fi

DEPLOY_BASE_DIR=deploy
DEPLOY_DIRS=("node_1" "node_2" "node_3" "node_4" "node_5")

echo "Deploy app ${APP_NAME}"

# shellcheck disable=SC2068
for deployDir in ${DEPLOY_DIRS[@]}; do
  cp -uf ${APP_BUILD} ${DEPLOY_BASE_DIR}/"${deployDir}"/${APP_NAME}
done


