FROM vineelsai/ubuntu

RUN apt update && apt upgrade -y

RUN apt install curl wget tar -y

# install docker
RUN apt install ca-certificates gnupg lsb-release -y
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
RUN echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt update
RUN apt install docker-ce docker-ce-cli containerd.io -y

WORKDIR /tmp

RUN wget https://dl.google.com/go/go1.21.3.linux-$(dpkg --print-architecture).tar.gz
RUN tar -xvf go1.21.3.linux-*.tar.gz
RUN mv go /usr/local

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

WORKDIR /usr/src/app

COPY . .

RUN go get && go build -ldflags "-s -w" && go clean -modcache

COPY etc/docker/daemon.json /etc/docker/daemon.json
COPY ./entrypoint ./entrypoint
COPY ./docker-entrypoint.d/* ./docker-entrypoint.d/

ENV DOCKER_TMPDIR=/data/docker/tmp

RUN chmod +x ./entrypoint
RUN chmod +x ./run.sh

ENTRYPOINT ["./entrypoint"]

CMD ["./run.sh"]

# docker run -it -p 3000:3000 -v "/var/run/docker.sock:/var/run/docker.sock" -v "/usr/src/app/runs:/usr/src/app/runs" vineelsai/rce
