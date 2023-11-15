rm -f /etc/sysconfig/network/ifcfg-eth*
rm -f /etc/udev/rules.d/51-qeth-0.0.*
> /boot/zipl/active_devices.txt
cat 0000 > /etc/sysconfig/network/ifcfg-eth1000
cat 0001 > /etc/udev/rules.d/51-qeth-0.0.1000.rules
cat 0002 > /etc/udev/rules.d/70-persistent-net.rules
cat 0003 > /tmp/znetconfig.sh
sleep 2
/bin/bash /tmp/znetconfig.sh
rm -rf invokeScript.sh
