FROM 98a2693c996c2294358552f48373498d
ARG DEBIAN_FRONTEND=noninteractive

# 下载libpcap用于从网卡发包
RUN apt update
RUN apt install -y libpcap-dev
RUN apt install -y net-tools


# 更新apt源，需要使用高版本的golang
RUN sed -i 's/focal/noble/g' /etc/apt/sources.list
RUN apt update

# 安装Golang
RUN apt install -y golang

# 设置GOPROXY
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制Bot目录
WORKDIR /root/bot
COPY bot .

# 将服务器IP文件复制进去
COPY c2.txt /root/bot/serverfile/c2.txt
COPY reflector.txt /root/bot/serverfile/reflector.txt
COPY resolver.txt /root/bot/serverfile/resolver.txt

# 运行go mod tidy
RUN go mod tidy

# 接入Seedemu
WORKDIR /
COPY 082b96ec819c95ae773daebde675ef80 /start.sh
COPY d18858afc6bb66ec3a19d872077acfd2 /seedemu_sniffer
COPY 17ac2d812a99a91e7f747e1defb72a29 /seedemu_worker
RUN chmod +x /start.sh
RUN chmod +x /seedemu_sniffer
RUN chmod +x /seedemu_worker
COPY e01e36443f9f72c6204189260d0bd276 /ifinfo.txt
COPY d3d51fdf7f4bad30dc5db560a01ce629 /interface_setup

CMD ["/start.sh"]

