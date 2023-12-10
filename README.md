# Quicksilver

| [![codecov](https://codecov.io/gh/ingenuity-build/quicksilver/branch/develop/graph/badge.svg)](https://codecov.io/gh/ingenuity-build/quicksilver) | [![Go Report Card](https://goreportcard.com/badge/github.com/quicksilver-zone/quicksilver)](https://goreportcard.com/report/github.com/quicksilver-zone/quicksilver) | [![license](https://img.shields.io/github/license/ingenuity-build/quicksilver.svg)](https://github.com/quicksilver-zone/quicksilver/blob/main/LICENSE) | [![GolangCI](https://golangci.com/badges/github.com/quicksilver-zone/quicksilver.svg)](https://golangci.com/r/github.com/quicksilver-zone/quicksilver) | [![Discord](https://badgen.net/badge/icon/discord?icon=discord&label)](https://discord.gg/quicksilverprotocol) |  
|---------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------|

## Introduction
Quicksilver is the Cosmos Liquid Staking Zone.

Many existing Liquid Staking providers take power and agency _away_ from delegators, permitting only a small whitelist
of validators to delegate to, and stripping away users voting rights. Quicksilver aims to right these wrongs, with
active measures to reward and incentivise decentralisation and governance participation.

## Table of Contents

- [Quicksilver App](https://app.quicksilver.zone)
- [Project Documentation](https://docs.quicksilver.zone)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Contributing](CONTRIBUTING.md)
- [Security/Bug Reporting](SECURITY.md)
- [Changelog](CHANGELOG.md)

## System Requirements
This system spec has been found to be optimal:

1. Quad Core AMD or Intel (amd64) CPU; higher clock speed is preferential to more cores, as Tendermint is largely single-threaded.
2. 32GB RAM 
3. 1TB NVMe Storage (mechanical disk are insufficient)
4. 100Mbps bidirectional Internet connection

## Software Dependencies
1. The Go programming language - <https://go.dev/>
2. Git distributed version control - <https://git-scm.com/>
3. Docker - <https://www.docker.com/get-started/>
4. GNU Make - <https://www.gnu.org/software/make/>

Make sure that the above software is installed on your system. Follow the instructions for your particular platform or use your preferred platform package manager;

In addition, install `jq` (a command line JSON processor):

 - Debian based systems:
`apt-get install jq`

 - Arch based systems:
`pacman -S jq`

 - Mac based systems:
`brew install jq`

## Clone & Run Quicksilver (dev)

_NB!! Use a fork of the repository when you plan to create Pull Requests;_

Clone the repository from GitHub and enter the directory:

    git clone https://github.com/quicksilver-zone/quicksilver.git
    cd quicksilver

Then run:

    make build-docker
    make test-docker

For subsequent tests run the following if you want to start with fresh state:

    make build-docker
    make test-docker-regen


## Origination and Attribution

Quicksilver utilises code and logic that originated from other projects; as an open-source project ourselves, we believe that appropriate attribution is necessary, in order to combat plagiarism.

The following modules and packages were lifted and reused in almost entirety from Osmosis (<https://github.com/osmosis-labs/osmosis>), under the terms of the Apache 2.0 License, and we are grateful for their contribution:

    x/mint
    x/epochs
    x/tokenfactory
    test/e2e

We're also using CosmWasm, developed over the course of years with lead from Confio and support from the whole of Cosmos. 

