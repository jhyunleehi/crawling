
VERSION_INFO = "-X 'main.version=$$(git describe --tags --abbrev=0)' -X 'main.date=$$(date +%F)' -X 'main.commit=$$(git rev-parse --short HEAD)'"

MY_VER = "$$(git describe --tags --abbrev=0)"
IMG_VER = "$$(echo $(MY_VER)| sed -e s/v//g)"
IMG_NAME = docker.io/jhyunleehi/gtrend:${IMG_VER}

main:	
	go build -mod vendor --ldflags ${VERSION_INFO} main.go

# Delete files
clean:
	rm -f main
	rm -f main.exe

# Default
all: main