python.lint:
	$(PYTHON) -m black --exclude venv/ .
	venv/bin/isort -rc --skip venv .
	$(PYTHON) -m autopep8 --in-place --aggressive --aggressive --exclude './venv/*' --recursive .
	$(PYTHON) -m autoflake --exclude './venv/*' --remove-all-unused-imports --ignore-init-module-imports --remove-duplicate-keys --remove-unused-variables .
	$(PYTHON) -m flake8 --exclude venv/ --max-line-length 250 .
	$(PYTHON) -m safety check