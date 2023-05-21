

# Commands

## PRINT

```
route PRINT -4
```

"-4" -> IPv4 only

```
PS C:\Users\tider\Projects\IP-Tool> route PRINT -4
===========================================================================
Interface List
  6...00 d8 61 1a 8e 7b ......Intel(R) Ethernet Connection (7) I219-V
  1...........................Software Loopback Interface 1
===========================================================================

IPv4 Route Table
===========================================================================
Active Routes:
Network Destination        Netmask          Gateway       Interface  Metric
          0.0.0.0          0.0.0.0    192.168.1.254      192.168.1.9     26
        127.0.0.0        255.0.0.0         On-link         127.0.0.1    331
        127.0.0.1  255.255.255.255         On-link         127.0.0.1    331
  127.255.255.255  255.255.255.255         On-link         127.0.0.1    331
      192.168.1.0    255.255.255.0         On-link       192.168.1.9    281
      192.168.1.9  255.255.255.255         On-link       192.168.1.9    281
    192.168.1.255  255.255.255.255         On-link       192.168.1.9    281
        224.0.0.0        240.0.0.0         On-link         127.0.0.1    331
        224.0.0.0        240.0.0.0         On-link       192.168.1.9    281
  255.255.255.255  255.255.255.255         On-link         127.0.0.1    331
  255.255.255.255  255.255.255.255         On-link       192.168.1.9    281
===========================================================================
Persistent Routes:
  Network Address          Netmask  Gateway Address  Metric
          0.0.0.0          0.0.0.0    192.168.1.254       1
===========================================================================
```

## ADD

route -p ADD 157.0.0.0 MASK 255.0.0.0 157.55.80.1 METRIC 3 IF 2
          destination^      ^mask     ^gateway     metric^    ^interface
            
"-p" -> persistent across boots of the system

If IF is not given, it tries to find the best interface for a given gateway.

```
PS C:\Users\tider\Projects\IP-Tool> route ADD 192.168.2.0 MASK 255.255.255.0 192.168.1.254 METRIC 1
 OK!
```

## CHANGE

route CHANGE 157.0.0.0 MASK 255.0.0.0 157.55.80.5 METRIC 2 IF 2
          destination^      ^mask     ^gateway     metric^    ^interface

CHANGE is used to modify gateway and/or metric only.

```
PS C:\Users\tider\Projects\IP-Tool> route CHANGE 192.168.2.0 MASK 255.255.255.0 192.168.1.254 METRIC 1
 OK!
```

"-p" -> persistent across boots of the system

If the route you change is Persistent the "-p" needs to be in the command otherwise it will be removed from the Persistent Routes. The reverse is also true. If the route you change wasn't added with the -p flag it can be CHANGED to be Persistent

## DELETE

route DELETE 157.0.0.0
          destination^

```
PS C:\Users\tider\Projects\IP-Tool> route DELETE 157.0.0.0
 OK!
```

## RESET routing table (requires reboot)

route /F

```
IPv4 Route Table
===========================================================================
Active Routes:
  None
```

you can get some routes back by changing the IP of an interface, but the loopback routes need added back manually or reboot

```
IPv4 Route Table
===========================================================================
Active Routes:
Network Destination        Netmask          Gateway       Interface  Metric
          0.0.0.0          0.0.0.0    192.168.1.254    192.168.1.110     26
      192.168.1.0    255.255.255.0         On-link     192.168.1.110    281
    192.168.1.110  255.255.255.255         On-link     192.168.1.110    281
    192.168.1.255  255.255.255.255         On-link     192.168.1.110    281
===========================================================================
```

## Reboot PC

shutdown /r

## RESET routing table (doesn't work?)

netsh interface ip delete destinationcache

## OTHER

tracert 192.168.1.95

arp -a


