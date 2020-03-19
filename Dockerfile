FROM golang:1.13-alpine AS builder

# installing usual dependencies
RUN apk add bash ca-certificates git gcc g++ libc-dev curl openssh-client

# setting up builder environment
WORKDIR /src/tetest
COPY . .

# generating a "throwaway" key
# NOTE: such approach should be avoided in real scenarios
RUN mkdir ~/.ssh && ssh-keygen -t rsa -f ~/.ssh/id_rsa

# recognizing hostnames
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts

# configuring git
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

# building
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /bin/tetest

# second stage; producing the final image
FROM alpine
RUN apk --no-cache add bash ca-certificates
COPY --from=builder /bin/tetest /bin/tetest
EXPOSE 8080
ENTRYPOINT ["/bin/tetest", "start"]
