# innuendo-1 Genesis Transation

In order to become a genesis validator on `innuendo-1`, you must generate a gentx (genesis transaction) in the following way:

For generating the gentx, any version of `quicksilverd` will suffice. A new release (v0.6.x) will be generated for genesis.


## Init a new Quicksilver instance, using: 

```quicksilverd init --chain_id 'innuendo-1' <your_moniker>```

- Create, or recover a key, using 

```quicksilverd keys add <keyname>` or `quicksilver keys add --recover <keyname>```

## Fund your account 

```quicksilverd add-genesis-account <keyname> 10000000000uqck```

## Create the gentx with 

```quicksilverd gentx <keyname> 10000000000uqck --moniker <validatorname> --chain-id innuendo-1```  

You may use `quicksilverd gentx -h` to see additional flags for setting commission / identities / etc.)

You will receive output such as: `Genesis transaction written to "~/.quicksilver/config/gentx/gentx-41eb6aa3ce902adf9603b1a1c55d005790a099c1.json"`

Copy the contents of this file into the registration form.

