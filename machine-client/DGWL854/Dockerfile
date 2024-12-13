FROM golang:1.22

RUN apt-get update && apt-get install sudo openssh-server ufw -y
RUN useradd ssh-user && echo 'ssh-user:passwd' | chpasswd
RUN ufw allow 22

WORKDIR /app

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /main
RUN chmod +x ./entrypoint.sh && chmod u+w ./data

EXPOSE 22

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./config.json"]