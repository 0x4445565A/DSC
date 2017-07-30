FROM golang:onbuild

# Run our tests
RUN go test -v

#ADD access.log /var/log/nginx/access.log
