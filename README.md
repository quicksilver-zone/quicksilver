# Quicksilver

Quicksilver is the Cosmos Liquid Staking Zone.

Many existing Liquid Staking providers take power and agency _away_ from delegators, permitting only a small whitelist
of validators to delegate to, and stripping away users voting rights. Quicksilver aims to right these wrongs, with
active measures to reward and incentivise decentralisation and governance participation.

## System Requirements
This system spec has been found to be optimal:

1. Quad Core AMD or Intel (amd64) CPU; higher clock speed is preferrential to more cores, as Tendermint is largely single-threaded.
2. 32GB RAM 
3. 1TB NVMe Storage (mechanical disk are insufficient)
4. 100Mbps bidirectional Internet connection

## Software Dependencies
1. The Go programming language - https://go.dev/
2. Git distributed version control - https://git-scm.com/
3. Docker - https://www.docker.com/get-started/
4. GNU Make - https://www.gnu.org/software/make/

Make sure that the above software is installed on your system. Follow the instructions for your particular platform or use your preferred platform package manager;

In addition install `jq` (a command line JSON processor):

 - Debian based systems:  
`apt-get install jq`

 - Arch based systems:  
`pacman -S jq`

 - Mac based systems:  
`brew install jq`

 - Windows based systems (using [Chocolatey NuGet](https://chocolatey.org/)):  
`chocolatey install jq`

## Clone & Run Quicksilver (dev)

_NB!! Use a fork of the repository when you plan to create Pull Requests;_

Clone the repository from GitHub and enter the directory:

    git clone https://github.com/ingenuity-build/quicksilver.git
    cd quicksilver

Then run:

    make build-docker
    make test-docker

For subsequent tests run the following if you want to start with fresh state:

    make build-docker
    make test-docker-regen


## Origination and Attribution

Quicksilver utilises code and logic that originated from other projects; as an open-source project ourselves, 
we believe that appropriate attribution is neccessary, in order to combat plagarism.

The following modules were lifted and reused in almost entirety from Osmosis (https://github.com/osmosis-labs/osmosis), 
under the terms of the Apache 2.0 License, and we are grateful for their contribution:

```
x/mint
x/epochs
```

