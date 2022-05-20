#!/bin/bash
echo cdgs$8133Tong|/openconnect/openconnect --protocol=gp extrapass.cdg.co.th --user=003459 --no-dtls -b

sleep 5

iptables -t nat -A POSTROUTING -o tun0 -j MASQUERADE
iptables -A FORWARD -i eth0 -j ACCEPT

/bin/bash