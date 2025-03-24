FROM 98a2693c996c2294358552f48373498d
ARG DEBIAN_FRONTEND=noninteractive

RUN apt update
RUN apt install -y tmux

# Upgrade to a higher version of ollama
RUN sed -i 's/focal/noble/g' /etc/apt/sources.list
RUN apt update

RUN curl -fsSL https://ollama.com/install.sh -o ollama_install.sh
RUN sed -i 's|https://ollama.com/download/|https://github.com/ollama/ollama/releases/download/v0.5.7/|' ollama_install.sh
RUN chmod +x ollama_install.sh
RUN sh ollama_install.sh
RUN OLLAMA_HOST=0.0.0.0 ollama serve & \
    sleep 5 && \
    ollama pull deepseek-r1:1.5b && \
    pkill ollama

# Integrate Seedemu
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
