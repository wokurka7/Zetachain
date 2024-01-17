
# Prerequisites 
- Have the zetacored binary in you path

# Setup a node folder
- Run `zetacored init` to create a new node folder. Note this is just a temporary node just to create the authz grants 

```shell
    zetacored init <moniker> --chain-id=zetachain_7000-1
```
- Add your operator key to keyring. Using test backend is fine as this key can be deleted after the grants are created.

```shell
    echo "<mnemonic>" | zetacored keys add operator --algo=secp256k1 --recover --keyring-backend=test
```

- Add a new hotkey to the keyring . Using test backend is fine as this key can be deleted after the grants are created.

```shell
    zetacored keys add hotkey --algo=secp256k1 --keyring-backend=test
```
NOTE : the console will output the mnemonic for the hotkey. Save this mnemonic .
    
# Create the authz grants
- Assign enviorment variables 

```shell
    export hotkey_address=$(zetacored keys show hotkey -a --keyring-backend=test)
    export operator_address=$(zetacored keys show operator -a --keyring-backend=test)
    export chain_id=zetachain_7000-1
    export node_ip=<node_ip> # the ip of the node you are broadcasting the tx to
    export min_gas_price=10000000000azeta
```

- Create grants 
```shell
zetacored tx authz grant $hotkey_address 'generic' --msg-type='/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx' --from=operator --keyring-backend=test --chain-id=$chain_id --node=$node_ip --gas=auto --gas-adjustment=1.5 --gas-prices=$min_gas_price --yes
zetacored tx authz grant $hotkey_address 'generic' --msg-type='/zetachain.zetacore.crosschain.MsgGasPriceVoter' --from=operator --keyring-backend=test --chain-id=$chain_id --node=$node_ip --gas=auto --gas-adjustment=1.5 --gas-prices=$min_gas_price --yes
zetacored tx authz grant $hotkey_address 'generic' --msg-type='/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx' --from=operator --keyring-backend=test --chain-id=$chain_id --node=$node_ip --gas=auto --gas-adjustment=1.5 --gas-prices=$min_gas_price --yes
zetacored tx authz grant $hotkey_address 'generic' --msg-type='/zetachain.zetacore.crosschain.MsgCreateTSSVoter' --from=operator --keyring-backend=test --chain-id=$chain_id --node=$node_ip --gas=auto --gas-adjustment=1.5 --gas-prices=$min_gas_price --yes
zetacored tx authz grant $hotkey_address 'generic' --msg-type='/zetachain.zetacore.crosschain.MsgAddToOutTxTracker' --from=operator --keyring-backend=test --chain-id=$chain_id --node=$node_ip --gas=auto --gas-adjustment=1.5 --gas-prices=$min_gas_price --yes
zetacored tx authz grant $hotkey_address 'generic' --msg-type='/zetachain.zetacore.observer.MsgAddBlameVote' --from=operator --keyring-backend=test --chain-id=$chain_id --node=$node_ip --gas=auto --gas-adjustment=1.5 --gas-prices=$min_gas_price --yes
zetacored tx authz grant $hotkey_address 'generic' --msg-type='/zetachain.zetacore.observer.MsgAddBlockHeader' --from=operator --keyring-backend=test --chain-id=$chain_id --node=$node_ip --gas=auto --gas-adjustment=1.5 --gas-prices=$min_gas_price --yes

```
- View grants 
```shell
zetacored q authz grants $operator_address $hotkey_address --node=$node_ip
```

- Remove the node folder once the transactions are confirmed

```shell
    rm -rf ~/.zetacored
``` 