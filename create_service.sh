#!/bin/Bash

#1단계 : 작업 디렉토리 & 서비스 명 받기
read -p "설치된 디렉토리 (ex: /mnt/openNAMU) : " working_directory
echo $working_directory

read -p "원하는 서비스명 (ex: opennamu): " service_name
echo $service_name

read -p "설명 (ex: OpenNAMU 서비스입니다): " description
echo $description

read -p "로그위치 (ex: /var/log/openNAMU.log) : " log_path
echo $log_path


#2단계 : 파일제작
cat <<EOF > /lib/systemd/system/${service_name}.service
[Unit]
Description=$description

[Service]
Type=simple

WorkingDirectory=$working_directory
ExecStart=/usr/bin/python3 $working_directory/app.py
Restart=on-failure
PIDFile=/run/$service_name.pid

#rsyslog 사용
#StandardOutput=syslog
#StandardError=syslog
#SyslogIdentifier=$service_name

#systemctl 245 이후 로깅
StandardOutput=append:$log_path
StandardError=append:$log_path

[Install]
WantedBy=multi-user.target
EOF


#3단계 : 서비스 확인
systemctl daemon-reload 
systemctl start $service_name
systemctl status $service_name
