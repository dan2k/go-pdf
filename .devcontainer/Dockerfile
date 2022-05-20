FROM golang:1.18

RUN apt-get update && apt-get install -y openconnect iptables git build-essential libtool autoconf automake gcc libxml2 vim autoconf gcc libxml2 pkg-config gettext liboath0 libopenconnect3 libproxy1 libstoken1 libtomcrypt0 libtommath0 libxml2-dev p11-kit libp11-2 libproxy1 trousers libstoken1 libpcsclite1 libssl-dev libxml2-dev iputils-ping
RUN git clone https://github.com/dlenski/openconnect.git && cd openconnect && ./autogen.sh && ./configure --with-vpnc-script=/usr/share/vpnc-scripts/vpnc-script --without-openssl-version-check && make

COPY scrics/connect.sh /root
RUN chmod +x /root/connect.sh


CMD ["/root/connect.sh"]

ENV GO111MODULE on
WORKDIR /var/www

RUN go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
  go install github.com/ramya-rao-a/go-outline@latest && \
  go install github.com/nsf/gocode@latest && \
  go install github.com/acroca/go-symbols@latest && \
  go install github.com/fatih/gomodifytags@latest && \
  go install github.com/josharian/impl@latest && \
  go install github.com/haya14busa/goplay/cmd/goplay@latest && \
  go install github.com/go-delve/delve/cmd/dlv@latest && \
  go install golang.org/x/lint/golint@latest && \
  go install golang.org/x/tools/gopls@latest && \
  go install github.com/cosmtrek/air@latest  

   
