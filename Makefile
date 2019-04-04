all: go-bin
go-bin:
	go build -o awesome-ssh-tarpit -mod vendor .
deploy: go-bin
	scp awesome-ssh-tarpit marv:awesome-ssh-tarpit/
make clean:
	go clean