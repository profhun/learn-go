# learn-go
this is a tiny POC application, written in GO, which list the existing airports within a range.

## pre-requirements
- install docker

## running the code
issue this commands form the application root:

### build the docker image
`docker image build -t go-http-server .`

### run the docker image
`docker run -ti -p 8080:8080 go-http-server`

### get list
open the page http://localhost:8080/getlist and set the query params:
- lon - the longitude (float)
- lat - the lattitude (float)
- r (like radius) in kilometers

Example:
http://localhost:8080/getlist?lon=34&lat=33&r=100
