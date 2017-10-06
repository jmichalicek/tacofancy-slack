
install:
	go install github.com/jmichalicek/tacofancy-slack/tacobot
installcli:
	go install github.com/jmichalicek/tacofancy-slack/tacocli
runbot:
	${GOPATH}/bin/tacobot
