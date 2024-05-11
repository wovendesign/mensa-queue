############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder 

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create appuser.
ENV USER=appuser
ENV UID=10001 

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
	--disabled-password \    
	--gecos "" \    
	--home "/nonexistent" \    
	--shell "/sbin/nologin" \    
	--no-create-home \    
	--uid "${UID}" \    
	"${USER}"

WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

# Fetch dependencies.

# Using go get.
RUN go get -d -v
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/mensaqueue

############################
# STEP 2 build a small image
############################
FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the SSL root certificates from the builder image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable.
COPY --from=builder /go/bin/mensaqueue /go/bin/mensaqueue

# Use an unprivileged user.
USER appuser:appuser

# Run the mensaqueue binary.
ENTRYPOINT ["/go/bin/mensaqueue"]