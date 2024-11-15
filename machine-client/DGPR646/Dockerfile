FROM golang:1.22

RUN apt-get update && apt-get install sudo openssh-server ufw -y

RUN useradd ssh-user && echo 'ssh-user:passwd' | chpasswd

RUN ufw allow 22

WORKDIR /app

#COPY go.mod go.sum ./
#RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

EXPOSE 22

#ENTRYPOINT service ssh start && /main
ENTRYPOINT ["service", "ssh", "start", "&&", "/main"]
CMD ["./config.json"]