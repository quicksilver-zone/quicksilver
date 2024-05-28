#!/bin/sh

mkdir /icq/.icq-relayer -p
wget https://raw.githubusercontent.com/notional-labs/nmisc/main/icq-relayer/config.toml -O /icq/.icq-relayer/config.toml

exec $@
