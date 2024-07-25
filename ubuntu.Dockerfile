FROM ubuntu:24.04

RUN apt update && apt install -y curl dnsutils iputils-ping iptables

CMD ["sh", "-c", "sleep 20000000"]
