FROM golang:1.10.3-alpine3.7

#新增 GLIBC
ENV GLIBC_VERSION 2.27-r0

# Download and install glibc
RUN apk add --update && \
    apk add --no-cache --upgrade \
    ca-certificates \
    gcc \
    g++ \
    make \
    curl \
    git

RUN curl -Lo /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub && \
    curl -Lo /var/glibc.apk "https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-${GLIBC_VERSION}.apk" && \
    curl -Lo /var/glibc-bin.apk "https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-bin-${GLIBC_VERSION}.apk" && \
    apk add /var/glibc-bin.apk /var/glibc.apk && \
    /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib && \
    echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf

#掛載 calibre 最新3.x

ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:/opt/calibre/lib
ENV PATH $PATH:/opt/calibre/bin

RUN curl -Lo /var/linux-installer.py https://download.calibre-ebook.com/linux-installer.py

#RUN mkdir -p /go/src/yy.com/ && cd /go/src/yy.com/ && git clone https://yy.com/mindoc.git && cd mindoc

ADD . /go/src/yy.com/mindoc

WORKDIR /go/src/yy.com/mindoc

RUN	 go get -u github.com/golang/dep/cmd/dep && dep ensure  && \
	CGO_ENABLE=1 go build -v -o mindoc_linux_amd64 -ldflags="-w -X main.VERSION=$TAG -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" && \
    rm -rf commands controllers models modules routers tasks vendor docs search data utils graphics .git Godeps uploads/* .gitignore .travis.yml Dockerfile gide.yaml LICENSE main.go README.md conf/enumerate.go conf/mail.go install.lock simsun.ttc

ADD start.sh /go/src/yy.com/mindoc
ADD simsun.ttc /usr/share/fonts/win/

FROM alpine:latest

LABEL maintainer="longfei6671@163.com"

RUN apk add --update && \
    apk add --no-cache --upgrade \
    ca-certificates \
    mesa-gl \
    python \
    qt5-qtbase-x11 \
    xdg-utils \
    libxrender \
    libxcomposite \
    xz \
    imagemagick \
    imagemagick-dev \
    msttcorefonts-installer \
    fontconfig && \
    update-ms-fonts && \
    fc-cache -f

COPY --from=0 /var/glibc.apk .
COPY --from=0 /var/glibc-bin.apk .
COPY --from=0 /etc/apk/keys/sgerrand.rsa.pub /etc/apk/keys/sgerrand.rsa.pub
COPY --from=0 /var/linux-installer.py .
COPY --from=0 /usr/share/fonts/win/simsun.ttc /usr/share/fonts/win/

RUN  apk add glibc-bin.apk glibc.apk && \
    /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib && \
    echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf && \
    rm -rf glibc.apk glibc-bin.apk /var/cache/apk/* && \
    chmod a+r /usr/share/fonts/win/simsun.ttc


ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:/opt/calibre/lib
ENV PATH $PATH:/opt/calibre/bin

RUN cat linux-installer.py | python -c "import sys; main=lambda x,y:sys.stderr.write('Download failed\n'); exec(sys.stdin.read()); main(install_dir='/opt', isolated=True)" && \
    rm -rf /tmp/* linux-installer.py

WORKDIR /mindoc

COPY --from=0 /go/src/yy.com/mindoc .


ENV ZONEINFO=/mindoc/lib/time/zoneinfo.zip
RUN chmod +x start.sh

CMD ["./start.sh"]