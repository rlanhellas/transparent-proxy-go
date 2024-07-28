FROM ubuntu:24.04

RUN apt update && apt install -y curl dnsutils iputils-ping iptables tcpdump strace

CMD ["sh", "-c", "sleep 20000000"]
