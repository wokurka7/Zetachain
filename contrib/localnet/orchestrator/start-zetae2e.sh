#!/bin/bash

# The script run the zetae2e CLI to run local end-to-end tests
# First argument is the command to run the local e2e
# A second optional argument can be passed and can have the following value:
# upgrade: run the local e2e once, then restart zetaclientd at upgrade height and run the local e2e again

ZETAE2E_CMD=$1
OPTION=$2

echo "waiting for geth RPC to start..."
sleep 2

### Create the accounts and fund them with Ether on local Ethereum network

# unlock the deployer account
echo "funding deployer address 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC with 10000 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", value: web3.toWei(10000,"ether")})' attach http://eth:8545

# unlock erc20 tester accounts
echo "funding deployer address 0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6 with 10000 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6", value: web3.toWei(10000,"ether")})' attach http://eth:8545

# unlock zeta tester accounts
echo "funding deployer address 0x5cC2fBb200A929B372e3016F1925DcF988E081fd with 10000 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x5cC2fBb200A929B372e3016F1925DcF988E081fd", value: web3.toWei(10000,"ether")})' attach http://eth:8545

# unlock bitcoin tester accounts
echo "funding deployer address 0x283d810090EdF4043E75247eAeBcE848806237fD with 10000 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x283d810090EdF4043E75247eAeBcE848806237fD", value: web3.toWei(10000,"ether")})' attach http://eth:8545

# unlock ethers tester accounts
echo "funding deployer address 0x8D47Db7390AC4D3D449Cc20D799ce4748F97619A with 10000 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x8D47Db7390AC4D3D449Cc20D799ce4748F97619A", value: web3.toWei(10000,"ether")})' attach http://eth:8545

# unlock miscellaneous tests accounts
echo "funding deployer address 0x90126d02E41c9eB2a10cfc43aAb3BD3460523Cdf with 10000 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x90126d02E41c9eB2a10cfc43aAb3BD3460523Cdf", value: web3.toWei(10000,"ether")})' attach http://eth:8545

# unlock admin erc20 tests accounts
echo "funding deployer address 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 with 10000 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", value: web3.toWei(10000,"ether")})' attach http://eth:8545

# unlock the TSS account
echo "funding TSS address 0xF421292cb0d3c97b90EEEADfcD660B893592c6A2 with 10000 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xF421292cb0d3c97b90EEEADfcD660B893592c6A2", value: web3.toWei(10000,"ether")})' attach http://eth:8545

### Run zetae2e command depending on the option passed

if [ "$OPTION" == "upgrade" ]; then

  # Run the e2e tests, then restart zetaclientd at upgrade height and run the e2e tests again

  # Fetch the height of the upgrade, default is 225, if arg3 is passed, use that value
  UPGRADE_HEIGHT=${3:-225}

  # Run zetae2e, if the upgrade height is lower than 100, we use the setup-only flag
  if [ "$UPGRADE_HEIGHT" -lt 100 ]; then
    echo "running E2E command to setup the networks..."
    zetae2e "$ZETAE2E_CMD" --setup-only --config-out deployed.yml
  else
    echo "running E2E command to setup the networks and populate the state..."

    # Use light flag to ensure tests can complete before the upgrade height
    zetae2e "$ZETAE2E_CMD" --config-out deployed.yml --light
 fi
  ZETAE2E_EXIT_CODE=$?

  if [ $ZETAE2E_EXIT_CODE -ne 0 ]; then
    echo "E2E setup failed"
    exit 1
  fi

  echo "E2E setup passed, waiting for upgrade height..."

  # Restart zetaclients at upgrade height
  /work/restart-zetaclientd-at-upgrade.sh -u "$UPGRADE_HEIGHT" -n 2

  echo "waiting 10 seconds for node to restart..."

  sleep 10

  echo "running E2E command to test the network after upgrade..."

  # Run zetae2e again
  # When the upgrade height is greater than 100 for upgrade test, the Bitcoin tests have been run once, therefore the Bitcoin wallet is already set up
  # Use light flag to skip advanced tests
  if [ "$UPGRADE_HEIGHT" -lt 100 ]; then
    zetae2e "$ZETAE2E_CMD" --skip-setup --config deployed.yml --light
  else
    zetae2e "$ZETAE2E_CMD" --skip-setup --config deployed.yml --skip-bitcoin-setup --light
  fi

  ZETAE2E_EXIT_CODE=$?
  if [ $ZETAE2E_EXIT_CODE -eq 0 ]; then
    echo "E2E passed after upgrade"
    exit 0
  else
    echo "E2E failed after upgrade"
    exit 1
  fi

else

  # Run the e2e tests normally

  echo "running e2e tests..."

  eval "zetae2e $ZETAE2E_CMD"
  ZETAE2E_EXIT_CODE=$?

  # if e2e passed, exit with 0, otherwise exit with 1
  if [ $ZETAE2E_EXIT_CODE -eq 0 ]; then
    echo "e2e passed"
    exit 0
  else
    echo "e2e failed"
    exit 1
  fi

fi
