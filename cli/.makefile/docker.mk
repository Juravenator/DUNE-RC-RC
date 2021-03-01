##
### Docker commands
##

docker.bash: ## open a shell in a builder container
	docker build -t run-control-cli-builder - < docker/builder.Dockerfile
	docker run -it --mount type=bind,source=$(shell pwd),target=/mnt --workdir /mnt --rm --entrypoint "" run-control-cli-builder bash

docker.%: ## run any make target in a docker container
	docker build -t run-control-cli-builder - < docker/builder.Dockerfile
	docker run --mount type=bind,source=$(shell pwd),target=/mnt --workdir /mnt --rm --entrypoint "" run-control-cli-builder make "$(subst docker.,,$@)"