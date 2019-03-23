FROM melvinodsa/go-web-application:latest
LABEL maintainer="melvinodsa@gmail.com"

# go get the dependencies and clone the repo
COPY . $GOPATH/src/github.com/shredx/golang-redis-rate-limiter
RUN go get -u github.com/go-redis/redis \
	&& cd $GOPATH/src/github.com/shredx/golang-redis-rate-limiter

EXPOSE 8085/tcp
EXPOSE 8085/udp

CMD ["revel", "run", "-v", "github.com/shredx/golang-redis-rate-limiter"]
