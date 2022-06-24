#/bin/sh
if [ `ps -ef | grep '/home/samba/espeniel/ddaom_mainapi/ddaom' | wc | awk '{print $1}'` -gt 1 ]
then
 echo "RUNNING"
else
 nohup /home/samba/espeniel/ddaom_mainapi/ddaom &
fi
