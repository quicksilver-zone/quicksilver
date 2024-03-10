#!/usr/bin/python3
import json

'''
This script exists to take raw exported state from quicksilver-1 chain, zero the interchainstaking and interchainquery state, and ibc connections/channels/clients.
It also burns the uqatom balance (users have already been reimbursed the underlying atoms) and ibc/ denoms (1.12 osmo, and 0.888 swth), which will be reimbursed.
This is cheaper, quicker and easier than spending 100s of osmo and going through governance to recover a channel with $1.20 locked in it.

Min commission rates are adjusted up to 5%, and the global min commission rate param is set to 5%.

Chain ID is set to quicksilver-2 and the genesis time is set to 1700 UTC on 03/01/2023.
'''

with open('export-quicksilver-1-115000.json') as file:
  input = json.load(file)
  file.close()

coins_to_burn = {}
## remove qatoms and ibc denoms
print("⚛️  Removing uqatom and ibc denoms from account balances")
balances = input.get('app_state').get('bank').get('balances')
for account_index, account in enumerate(balances):
  balance = account.get('coins', [])
  for coin_index, coin in enumerate(account.get('coins', [])):
    if coin.get('denom') != "uqck":
      coins_to_burn.update({coin.get('denom'): coins_to_burn.get(coin.get('denom'), 0)+int(coin.get('amount'))})
      print("  ⚛ Removing {} from {}".format(coin, account.get('address')))
      balance.remove(coin)
      account.update({'coins': balance})

print("⚛️  Coins to remove from supply")
print(coins_to_burn)

supply = input.get('app_state').get('bank').get('supply')
print("⚛️  Supply before migration")
print(supply)
for denom, amount in coins_to_burn.items():
  print("  ⚛ Removing {} from supply".format({"amount": str(amount), "denom": denom}))
  supply.remove({"amount": str(amount), "denom": denom})

print("⚛️  Supply post migration")
print(input.get('app_state').get('bank').get('supply'))

print("⚛️  Removing ibc channels, clients, connections and capabilities")
## remove ibc capabilities, clients, connections, channels
input.get('app_state').update({'capability': {"index": "1"}})
input.get('app_state').get('ibc').update({"channel_genesis": {"channels": [],"acknowledgements": [],"commitments": [],"receipts": [],"send_sequences": [],"recv_sequences": [],"ack_sequences": [],"next_channel_sequence": "0"}, "client_genesis": {"clients": [], "clients_consensus": [], "clients_metadata": [], "params": {"allowed_clients": ["06-solomachine","07-tendermint"]}, "create_localhost": False, "next_client_sequence": "0"}, "connection_genesis": {"connections": [],"client_connection_paths": [],"next_connection_sequence": "0","params": {"max_expected_time_per_block": "30000000000"}}})
input.get('app_state').get('transfer').update({'denom_traces': []})
input.get('app_state').get('interchainaccounts').update({"controller_genesis_state": {"active_channels": [],"interchain_accounts": [],"ports": [],"params": {"controller_enabled": True}},"host_genesis_state": {"active_channels": [],"interchain_accounts": [],"port": "icahost","params": {"host_enabled": False,"allow_messages": []}}})

## remove interchainstaking / interchain query entries
print("⚛️  Zeroing interchainstaking and interchainquery state")
input.get('app_state').get('interchainquery').update({'queries': []})
input.get('app_state').update({'interchainstaking': {'params': input.get('app_state').get('interchainstaking').get('params')}})

## reset epochs
print("⚛️  Zeroing epoch state")
input.get('app_state').get('epochs').get('epochs')[0].update({"start_time": "0001-01-01T00:00:00Z", "current_epoch": "0", "current_epoch_start_time": "0001-01-01T00:00:00Z", "epoch_counting_started": False, "current_epoch_start_height": "0"})
input.get('app_state').get('epochs').get('epochs')[1].update({"start_time": "0001-01-01T00:00:00Z", "current_epoch": "0", "current_epoch_start_time": "0001-01-01T00:00:00Z", "epoch_counting_started": False, "current_epoch_start_height": "0"})

## update min commission to 5%
print("⚛️  Updating validator min commission rate to 5%")
for validator in input.get('app_state').get('staking').get('validators'):
  if float(validator.get('commission').get('commission_rates').get('rate')) < 0.05:
    print("  ⚛ Found {} with commission rate < 0.05; updating".format(validator.get('description').get('moniker')))
    validator.get('commission').get('commission_rates').update({'rate': "0.050000000000000000"})
    if float(validator.get('commission').get('commission_rates').get('max_rate')) < 0.05:
      validator.get('commission').get('commission_rates').update({'max_rate': "0.050000000000000000"})

## update param
print("⚛️  Updating min_commission_rate param to 5%")
input.get('app_state').get('staking').get('params').update({'min_commission_rate': "0.050000000000000000"})

## chain id and genesis time
print("⚛️  Setting chain id and genesis time")
input.update({'chain_id': 'quicksilver-2', 'genesis_time': '2023-01-03T17:00:00Z'})

with open("genesis.json", "w+") as file:
  json.dump(input, file)
  file.close()
