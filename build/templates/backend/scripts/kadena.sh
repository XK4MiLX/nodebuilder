#!/bin/sh

{{define "main" -}}

set -e
export CHAINWEB_NETWORK=${CHAINWEB_NETWORK:-mainnet01}
export CHAINWEB_P2P_PORT=${CHAINWEB_P2P_PORT:-1789}
export CHAINWEB_SERVICE_PORT=${CHAINWEB_SERVICE_PORT:-1848}
export LOGLEVEL=${LOGLEVEL:-warn}
export MINER_KEY=${MINER_KEY:-}
export MINER_ACCOUNT=${MINER_ACCOUNT:-$MINER_KEY}
export ENABLE_ROSETTA=${ENABLE_ROSETTA:-}

if [[ -z "$CHAINWEB_P2P_HOST" ]] ; then
    CHAINWEB_P2P_HOST="0.0.0.0"
fi
export CHAINWEB_P2P_HOST

INSTALL_DIR={{.Env.BackendInstallPath}}/{{.Coin.Alias}}
DATA_DIR={{.Env.BackendDataPath}}/{{.Coin.Alias}}/backend
BIN=$INSTALL_DIR/chainweb-node

mkdir -p "$DATA_DIR/0"

if [[ -z "$MINER_KEY" ]] ; then
export MINER_CONFIG="
chainweb:
  mining:
    coordination:
      enabled: ${MINING_ENABLED:-false}
"
else
export MINER_CONFIG="
chainweb:
  mining:
    coordination:
      enabled: true
      miners:
        - account: $MINER_ACCOUNT
          public-keys: [ $MINER_KEY ]
          predicate: keys-all
"
fi

if [[ -n "$ROSETTA" ]] ; then
    ROSETTA_FLAG="--rosetta"
else
    ROSETTA_FLAG="--no-rosetta"
fi

    $BIN \
    --config-file=kadena.conf \
    --config-file <(echo "$MINER_CONFIG") \
    --bootstrap-reachability=0 \
    --p2p-hostname="$CHAINWEB_P2P_HOST" \
    --p2p-port="$CHAINWEB_P2P_PORT" \
    --service-port="$CHAINWEB_SERVICE_PORT" \
    "$ROSETTA_FLAG" \
    --log-level="$LOGLEVEL" \
    +RTS -N -t -A64M -H500M -RTS \
    "$@"

{{end}}
