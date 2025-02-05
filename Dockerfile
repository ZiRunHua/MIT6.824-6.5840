FROM golang:1.19

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
ENV TZ=Asia/Shanghai

WORKDIR /mit6.824

COPY test/go.mod test/go.mod
COPY test/go.sum test/go.sum
RUN cd test && go mod download

COPY src/go.mod src/go.mod
COPY src/go.sum src/go.sum
RUN cd src && go mod download

RUN cd ../

COPY . .

CMD ["tail", "-f", "/dev/null"]