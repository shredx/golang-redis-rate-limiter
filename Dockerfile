FROM melvinodsa/go-web-application:latest
LABEL maintainer="melvinodsa@gmail.com"

# go get the dependencies and clone the repo
RUN go get -u github.com/go-redis/redis \
	github.com/shredx/golang-redis-rate-limiter \
	&& cd $GOPATH/src/github.com/shredx/golang-redis-rate-limiter \
	&& dep ensure

EXPOSE 8085/tcp
EXPOSE 8085/udp

CMD ["revel", "run", "-v", "github.com/shredx/golang-redis-rate-limiter"]
