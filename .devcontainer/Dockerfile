FROM golang:1.18

RUN apt-get update
RUN apt install -y libthai0
RUN apt install -y xfonts-thai
RUN apt install -y wkhtmltopdf
RUN apt-get install -y tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata
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

   
