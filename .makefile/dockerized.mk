##
### Docker commands
##

docker.images: ## build docker images used for local running
	cd docker && docker build -t dune-rc-march-ssh-server -f ssh-server.Dockerfile .
	cd docker && docker build -t dune-rc-march-ansible -f ansible.Dockerfile .

docker.ansible: ## run ansible on docker images
	docker-compose -f docker/docker-compose.yml run --rm ansible -i /mnt/ansible/hosts.yaml /mnt/ansible/playbook.yaml

docker.start: ## start Run Control setup in docker containers
	docker-compose -f docker/docker-compose.yml up

docker.bash: ## open a shell in a runner container
	docker-compose -f docker/docker-compose.yml run --rm --entrypoint "/bin/bash -c" ansible bash

docker.%: ## run any make target in a docker container
	docker-compose -f .makefile/docker-compose.yml run runner "make $(subst docker.,,$@)"