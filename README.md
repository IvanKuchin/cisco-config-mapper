# Cisco config mapper 

This tool designed for re-using production config in lab with minor tweaks.

- Replace interface names
- Anonymize users
- Hide any password from running config 
- Remove prduction SSL-certificates
- and so on, as long as config represented in a text format

## Workflow
1. Save `sh run` output to `src/` folder. For convenience use hostname as a filename
2. Craft config files in the `config/` folder
3. Run `cisco-config-mapper`

## Configuration explained
Assume `sh run` saved into `src/r1`

```cisco
feature bfd
service timestamps debug datetime msec
service timestamps log datetime msec

hostname R1

aaa new-model
aaa authentication login default local
aaa authentication login CONSOLE none

aaa group server tacacs+ NetAdmin
 server name ISE1
 server name ISE2

interface GigabitEthernet1.100
 ip address 10.0.100.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled
interface GigabitEthernet1.200
 ip address 10.0.200.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled

interface Vlan90
  no ip redirects
  ip address 10.90.90.3/24
  no ipv6 redirects
  ip ospf passive-interface
  ip router ospf 10 area 0.0.0.90
  hsrp version 2
  hsrp 90 
    priority 100 forwarding-threshold lower 0 upper 100
    timers  1  3
    ip 10.90.90.1
  description uplink
  no shutdown

```

### Cleanup - remove lines
`config/r1`
```yaml
remove-lines:
- 'feature bfd'
- 'hostname R1'
```

will produce the output
```cisco
service timestamps debug datetime msec
service timestamps log datetime msec

aaa new-model
aaa authentication login default local
aaa authentication login CONSOLE none

aaa group server tacacs+ NetAdmin
 server name ISE1
 server name ISE2

interface GigabitEthernet1.100
 ip address 10.0.100.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled
interface GigabitEthernet1.200
 ip address 10.0.200.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled

interface Vlan90
  no ip redirects
  ip address 10.90.90.3/24
  no ipv6 redirects
  ip ospf passive-interface
  ip router ospf 10 area 0.0.0.90
  hsrp version 2
  hsrp 90 
    priority 100 forwarding-threshold lower 0 upper 100
    timers  1  3
    ip 10.90.90.1
  description uplink
  no shutdown
```

### Cleanup - remove prefixes
`config/r1`
```yaml
remove-prefixes:
- 'aaa'
- 'service'
```

Notice removal of `aaa` also removes `aaa`-blocks
will produce the output.

```cisco
interface GigabitEthernet1.100
 ip address 10.0.100.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled
interface GigabitEthernet1.200
 ip address 10.0.200.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled

interface Vlan90
  no ip redirects
  ip address 10.90.90.3/24
  no ipv6 redirects
  ip ospf passive-interface
  ip router ospf 10 area 0.0.0.90
  hsrp version 2
  hsrp 90 
    priority 100 forwarding-threshold lower 0 upper 100
    timers  1  3
    ip 10.90.90.1
  description uplink
  no shutdown
```

### Cleanup - remove under interfaces
`config/r1`
```yaml
- from: Vlan90
  to: Vlan90
  remove-prefixes:
  - '  hsrp'
```

will produce the output.

```cisco
interface GigabitEthernet1.100
 ip address 10.0.100.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled
interface GigabitEthernet1.200
 ip address 10.0.200.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled

interface Vlan90
  no ip redirects
  ip address 10.90.90.3/24
  no ipv6 redirects
  ip ospf passive-interface
  ip router ospf 10 area 0.0.0.90
  description uplink
  no shutdown
```

### Rename interface
`config/r1`
```yaml
- from: GigabitEthernet1
  to: GigabitEthernet23/45
```

will produce the output.

```cisco
interface GigabitEthernet23/45.100
 ip address 10.0.100.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled
interface GigabitEthernet23/45.200
 ip address 10.0.200.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled

interface Vlan90
  no ip redirects
  ip address 10.90.90.3/24
  no ipv6 redirects
  ip ospf passive-interface
  ip router ospf 10 area 0.0.0.90
  description uplink
  no shutdown
```

### Change interfaces IP address
`config/r1`
```yaml
- from: Vlan90
  to: Vlan90
  remove-lines:
  - '  ip address 10.90.90.3/24'
  append: ip address 10.90.90.2/24
```

will produce the output.

```cisco
interface GigabitEthernet23/45.100
 ip address 10.0.100.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled
interface GigabitEthernet23/45.200
 ip address 10.0.200.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled

interface Vlan90
  no ip redirects
  ip address 10.90.90.2/24
  no ipv6 redirects
  ip ospf passive-interface
  ip router ospf 10 area 0.0.0.90
  description uplink
  no shutdown
```

### Prepend / append
`config/r1`
```yaml
prepend: hostname lab-r1
append: |
  line vty 0 15
    login
    password cisco
```

will produce the output.

```cisco
hostname lab-r1

interface GigabitEthernet23/45.100
 ip address 10.0.100.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled
interface GigabitEthernet23/45.200
 ip address 10.0.200.2 255.255.255.0
 negotiation auto
 bfd interval 1000 min_rx 1000 multiplier 10
 no mop enabled

interface Vlan90
  no ip redirects
  ip address 10.90.90.2/24
  no ipv6 redirects
  ip ospf passive-interface
  ip router ospf 10 area 0.0.0.90
  description uplink
  no shutdown

line vty 0 15
  login
  password cisco
```

