

install: installbot installcli

installbot:
	go install github.com/jmichalicek/tacofancy-slack/tacobot

installcli:
	go install github.com/jmichalicek/tacofancy-slack/tacocli

tacocli:
	go build ./tacocli

tacobot:
	go build ./tacobot

runbot:
	${GOPATH}/bin/tacobot

test:
	go test ./... -v

testtacocli:
	go test ./tacocli -v

testslack:
	go test ./slack -v
