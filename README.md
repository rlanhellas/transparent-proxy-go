Iptables rules to run on bastion container (access the container and run these commands by yourself)
````
iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to 8080
iptables -t nat -A OUTPUT -p tcp --dport 443 -j REDIRECT --to 8080
iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to 8080
iptables -t nat -A OUTPUT -p tcp --dport 80 -j REDIRECT --to 8080
````

How to test ?
- Follow the logs with `docker logs -f proxy`, then you will see the host returned by socket connection
- Access the bastion container and execute `curl https://google.com`.