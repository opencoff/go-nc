
.PHONY: all clean realclean

all:
	-mkdir -p bin
	go build -o bin/gonc

realclean clean:
	rm -f bin/gonc

