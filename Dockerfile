FROM melvinodsa/go-web-application:latest
LABEL maintainer="melvinodsa@gmail.com"

WORKDIR $GOPATH/src/github.com/shredx/golang-redis-rate-limiter

# go get the dependencies and clone the repo
COPY . $GOPATH/src/github.com/shredx/golang-redis-rate-limiter
RUN cd $GOPATH/src/github.com/shredx/golang-redis-rate-limiter \
	&& dep ensure

EXPOSE 8085/tcp
EXPOSE 8085/udp

ENTRYPOINT ["revel", "run", "-v", "github.com/shredx/golang-redis-rate-limiter"]