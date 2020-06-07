# Overview
The Superman API inspects user login events for suspicious account activity
based on IP address, geographical location, and time between login events. When
a series of login events has been found to exceed a threshold speed of 500 MPH,
Superman flags the activity as potentially suspicious. Exceeding the speed threshold
indicates suspicious activity because 500MPH is about the speed of a commercial
airplane in flight. The assumption that a login attempt was by the same individual
is broken when physical distance and the time between logins indicates a travel
speed greater than in-flight travel speed.

## Build and Test

### Dependencies
* If you're on Mac, you may need to ensure openssl is on your PATH.
* docker-compose version 1.25.5
* docker version 19.03.8
* Make

### Build with Make
**Known issue**: the volume mount may mount as root, making the database
file unusable for the docker-run step.
If this is the case for you, skip to [Build with Docker](#build-with-docker)
```shell
# local binary, requires go installation, cgo dependencies
make build-api
# or on mac
make build-api-mac

# Build the docker container
make docker-build

# Test
make docker-test

# Run the container (not detached)
make docker-run
```

### Build with docker-compose
**Known issue**: the volume mount may mount as root, making the database
file unusable for the build & run step.
If this is the case for you, skip to [Build with Docker](#build-with-docker)
```shell
# force a clean build
docker-compose -f build/package/docker-compose.yaml build --force-rm

# build & run
docker-compose -f build/package/docker-compose.yaml up --build

# test
docker-compose -f build/test/docker-compose.yaml up --exit-code-from superman-api-test
```

### Build with docker
```shell
# build
docker build -f Dockerfile -t superman-api:latest .
# run with default args
docker run -p 8080:8080 superman-api:latest
# clean up
docker rm superman-api:latest

# build the test image
docker build -f Dockerfile-test -t superman-api:test .
# run tests
docker run superman-api:test make test
# clean up
docker rm superman-api:test
```


# Example

For example, if user bob was seen to log in from three different IPs all at different times:
1. IP: 206.81.252.7, Timestamp: 1514763200 (12/31/2017 11:33PM GMT)
2. IP: 206.81.252.200, Timestamp: 1514764800 (01/01/2018 12AM GMT)
3. IP: 42.222.21.19, Timestamp: 1514769200 (01/01/2018 1:13AM GMT)

When querying the second temporal event, we will see that the IP access from the
first event indicates an inhuman travel speed (even by airplane) for the user
identity, which indicates suspicious login activity.

```shell
# start the application
make docker-run
# or
# docker run -p 8080:8080 superman-api:latest

# In another terminal,
# Provide an initial login event for bob
curl -X POST -d '{"username": "bob","unix_timestamp": 1514764800,"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42","ip_address": "206.81.252.200"}' localhost:8080/v1/
# {
#    "currentGeo":{
#       "lat":38.9206,
#       "lon":-76.8787,
#       "radius":1000
#    },
#    "travelToCurrentGeoSuspicious":false,
#    "travelFromCurrentGeoSuspicious":false
# }

# Some time later, bob logs in from a different IP address
curl -X POST -d '{"username": "bob","unix_timestamp": 1514769200,"event_uuid": "8ae38b0c-a8bf-11ea-bb37-0242ac130002","ip_address": "42.222.21.19"}' localhost:8080/v1/
# {
#    "currentGeo":{
#       "lat":34.7725,
#       "lon":113.7266,
#       "radius":50
#    },
#    "travelToCurrentGeoSuspicious":true,
#    "travelFromCurrentGeoSuspicious":false,
#    "precedingIpAccess":{
#       "lat":38.9206,
#       "lon":-76.8787,
#       "radius":1000,
#       "ip":"206.81.252.200",
#       "speed":5971,
#       "timestamp":1514764800
#    }
# }

# We received this event out of order, but this login event occurred before the first curl request,
# again, from a different IP
curl -X POST -d '{"username": "bob","unix_timestamp": 1514763200,"event_uuid": "a547a38c-d23e-4990-be23-81cf212102b3","ip_address": "206.81.252.7"}' localhost:8080/v1/
# {
#    "currentGeo":{
#       "lat":38.9206,
#       "lon":-76.8787,
#       "radius":1000
#    },
#    "travelToCurrentGeoSuspicious":false,
#    "travelFromCurrentGeoSuspicious":false,
#    "subsequentIpAccess":{
#       "lat":38.9206,
#       "lon":-76.8787,
#       "radius":1000,
#       "ip":"206.81.252.200",
#       "speed":0,
#       "timestamp":1514764800
#    }
# }

# Try the first query again to see new preceding and subsequent access data
curl -X POST -d '{"username": "bob","unix_timestamp": 1514764800,"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42","ip_address": "206.81.252.200"}' localhost:8080/v1/
# {
#    "currentGeo":{
#       "lat":38.9206,
#       "lon":-76.8787,
#       "radius":1000
#    },
#    "travelToCurrentGeoSuspicious":false,
#    "travelFromCurrentGeoSuspicious":true,
#    "precedingIpAccess":{
#       "lat":38.9206,
#       "lon":-76.8787,
#       "radius":1000,
#       "ip":"206.81.252.7",
#       "speed":0,
#       "timestamp":1514763200
#    },
#    "subsequentIpAccess":{
#       "lat":34.7725,
#       "lon":113.7266,
#       "radius":50,
#       "ip":"42.222.21.19",
#       "speed":5971,
#       "timestamp":1514769200
#    }
# }

# expect: invalid ip address response
curl -X POST -d '{"username": "bob","unix_timestamp": 1514764800,"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42","ip_address": "206.81.252.432"}' localhost:8080/v1/
```

# References

- https://godoc.org/github.com/oschwald/geoip2-golang<br>
    For a ready-to-use API to interact with MaxMind DB
- https://stackoverflow.com/questions/15965166/what-is-the-maximum-length-of-latitude-and-longitude#:~:text=The%20valid%20range%20of%20latitude,of%20the%20Prime%20Meridian%2C%20respectively<br>
    For determining valid lat/long ranges
- https://gorm.io/docs/query.html<br>
    For general query building reference
- https://github.com/umahmood/haversine<br>
    For Go implementation of haversine formula
- http://www.jtrive.com/calculating-distance-between-geographic-coordinate-pairs.html<br>
    For general understanding of haversine distance and speed
- https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis<br>
    For a builder pattern building the response object
- https://yourbasic.org/golang/generate-uuid-guid/<br>
    For a random uuid generator pattern
- http://choly.ca/post/go-json-marshalling/<br>
    For data validation in JSON marshalling
- https://gist.github.com/PurpleBooth/ec81bad0a7b56ac767e0da09840f835a<br>
    For resolving a container build error "standard_init_linux.go:211: exec user process caused "no such file or directory""
- https://medium.com/@lmakarov/the-backlash-of-chmod-chown-mv-in-your-dockerfile-f12fe08c0b55<br>
    For resolving file local.db ownership issues with docker build
