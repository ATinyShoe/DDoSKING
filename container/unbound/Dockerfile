FROM 98a2693c996c2294358552f48373498d
ARG DEBIAN_FRONTEND=noninteractive

# Update apt sources to use a higher version of unbound
RUN sed -i 's/focal/noble/g' /etc/apt/sources.list
RUN apt update

# Install unbound-1.19.2
RUN apt install -y unbound

# Configure unbound
COPY unbound.conf /etc/unbound/unbound.conf

# Start unbound
RUN /etc/init.d/unbound start

# Integrate with Seedemu
COPY 082b96ec819c95ae773daebde675ef80 /start.sh
COPY d18858afc6bb66ec3a19d872077acfd2 /seedemu_sniffer
COPY 17ac2d812a99a91e7f747e1defb72a29 /seedemu_worker
RUN chmod +x /start.sh
RUN chmod +x /seedemu_sniffer
RUN chmod +x /seedemu_worker
COPY e01e36443f9f72c6204189260d0bd276 /ifinfo.txt
COPY d3d51fdf7f4bad30dc5db560a01ce629 /interface_setup

CMD ["/start.sh"]
