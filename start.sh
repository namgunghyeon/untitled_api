kill $(ps -ef | grep './main' | grep -v grep | awk '{print $2}')
nohup ./main > logs/log.log
