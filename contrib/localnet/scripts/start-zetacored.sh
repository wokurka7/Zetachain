#!/bin/bash

# This script is used to start the zetacored nodes
# It initializes the nodes and creates the genesis.json file
# It also starts the nodes
# The number of nodes is passed as an first argument to the script
# The second argument is optional and can have the following value:
# 1. upgrade : This is used to test the upgrade process, a proposal is created for the upgrade and the nodes are started using cosmovisor

/usr/sbin/sshd

if [ $# -lt 1 ]
then
  echo "Usage: genesis.sh <num of nodes> [option]"
  exit 1
fi
NUMOFNODES=$1
OPTION=$2

# create keys
CHAINID="athens_101-1"
KEYRING="test"
HOSTNAME=$(hostname)
INDEX=${HOSTNAME:0-1}

# Environment variables used for upgrade testing
export DAEMON_HOME=$HOME/.zetacored
export DAEMON_NAME=zetacored
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export DAEMON_RESTART_AFTER_UPGRADE=true
export CLIENT_DAEMON_NAME=zetaclientd
export CLIENT_DAEMON_ARGS="-enable-chains,GOERLI,-val,operator"
export DAEMON_DATA_BACKUP_DIR=$DAEMON_HOME
export CLIENT_SKIP_UPGRADE=true
export CLIENT_START_PROCESS=false
export UNSAFE_SKIP_BACKUP=true
export UpgradeName=${NEW_VERSION}

# upgrade name used for upgrade testing
export UpgradeName=${NEW_VERSION}

# generate node list
START=1
# shellcheck disable=SC2100
END=$((NUMOFNODES - 1))
NODELIST=()
for i in $(eval echo "{$START..$END}")
do
  NODELIST+=("zetacore$i")
done

echo "HOSTNAME: $HOSTNAME"

# Init a new node to generate genesis file .
# Copy config files from existing folders which get copied via Docker Copy when building images
mkdir -p ~/.backup/config
zetacored init Zetanode-Localnet --chain-id=$CHAINID
rm -rf ~/.zetacored/config/app.toml
rm -rf ~/.zetacored/config/client.toml
rm -rf ~/.zetacored/config/config.toml
cp -r ~/zetacored/common/app.toml ~/.zetacored/config/
cp -r ~/zetacored/common/client.toml ~/.zetacored/config/
cp -r ~/zetacored/common/config.toml ~/.zetacored/config/
sed -i -e "/moniker =/s/=.*/= \"$HOSTNAME\"/" "$HOME"/.zetacored/config/config.toml

# Add two new keys for operator and hotkey and create the required json structure for os_info
source ~/add-keys.sh

# Pause other nodes so that the primary can node can do the genesis creation
if [ $HOSTNAME != "zetacore0" ]
then
  echo "Waiting for zetacore0 to create genesis.json"
  sleep 10
  echo "genesis.json created"
fi

# Genesis creation following steps
# 1. Accumulate all the os_info files from other nodes on zetcacore0 and create a genesis.json
# 2. Add the observers , authorizations and required params to the genesis.json
# 3. Copy the genesis.json to all the nodes .And use it to create a gentx for every node
# 4. Collect all the gentx files in zetacore0 and create the final genesis.json
# 5. Copy the final genesis.json to all the nodes and start the nodes
# 6. Update Config in zetacore0 so that it has the correct persistent peer list
# 7. Start the nodes

# Start of genesis creation . This is done only on zetacore0
if [ $HOSTNAME == "zetacore0" ]
then
  # Misc : Copying the keyring to the client nodes so that they can sign the transactions
  ssh zetaclient0 mkdir -p ~/.zetacored/keyring-test/
  scp ~/.zetacored/keyring-test/* zetaclient0:~/.zetacored/keyring-test/
  ssh zetaclient0 mkdir -p ~/.zetacored/keyring-file/
  scp ~/.zetacored/keyring-file/* zetaclient0:~/.zetacored/keyring-file/

# 1. Accumulate all the os_info files from other nodes on zetcacore0 and create a genesis.json
  for NODE in "${NODELIST[@]}"; do
    INDEX=${NODE:0-1}
    ssh zetaclient"$INDEX" mkdir -p ~/.zetacored/
    scp "$NODE":~/.zetacored/os_info/os.json ~/.zetacored/os_info/os_z"$INDEX".json
    scp ~/.zetacored/os_info/os_z"$INDEX".json zetaclient"$INDEX":~/.zetacored/os.json
  done

  ssh zetaclient0 mkdir -p ~/.zetacored/
  scp ~/.zetacored/os_info/os.json zetaclient0:/root/.zetacored/os.json

# 2. Add the observers, authorizations, required params and accounts to the genesis.json
  zetacored collect-observer-info
  zetacored add-observer-list --keygen-block 55
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="500000000"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="100s"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["feemarket"]["params"]["min_gas_price"]="10000000000.0000"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json

# set admin account
  zetacored add-genesis-account zeta1srsq755t654agc0grpxj4y3w0znktrpr9tcdgk 100000000000000000000000000azeta
  zetacored add-genesis-account zeta1n0rn6sne54hv7w2uu93fl48ncyqz97d3kty6sh 100000000000000000000000000azeta # Funds the localnet_gov_admin account
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["authority"]["policies"]["items"][0]["address"]="zeta1srsq755t654agc0grpxj4y3w0znktrpr9tcdgk"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["authority"]["policies"]["items"][1]["address"]="zeta1srsq755t654agc0grpxj4y3w0znktrpr9tcdgk"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["authority"]["policies"]["items"][2]["address"]="zeta1srsq755t654agc0grpxj4y3w0znktrpr9tcdgk"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json

# give balance to runner accounts to deploy contracts directly on zEVM
# deployer
  zetacored add-genesis-account zeta1uhznv7uzyjq84s3q056suc8pkme85lkvhrz3dd 100000000000000000000000000azeta
# erc20 tester
  zetacored add-genesis-account zeta1datate7xmwm4uk032f9rmcu0cwy7ch7kg6y6zv 100000000000000000000000000azeta
# zeta tester
  zetacored add-genesis-account zeta1tnp0hvsq4y5mxuhrq9h3jfwulxywpq0ads0rer 100000000000000000000000000azeta
# bitcoin tester
  zetacored add-genesis-account zeta19q7czqysah6qg0n4y3l2a08gfzqxydla492v80 100000000000000000000000000azeta
# ethers tester
  zetacored add-genesis-account zeta134rakuus43xn63yucgxhn88ywj8ewcv6ezn2ga 100000000000000000000000000azeta

# 3. Copy the genesis.json to all the nodes .And use it to create a gentx for every node
  zetacored gentx operator 1000000000000000000000azeta --chain-id=$CHAINID --keyring-backend=$KEYRING --gas-prices 20000000000azeta
  # Copy host gentx to other nodes
  for NODE in "${NODELIST[@]}"; do
    ssh $NODE mkdir -p ~/.zetacored/config/gentx/peer/
    scp ~/.zetacored/config/gentx/* $NODE:~/.zetacored/config/gentx/peer/
  done
  # Create gentx files on other nodes and copy them to host node
  mkdir ~/.zetacored/config/gentx/z2gentx
  for NODE in "${NODELIST[@]}"; do
      ssh $NODE rm -rf ~/.zetacored/genesis.json
      scp ~/.zetacored/config/genesis.json $NODE:~/.zetacored/config/genesis.json
      ssh $NODE zetacored gentx operator 1000000000000000000000azeta --chain-id=$CHAINID --keyring-backend=$KEYRING
      scp $NODE:~/.zetacored/config/gentx/* ~/.zetacored/config/gentx/
      scp $NODE:~/.zetacored/config/gentx/* ~/.zetacored/config/gentx/z2gentx/
  done

# 4. Collect all the gentx files in zetacore0 and create the final genesis.json
  zetacored collect-gentxs
  zetacored validate-genesis
# 5. Copy the final genesis.json to all the nodes
  for NODE in "${NODELIST[@]}"; do
      ssh $NODE rm -rf ~/.zetacored/genesis.json
      scp ~/.zetacored/config/genesis.json $NODE:~/.zetacored/config/genesis.json
  done
# 6. Update Config in zetacore0 so that it has the correct persistent peer list
   pp=$(cat $HOME/.zetacored/config/gentx/z2gentx/*.json | jq '.body.memo' )
   pps=${pp:1:58}
   sed -i -e 's/^persistent_peers =.*/persistent_peers = "'$pps'"/' "$HOME"/.zetacored/config/config.toml
fi
# End of genesis creation steps . The steps below are common to all the nodes

# Update persistent peers
if [ $HOSTNAME != "zetacore0" ]
then
  # Misc : Copying the keyring to the client nodes so that they can sign the transactions
  ssh zetaclient"$INDEX" mkdir -p ~/.zetacored/keyring-test/
  scp ~/.zetacored/keyring-test/* "zetaclient$INDEX":~/.zetacored/keyring-test/
  ssh zetaclient"$INDEX" mkdir -p ~/.zetacored/keyring-file/
  scp ~/.zetacored/keyring-file/* "zetaclient$INDEX":~/.zetacored/keyring-file/

  pp=$(cat $HOME/.zetacored/config/gentx/peer/*.json | jq '.body.memo' )
  pps=${pp:1:58}
  sed -i -e "/persistent_peers =/s/=.*/= \"$pps\"/" "$HOME"/.zetacored/config/config.toml
fi

# 7 Start the nodes
# If upgrade option is passed, use cosmovisor to start the nodes and create a governance proposal for upgrade
if [ "$OPTION" != "upgrade" ]; then

  exec zetacored start --pruning=nothing --minimum-gas-prices=0.0001azeta --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable --home /root/.zetacored

else

  # Setup cosmovisor
  mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
  mkdir -p $DAEMON_HOME/cosmovisor/upgrades/"$UpgradeName"/bin

  # Genesis
  cp $GOPATH/bin/old/zetacored $DAEMON_HOME/cosmovisor/genesis/bin
  cp $GOPATH/bin/zetaclientd $DAEMON_HOME/cosmovisor/genesis/bin

  #Upgrades
  cp $GOPATH/bin/new/zetacored $DAEMON_HOME/cosmovisor/upgrades/$UpgradeName/bin/

  #Permissions
  chmod +x $DAEMON_HOME/cosmovisor/genesis/bin/zetacored
  chmod +x $DAEMON_HOME/cosmovisor/genesis/bin/zetaclientd
  chmod +x $DAEMON_HOME/cosmovisor/upgrades/$UpgradeName/bin/zetacored

  # Start the node using cosmovisor
  cosmovisor run start --pruning=nothing --minimum-gas-prices=0.0001azeta --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable --home /root/.zetacored >> zetanode.log 2>&1  &
  sleep 20
  echo

  # Fetch the height of the upgrade, default is 225, if arg3 is passed, use that value
  UPGRADE_HEIGHT=${3:-225}

  # If this is the first node, create a governance proposal for upgrade
  if [ $HOSTNAME = "zetacore0" ]
  then
    /root/.zetacored/cosmovisor/genesis/bin/zetacored tx gov submit-legacy-proposal software-upgrade $UpgradeName --from operator --deposit 100000000azeta --upgrade-height "$UPGRADE_HEIGHT" --title $UpgradeName --description $UpgradeName --keyring-backend test --chain-id $CHAINID --yes --no-validate --fees=2000000000000000azeta --broadcast-mode block
  fi

  # Wait for the proposal to be voted on
  sleep 8
  /root/.zetacored/cosmovisor/genesis/bin/zetacored tx gov vote 1 yes --from operator --keyring-backend test --chain-id $CHAINID --yes --fees=2000000000000000azeta --broadcast-mode block
  sleep 7
  /root/.zetacored/cosmovisor/genesis/bin/zetacored query gov proposal 1

  # We use tail -f to keep the container running
  tail -f zetanode.log

fi

