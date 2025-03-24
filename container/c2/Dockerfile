FROM 98a2693c996c2294358552f48373498d
ARG DEBIAN_FRONTEND=noninteractive

# Install libpcap for packet capturing from the network interface
RUN apt update
RUN apt install -y libpcap-dev

# Update apt sources to use a higher version of Golang
RUN sed -i 's/focal/noble/g' /etc/apt/sources.list
RUN apt update

# Install Golang
RUN apt install -y golang

# Set GOPROXY
RUN go env -w GOPROXY=https://goproxy.cn,direct

# Copy the c2 directory
WORKDIR /root/c2
COPY c2 .

# Run go mod tidy
RUN go mod tidy

# Integrate with Seedemu
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
