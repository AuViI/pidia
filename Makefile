install:
	go build .
doc:
	pandoc -V geometry:margin=.8in -sS  README.md -o README.pdf
