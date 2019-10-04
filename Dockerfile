FROM golang:1.12 as launcher-env
WORKDIR /usr/src/invoke
COPY . /usr/src/invoke
RUN CGO_ENABLED=0 GOOS=linux go get "github.com/gorilla/mux"
RUN CGO_ENABLED=0 GOOS=linux go build -v -o invoke

FROM alpine:3.5 as build-env

# install deps, compile, install and remove unnecessary stuff
RUN apk add --update --no-cache git coreutils build-base autoconf automake bash \
                                boost-dev zlib-dev libpng-dev jpeg-dev tiff-dev openexr-dev && \
                                git clone https://github.com/POV-Ray/povray.git && \
                                cd /povray/unix && ./prebuild.sh && \
                                cd /povray && ./configure --prefix=/opt/povray --enable-static && make -j 7 && make install && \
                                apk del build-base autoconf automake && \
                                rm -rf /var/cache/apk/* && \
                                rm -rf /povray

FROM alpine:3.5
COPY --from=launcher-env /usr/src/invoke/invoke /
COPY --from=build-env  /opt /opt
COPY test.pov /test.pov
CMD ["/invoke"] 