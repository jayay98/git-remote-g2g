FROM golang
WORKDIR /app
COPY cmd /app/cmd
COPY pkg /app/pkg
COPY go.mod go.sum Makefile /app
RUN make install