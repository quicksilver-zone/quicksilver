#!/bin/bash

DEPENDENCIES="jq"

echo -en "\nChecking dependencies... "
for name in $DEPENDENCIES
do
    [[ $(type $name 2>/dev/null) ]] || { echo -en "\n    * $name is required to run this script;";deps=1; }
done
[[ $deps -ne 1 ]] && echo -e "OK\n" || { echo -e "\nInstall the missing dependencies and rerun this script...\n"; exit 1; }

set -xe

SED="sed -i"

if [[ "$OSTYPE" == "darwin"* ]]; then
    SED="sed -i ''"
    TIME="$(TZ=GMT0 date -v-2M +%Y-%m-%dT%H:%M:00Z)"
else
    TIME="$(date --date '-2 minutes' +%Y-%m-%dT%H:%M:00Z -u)"
fi

echo $SED
echo $TIME

QS_IMAGE=quicksilverzone/quicksilver
QS_VERSION=latest
TZ_IMAGE=quicksilverzone/testzone
TZ_VERSION=latest
RLY_IMAGE=quicksilverzone/relayer
RLY_VERSION=v2.1.1

CHAIN_DIR=data
CHAINID_0=qstest-1
CHAINID_1=lstest-1
CHAINID_2=lstest-2

QS1_RUN="docker-compose --ansi never run --rm -T quicksilver quicksilverd"
QS2_RUN="docker-compose --ansi never run --rm -T quicksilver2 quicksilverd"
QS3_RUN="docker-compose --ansi never run --rm -T quicksilver3 quicksilverd"
TZ1_1_RUN="docker-compose --ansi never run --rm -T testzone1-1 icad"
TZ1_2_RUN="docker-compose --ansi never run --rm -T testzone1-2 icad"
TZ1_3_RUN="docker-compose --ansi never run --rm -T testzone1-3 icad"
TZ1_4_RUN="docker-compose --ansi never run --rm -T testzone1-4 icad"
TZ2_1_RUN="docker-compose --ansi never run --rm -T testzone2-1 osmosisd"
TZ2_2_RUN="docker-compose --ansi never run --rm -T testzone2-2 osmosisd"
TZ2_3_RUN="docker-compose --ansi never run --rm -T testzone2-3 osmosisd"
TZ2_4_RUN="docker-compose --ansi never run --rm -T testzone2-4 osmosisd"
RLY_RUN="docker-compose --ansi never run --rm -T relayer rly"
HERMES_RUN="docker-compose --ansi never run --rm -T hermes hermes --config /tmp/hermes.toml"

QS1_EXEC="docker-compose --ansi never exec -T quicksilver quicksilverd"
QS2_EXEC="docker-compose --ansi never exec -T quicksilver2 quicksilverd"
QS3_EXEC="docker-compose --ansi never exec -T quicksilver3 quicksilverd"
TZ1_1_EXEC="docker-compose --ansi never exec -T testzone1-1 icad"
TZ1_2_EXEC="docker-compose --ansi never exec -T testzone1-2 icad"
TZ1_3_EXEC="docker-compose --ansi never exec -T testzone1-3 icad"
TZ1_4_EXEC="docker-compose --ansi never exec -T testzone1-4 icad"
TZ2_1_EXEC="docker-compose --ansi never exec -T testzone2-1 osmosisd"
TZ2_2_EXEC="docker-compose --ansi never exec -T testzone2-2 osmosisd"
TZ2_3_EXEC="docker-compose --ansi never exec -T testzone2-3 osmosisd"
TZ2_4_EXEC="docker-compose --ansi never exec -T testzone2-4 osmosisd"
RLY_EXEC="docker-compose --ansi never exec -T relayer"

ICQ_RUN="docker-compose --ansi never run --rm -T icq interchain-queries"
ICQ2_RUN="docker-compose --ansi never run --rm -T icq2 interchain-queries"

VAL_MNEMONIC_1="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"
VAL_MNEMONIC_2="angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice"
VAL_MNEMONIC_3="convince erupt tongue pet jeans leader boil mosquito unfair move dinosaur wrist ankle clog brown nerve next lunch speak source turtle fault gun fade"
VAL_MNEMONIC_4="cheese alarm easy kick now tattoo forward blast exercise abuse brisk race embrace cook august dwarf axis flat allow cup ripple measure keep flip"
VAL_MNEMONIC_5="ecology thank spot fork trust sorry speed april hood midnight put umbrella detail coin census crash ride fan know cup liar plastic kitten affair"
VAL_MNEMONIC_6="lock until swarm rival chaos intact style radio silent air ship siren garbage wheat runway tornado subway moral bench arrow phone medal bar feed"
VAL_MNEMONIC_7="castle quote local answer cheap crunch decrease average rare time piano income ticket weekend supply devote earth bunker exhaust network real claw require cool"
VAL_MNEMONIC_8="truth more return engage screen ramp rebuild twice core alcohol blue cactus bless text raven sure quarter north absurd among ranch text wide double"
VAL_MNEMONIC_9="parent tribe announce stuff round frame document final improve claw prison shy claim curve baby shoulder animal city thing grief ordinary twin bean unlock"
VAL_MNEMONIC_10="slab kick oven mail supply toast crisp woman erosion alpha attitude erupt pilot decrease retire version blood donor lyrics mixed same ice cotton choice"
VAL_MNEMONIC_11="onion saddle pencil fashion captain vendor firm goat mass indoor upon attend blush grocery desert inflict vocal best shell nasty barely census keep urge"
DEMO_MNEMONIC_1="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"
DEMO_MNEMONIC_2="veteran try aware erosion drink dance decade comic dawn museum release episode original list ability owner size tuition surface ceiling depth seminar capable only"
DEMO_MNEMONIC_3="snow cancel exhibit neutral cushion what bench bomb season hard mesh method virus enforce hip put voice toilet love head risk ankle toy fiscal"
DEMO_MNEMONIC_4="sustain stumble true ozone note engine unit dignity tip sheriff barrel connect fire ridge wealth echo behind will pledge coin joke mouse ripple battle"
DEMO_MNEMONIC_5="remain season shoot frog include erase august click rookie shine person oxygen pyramid table disagree language blossom island begin theory strike planet acid mad"
DEMO_MNEMONIC_6="dog remind design enrich kingdom village lottery sleep access impulse actual verb finger wreck main disorder erosion involve marriage cup quick meadow scale antenna"
DEMO_MNEMONIC_7="develop eagle toast brass table month biology fabric oven actor upper empty pigeon drum leave artist net defense excuse humor verb gown delay garden"
DEMO_MNEMONIC_8="rule casual squirrel drift mirror coast beach limb dutch tool wet small shed critic true exotic flat corn more beyond present rent mercy tomorrow"
DEMO_MNEMONIC_9="family govern swallow dignity garlic broken core expect reopen increase north recycle hair resemble dance fluid wreck paper ability various forget relax cradle rebuild"
DEMO_MNEMONIC_10="pudding another pony cancel timber brown exact valid jump glide umbrella baby joy blue next gentle pelican foster snack process unaware broom claim gentle"
DEMO_MNEMONIC_11="episode excite relax fortune devote grid staff ocean senior boss theory dutch hill laptop away fork tired bundle fan afraid parrot senior crowd indoor"
RLY_MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
RLY_MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"
RLY_MNEMONIC_3="pledge cable inform strong sadness exist favorite month illegal trial slight identify combine proud buffalo often ritual repair gown olympic bleak island worry tide"
