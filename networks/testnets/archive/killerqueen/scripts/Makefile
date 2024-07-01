#!/usr/bin/make -f

init:
	./quicksilver.sh

clean:
	rm -rf ~/rhapsody

keys:
	./keys.sh

balances:
	./balances.sh

reset:
	./reset.sh

start:
	./start.sh
stop:
	pkill quicksilverd

logs:
	tail -f qs.log

validate:
	./validate.sh

unjail:
	./unjail.sh

show-validators:
	./show-validators.sh

check-validator-rewards:
	./check-validator-rewards.sh

check-validator-voting-power:
	./check-validator-voting-power.sh
