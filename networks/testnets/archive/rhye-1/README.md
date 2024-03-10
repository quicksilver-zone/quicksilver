# Rhye-1 chain sync-up instructions 

 We will use `quicksilverd v1.4.2-rc7` to syncup the node from a snapshot and join the network. 
## 1. Build Instructions
 ```git fetch && git checkout v1.4.2-rc7```

 ```make install```


## 2. Init a new Quicksilver instance, using: 

```quicksilverd init --chain_id 'rhye-1' <your_moniker>```

## 3. Download the genesis.json and place in /config

``` wget -O ~/.quicksilverd/config/genesis.json https://raw.githubusercontent.com/ingenuity-build/testnets/main/rhye/genesis.json ```

 (be sure to replace ~/.quicksilverd with your node's HOME).

 ## 4. Download the snapshot and replace data folder
 
 ```wget -O data-rhye-1.tar.gz https://storage.googleapis.com/rhye-1_snapshots/rhye-1.tar.gz```

 ``` tar -xvzf data-rhye-1.tar.gz -C ~/.quicksilverd/```

 (be sure to replace ~/.quicksilverd with your node's HOME).

 ## 5. Add peers in config 
 
 Add peers from [PEERS.md](PEERS.md) in ```HOME/config/config.toml``` in ```[p2p]``` section

## 6. Start the node

```quicksilverd start ``` or, if using systemd, ```systemctl start quicksilver```


# For Validators:


**NOTE**: Faucet for ```rhye-1``` is available at https://faucet.test.quicksilver.zone

## Rhye-1 Instructions

Rhye-1  is a replacement for the innuendo-5 long term testnet for Quicksilver; it was created with the intent of onboarding public testnets against a fresh chain. It starts with v1.4.2-rc0 of the `quicksilverd` binary which can be downloaded from the releases section of this repository, or as a docker container from quicksilverzone/quicksilver:v1.4.2-rc0.

If you sync from genesis, you will need to ensure you follow the upgrade paths.



The network was started asyncronously. Accounts and validators from previous testnets were dropped. You will need to request tokens from faucet. Genesis can be downloaded from [here](genesis.json).

You may download a recent snapshot of innuendo-4 from here: https://storage.googleapis.com/rhye-1_snapshots/rhye-1.tar.gz . This snapshot is using v1.4.2-rc7 of quicksilverd.

The current peers are available: 
```
8e14e58b054248a04be96e4a40d6359e93b636ac@65.108.65.94:26656,5a3c424c19d9ab694190a7805a2b1a146460d752@65.108.2.27:26656,e6bf55bc9f08958b7518bea455423375db78d1ef@65.108.13.176:26656
```
