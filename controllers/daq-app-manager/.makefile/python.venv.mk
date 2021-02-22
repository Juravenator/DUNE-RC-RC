venv: requirements.txt requirements-dev.txt ## set up virtualenv
	python3 -m venv venv
	source venv/bin/activate; \
	pip3 install --upgrade pip; \
	pip3 install -r requirements-dev.txt

.PHONY: upgrade
upgrade: venv ## upgrade dependencies
	source venv/bin/activate; \
	pip-upgrade requirements.txt requirements-dev.txt