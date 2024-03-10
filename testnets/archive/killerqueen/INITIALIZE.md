# Initialize

To initialize your node, 

```go
# Replace moniker with your desired node's moniker.

$ quicksilverd init <moniker> --chain-id killerqueen-1 --home $HOME/.quicksilverd
```

This will create and initialize a .quicksilverd directory in your `$HOME` path. This directory contains all the configuration files needed to run your node, as well as a default `genesis.json` file. You can change the location of the directory by specifying a different `--home` flag.

Once initialized, overwrite the default `genesis.json` file with genesis state file for the particular network that you are joining. You may retrieve the genesis state file from the Quicksilver repository or another trusted source:

```go
$ cd ~/.quicksilverd/config

$ wget https://raw.githubusercontent.com/ingenuity-build/testnets/main/killerqueen/genesis.json

## verify the downloaded file matches the following hash:

shasum -a 256 genesis.json 
3510dd3310e3a127507a513b3e9c8b24147f549bac013a5130df4b704f1bac75  genesis.json

## alternatively, sort keys (-s) and remove whitespace (-C), and verify:

jq . -sC genesis.json | shasum -a 256
32686cb333bada56f0e2e101137654dd11d00f40a0c5ae265e589795b032a4f1  -
```
