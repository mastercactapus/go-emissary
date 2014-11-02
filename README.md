#Emissary
systemd integration for consul

Inspired by fleet, emissary is a tool to allow cluster management with systemd using consul.

##Emissary Unit File Options

```systemd-unit

[X-Emissary]

# Only deploy to a machine running redis.service
# Not using definitions under [Unit] to separate
# system deps from deploy deps
#
# If no machine is running required services
# they will be deployed along side
Requires=redis.service
Requires=mysql.service

# Locks this service to the given machine id(s)
Machine=core-01
Machine=core-02

# Locks this service to the given datacenter(s)
Datacenter=dc1

# Locks this service to machines with the given tag(s)
Tag=external-address
Tag=external-storage

# Keep this service running on existing and new nodes;
# respecting machine, datacenter, and tag rules
Global=false

# Link a consul check to the systemd state of this service
# anything other than active/running is considered critical
# the 'note' will reflect actual systemd status
Monitor=true

```
