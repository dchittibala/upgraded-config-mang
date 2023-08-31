GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
GOGET=${GOCMD} get
BINARY_NAME=upgraded-config-mang

binary: clean builder

clean:
	${GOCLEAN}
	rm -rf ${BINARY_NAME}
	rm -rf ./cmd/${BINARY_NAME}/${BINARY_NAME}

builder:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GOBUILD} -o ${BINARY_NAME}-linux -v
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin ${GOBUILD} -o ${BINARY_NAME}-darwin -v