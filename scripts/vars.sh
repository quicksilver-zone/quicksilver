#!/bin/bash

DEPENDENCIES="jq"

echo -en "\nChecking dependencies... "
for name in $DEPENDENCIES
do
    [[ $(type $name 2>/dev/null) ]] || { echo -en "\n    * $name is required to run this script;";deps=1; }
done
[[ $deps -ne 1 ]] && echo -e "OK\n" || { echo -e "\nInstall the missing dependencies and rerun this script...\n"; exit 1; }

set -xe

QS_IMAGE=quicksilverzone/quicksilver
QS_VERSION=latest
TZ_IMAGE=quicksilverzone/testzone
TZ_VERSION=latest

CHAIN_DIR=data
CHAINID_1=test-1
CHAINID_2=test-2

QS_RUN="docker-compose --ansi never run -T quicksilver quicksilverd"
TZ_RUN="docker-compose --ansi never run -T testzone icad"
TZ2_RUN="docker-compose --ansi never run -T testzone2 icad"
TZ3_RUN="docker-compose --ansi never run -T testzone3 icad"
TZ4_RUN="docker-compose --ansi never run -T testzone4 icad"

QS_EXEC="docker-compose --ansi never exec -T quicksilver quicksilverd"
TZ_EXEC="docker-compose --ansi never exec -T testzone icad"
TZ2_EXEC="docker-compose --ansi never exec -T testzone2 icad"
TZ3_EXEC="docker-compose --ansi never exec -T testzone3 icad"
TZ4_EXEC="docker-compose --ansi never exec -T testzone4 icad"

ICQ_RUN="docker-compose --ansi never run -T icq interchain-queries"

VAL_MNEMONIC_1="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"
VAL_MNEMONIC_2="angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice"
VAL_MNEMONIC_3="convince erupt tongue pet jeans leader boil mosquito unfair move dinosaur wrist ankle clog brown nerve next lunch speak source turtle fault gun fade"
VAL_MNEMONIC_4="cheese alarm easy kick now tattoo forward blast exercise abuse brisk race embrace cook august dwarf axis flat allow cup ripple measure keep flip"
VAL_MNEMONIC_5="ecology thank spot fork trust sorry speed april hood midnight put umbrella detail coin census crash ride fan know cup liar plastic kitten affair"
DEMO_MNEMONIC_1="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"
DEMO_MNEMONIC_2="veteran try aware erosion drink dance decade comic dawn museum release episode original list ability owner size tuition surface ceiling depth seminar capable only"
DEMO_MNEMONIC_3="snow cancel exhibit neutral cushion what bench bomb season hard mesh method virus enforce hip put voice toilet love head risk ankle toy fiscal"
DEMO_MNEMONIC_4="sustain stumble true ozone note engine unit dignity tip sheriff barrel connect fire ridge wealth echo behind will pledge coin joke mouse ripple battle"
DEMO_MNEMONIC_5="remain season shoot frog include erase august click rookie shine person oxygen pyramid table disagree language blossom island begin theory strike planet acid mad"
RLY_MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
RLY_MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

