import requests

zone_req = requests.get("https://quicksilver-2.lcd.quicksilver.zone/quicksilver/interchainstaking/v1/zones")
zone_res = zone_req.json()
for zone in zone_res.get("zones"):
  chain_id = zone.get("chain_id")
  deposit = zone.get("deposit_address").get("address")
  base_denom = zone.get("base_denom")
  balance_req = requests.get(f"https://{chain_id}.lcd.quicksilver.zone/cosmos/bank/v1beta1/balances/{deposit}")
  balance_res = balance_req.json()
  balances = balance_res.get("balances")
  for denom in balances:
    if denom.get("denom") == base_denom:
      if int(denom.get("amount", 0)) > 0:
        print(f"Balance for {chain_id} exceeds 0: {denom}")
