# xcclookup

xcclookup is a utility to assemble Quicksilver x/partiticipationrewards MsgSubmitClaim transactions for a user, in response to a REST request.

## Configuration and Installation

A configuration file is specified on the command-line and should contain the following:

```
source_chain: quicktest-1
source_lcd: https://lcd.dev.quicksilver.zone
chains:
  quickgaia-1: http://172.17.0.1:21401
  quickstar-1: http://172.17.0.1:22401
  quickosmo-1: http://172.17.0.1:23101
```

`source_chain` is the chain id of the Quicksilver chain to which submissions are made.
`source_lcd` is the LCD/API endpoint of the Quicksilver chain to which submissions are made. This is used to fetch ProtocolDatas from the Quicksilver chain.
`chains` is a map of RPC endpoints (typically port 26657) indexed by chain id, of the chains that should be checked for assets.

The config file is specified at the command line using the `-f` flag. This is presently the only command line argument.

Docker containers build from this repo are available from Docker Hub, under `quicksilverzone/xcclookup` tagged by the Git tag.

A sample docker-compose.yml file for `xcckookup` is as follows:

```
version: '3.7'
services:
  xcclookup:
    image: quicksilverzone/xcclookup:v0.3.0
    volumes:
      - /data/xcclookup:/config
    ports:
      - 8090:8090
    command:
      - /xcc
      - -f
      - /config/config.yaml
```

### Licensing and Attribution

`xcclookup` is made available under the Apache 2.0 License. 
