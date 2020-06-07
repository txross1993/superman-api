FROM golang:1.14-alpine as builder


RUN apk add --update --no-cache git build-base make

WORKDIR /go/src/github.com/txross1993/superman-api

# Cache dependencies
COPY GeoLite2-City_20200602/GeoLite2-City.mmdb GeoLite2-City.mmdb
ADD vendor vendor
COPY go.mod .
COPY go.sum .
RUN go mod download

# build
COPY . .
RUN make build-api

RUN chmod 744 GeoLite2-City.mmdb && \
    mkdir -p /local-db && \
	touch /local-db/local.db


#===============================================================================

FROM scratch

ENV HOST=0.0.0.0
ENV PORT=8080
ENV GEODB=GeoLite2-City.mmdb
ENV DBPATH=/local-db

USER 1000:1000

COPY --from=builder --chown=1000:1000 /local-db /local-db
COPY --from=builder /go/src/github.com/txross1993/superman-api/GeoLite2-City.mmdb /GeoLite2-City.mmdb
COPY --from=builder /go/src/github.com/txross1993/superman-api/app /app


EXPOSE 8080

CMD ["/app"]
