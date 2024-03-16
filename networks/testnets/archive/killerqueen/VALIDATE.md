### Validating from Genesis

If you submitted a genesis transaction (gentx) file, you have need to ensure that the `.quicksilverd/config/priv_validator_key.json` from where you created the genesis transaction is present in your full node. 

With the correct configuration and genesis file, you will see the following log message until genesis time:
```
5:54PM INF This node is a validator addr=29E8CB95478FA1478FB53D46600CD60704AD2099 module=consensus pubKey=nFWrwOkX4K0Qmc7H1oIdMUQrPS1NZqlqKZDt+CVafh0=
...
5:54PM INF Genesis time is in the future. Sleeping until then... genTime=2022-06-23T09:00:00Z
```

If you see `This node is not a validator`; either your gentx was not successfully included in the genesis file, or your `priv_validator_key.json` does not match that used to create your validator. You will be able to join the network after genesis, using the instructions below.

### Create Validator

If you did NOT submit a genesis transaction, or your genesis transaction was not included because it was invalid, once you have `quicksilverd` running and sync, you can create a validator on the Quicksilver network via a `MsgCreateValidator` transaction:

```go

$ quicksilverd tx staking create-validator \
--amount=<amount> \
--pubkey=$(quicksilverd tendermint show-validator) \
--moniker="<moniker>" \
--chain-id="killerqueen-1" \
--commission-rate="<commission>" \
--commission-max-rate="<max-commission>" \
--commission-max-change-rate="<max-commission-rate-change>" \
--min-self-delegation="<min-self-delegation>" \
--fees=<fees> \
--from=<key-name>
```
