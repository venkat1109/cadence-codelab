export PATH := $(GOPATH)/bin:$(PATH)

.PHONY: helloworld bins


COMMON_SRC := $(shell find ./common -name "*.go")

mkbins: 
	mkdir -p bins

vendor/glide.updated: glide.yaml
	glide install
	touch vendor/glide.updated

helloworld: vendor/glide.updated mkbins $(COMMON_SRC)
	go build -i -o bins/helloworld_starter helloworld/starter.go
	go build -i -o bins/helloworld_worker helloworld/worker.go

cli: vendor/glide.updated mkbins $(COMMON_SRC)
	go build -i -o bins/cli tools/cli.go

eats: vendor/glide.updated mkbins $(COMMON_SRC)
	go build -i -o bins/eats_worker eatsapp/worker/main.go
	go build -i -o bins/eats_server eatsapp/webserver/main.go

cron: vendor/glide.updated mkbins $(COMMON_SRC)
	go build -i -o bins/cron_worker cron/worker.go
	go build -i -o bins/cron_starter cron/starter.go

bins: cli helloworld eats cron

clean:
	rm -rf bins
