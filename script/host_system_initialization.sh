#!/bin/bash
# @Author: Brady
# @Description: Ubuntu TLS Security Initiate
# @Create Time:  2020年6月2日 13:53:00
# @Last Modified time: 2022-11-2 09:06:31
# @E-mail: master@brady.top
# @Blog: https://www.brady.top
# @Version: 3.2
#-------------------------------------------------#
# 脚本主要功能说明:
# (1) Ubuntu 系统初始化操作包括IP地址设置、基础软件包更新以及安装加固。
# (2) Ubuntu 系统基于等保3.0部分、安全加固。
#-------------------------------------------------#

## 系统全局变量定义
username=zcxtong
SSHPORT=20480
ROOTPASS='ZCxtongtest@2023!+'  # 密码建议12位以上且包含数字、大小写字母以及特殊字符。
APPPASS='ZCxtongtel@2022!+'   

## 名称: err 、info 、warning
## 用途：全局Log信息打印函数
## 参数: $@
logerr() {
  printf "[$(date +'%Y-%m-%dT%H:%M:%S')]: \033[31mERROR: $@ \033[0m\n"
}
loginfo() {
  printf "[$(date +'%Y-%m-%dT%H:%M:%S')]: \033[32mINFO: $@ \033[0m\n"
}
logwarning() {
  printf "[$(date +'%Y-%m-%dT%H:%M:%S')]: \033[33mWARNING: $@ \033[0m\n"
}


## 名称: os::Network
## 用途: 网络配置相关操作脚本包括(IP地址修改)
## 参数: 无
osNetwork () {
  loginfo "[-] 操作系统网络配置相关脚本,开始执行....."
# (1) 卸载多余软件，例如 snap 软件及其服务
systemctl stop snapd snapd.socket #停止snapd相关的进程服务
apt autoremove --purge -y snapd
systemctl daemon-reload
rm -rf ~/snap /snap /var/snap /var/lib/snapd /var/cache/snapd /run/snapd
}

## 名称: os::TimedataZone
## 用途: 操作系统时间与时区同步配置
## 参数: 无
osTimedataZone () {
  loginfo "[*] 操作系统系统时间时区配置相关脚本,开始执行....."

# (1) 时间同步服务端容器(可选也可以用外部ntp服务器) : docker run -d --rm --cap-add SYS_TIME -e ALLOW_CIDR=0.0.0.0/0 -p 123:123/udp geoffh1977/chrony
echo "同步前的时间: $(date -R)"

# 方式1.Chrony 客户端配置
sudo apt install chrony
sudo grep -q "aliyun" /etc/chrony/chrony.conf || sudo tee -a /etc/chrony/chrony.conf <<'EOF'
pool ntp.aliyun.com iburst maxsources 4
keyfile /etc/chrony/chrony.keys
driftfile /var/lib/chrony/chrony.drift
logdir /var/log/chrony
maxupdateskew 100.0
rtcsync
# 允许跳跃式校时 如果在前 3 次校时中时间差大于 1.0s
makestep 1 3
EOF
systemctl enable chronyd && systemctl restart chronyd && systemctl status chronyd -l


# (2) 时区与地区设置:
cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
timedatectl set-timezone Asia/Shanghai
# dpkg-reconfigure tzdata  # 修改确认
# bash -c "echo 'Asia/Shanghai' > /etc/timezone" # 与上一条命令一样
# 将当前的 UTC 时间写入硬件时钟 (硬件时间默认为UTC)
timedatectl set-local-rtc 0
# 启用NTP时间同步：
timedatectl set-ntp yes
# 校准时间服务器-时间同步(推荐使用chronyc进行平滑同步)
chronyc tracking
# 手动校准-强制更新时间
# chronyc -a makestep
# 系统时钟同步硬件时钟
# hwclock --systohc
hwclock -w

# (3) 重启依赖于系统时间的服务
systemctl restart rsyslog.service cron.service
loginfo "[*] Tie confmigure modifiy successful! restarting chronyd rsyslog.service crond.service........."
timedatectl
}


## 名称: os::Security
## 用途: 操作系统安全加固配置脚本(符合等保要求-三级要求)
## 参数: 无
osSecurity () {
  loginfo "正在进行->操作系统安全加固(符合等保要求-三级要求)配置"

# (0) 系统用户核查配置
  loginfo "[-] 锁定或者删除多余的系统账户以及创建低权限用户"
defaultuser=(root daemon bin sys games man lp mail news uucp proxy www-data backup list irc gnats nobody systemd-network systemd-resolve systemd-timesync messagebus syslog _apt tss uuidd tcpdump landscape pollinate usbmux sshd systemd-coredump _chrony)
for i in $(cat /etc/passwd | cut -d ":" -f 1,7);do
  flag=0; name=${i%%:*}; terminal=${i##*:}
  if [[ "${terminal}" == "/bin/bash" || "${terminal}" == "/bin/sh" ]];then
    log::warning "${i} 用户，shell终端为 /bin/bash 或者 /bin/sh"
  fi
  for j in ${defaultuser[@]};do
    if [[ "${name}" == "${j}" ]];then
      flag=1
      break;
    fi
  done
  if [[ $flag -eq 0 ]];then
    log::warning "${i} 非默认用户"
  fi
done
passwd -l adm&>/dev/null 2&>/dev/null; passwd -l daemon&>/dev/null 2&>/dev/null; passwd -l bin&>/dev/null 2&>/dev/null; passwd -l sys&>/dev/null 2&>/dev/null; passwd -l lp&>/dev/null 2&>/dev/null; passwd -l uucp&>/dev/null 2&>/dev/null; passwd -l nuucp&>/dev/null 2&>/dev/null; passwd -l smmsplp&>/dev/null 2&>/dev/null; passwd -l mail&>/dev/null 2&>/dev/null; passwd -l operator&>/dev/null 2&>/dev/null; passwd -l games&>/dev/null 2&>/dev/null; passwd -l gopher&>/dev/null 2&>/dev/null; passwd -l ftp&>/dev/null 2&>/dev/null; passwd -l nobody&>/dev/null 2&>/dev/null; passwd -l nobody4&>/dev/null 2&>/dev/null; passwd -l noaccess&>/dev/null 2&>/dev/null; passwd -l listen&>/dev/null 2&>/dev/null; passwd -l webservd&>/dev/null 2&>/dev/null; passwd -l rpm&>/dev/null 2&>/dev/null; passwd -l dbus&>/dev/null 2&>/dev/null; passwd -l avahi&>/dev/null 2&>/dev/null; passwd -l mailnull&>/dev/null 2&>/dev/null; passwd -l nscd&>/dev/null 2&>/dev/null; passwd -l vcsa&>/dev/null 2&>/dev/null; passwd -l rpc&>/dev/null 2&>/dev/null; passwd -l rpcuser&>/dev/null 2&>/dev/null; passwd -l nfs&>/dev/null 2&>/dev/null; passwd -l sshd&>/dev/null 2&>/dev/null; passwd -l pcap&>/dev/null 2&>/dev/null; passwd -l ntp&>/dev/null 2&>/dev/null; passwd -l haldaemon&>/dev/null 2&>/dev/null; passwd -l distcache&>/dev/null 2&>/dev/null; passwd -l webalizer&>/dev/null 2&>/dev/null; passwd -l squid&>/dev/null 2&>/dev/null; passwd -l xfs&>/dev/null 2&>/dev/null; passwd -l gdm&>/dev/null 2&>/dev/null; passwd -l sabayon&>/dev/null 2&>/dev/null; passwd -l named&>/dev/null 2&>/dev/null
userdel -r lxd
groupdel lxd

# (2) 用户密码设置和口令策略设置
  loginfo "[-]  配置满足策略的root管理员密码 "
#echo  ${ROOTPASS} | passwd --stdin root
echo "root:$ROOTPASS" | chpasswd
loginfo "[-] 配置满足策略的app普通用户密码(根据需求配置)"
groupadd application
gpasswd -a zcxtong application
useradd -m -s /bin/bash -c "application primary user" -g application $username
#echo ${APPPASS} | passwd --stdin $username
echo "$username:$APPPASS" | chpasswd


  loginfo "[-] 存储用户密码的文件，其内容经过sha512加密，所以非常注意其权限"
touch /etc/security/opasswd && chown root:root /etc/security/opasswd && chmod 600 /etc/security/opasswd
#userdel -r yungo


# (3) 用户sudo权限以及重要目录和文件的权限设置
  loginfo "[-] 用户sudo权限以及重要目录和文件的新建默认权限设置"
# 如uBuntu安装时您创建的用户 yungo 防止直接通过 passwd 修改root密码(此时必须要求输入yungo密码后才可修改root密码)
# Tips: Sudo允许授权用户权限以另一个用户（通常是root用户）的身份运行程序,
sed -i "/# Members of the admin/i ${username} ALL=(ALL) NOPASSWD:ALL" /etc/sudoers
  loginfo "[-] 配置用户 umask 为027 "
egrep -q "^\s*umask\s+\w+.*$" /etc/profile && sed -ri "s/^\s*umask\s+\w+.*$/umask 027/" /etc/profile || echo "umask 027" >> /etc/profile
egrep -q "^\s*umask\s+\w+.*$" /etc/bash.bashrc && sed -ri "s/^\s*umask\s+\w+.*$/umask 027/" /etc/bashrc || echo "umask 027" >> /etc/bash.bashrc
# loginfo "[-] 设置用户目录创建默认权限, (初始为077比较严格)在未设置umask为027 则默认为077"
# egrep -q "^\s*(umask|UMASK)\s+\w+.*$" /etc/login.defs && sed -ri "s/^\s*(umask|UMASK)\s+\w+.*$/UMASK 022/" /etc/login.defs || echo "UMASK 022" >> /etc/login.defs

  loginfo "[-] 设置或恢复重要目录和文件的权限"
chmod 755 /etc;
chmod 777 /tmp;
chmod 700 /etc/inetd.conf&>/dev/null 2&>/dev/null;
chmod 755 /etc/passwd;
chmod 755 /etc/shadow;
chmod 644 /etc/group;
chmod 755 /etc/security;
chmod 644 /etc/services;
chmod 750 /etc/rc*.d
chmod 600 ~/.ssh/authorized_keys

  loginfo "[-] 删除潜在威胁文件 "
find / -maxdepth 3 -name hosts.equiv | xargs rm -rf
find / -maxdepth 3 -name .netrc | xargs rm -rf
find / -maxdepth 3 -name .rhosts | xargs rm -rf


# (4) SSHD 服务安全加固设置
loginfo "[-] sshd 服务安全加固设置"
# 严格模式
egrep -q "^\s*StrictModes\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*StrictModes\s+.+$/StrictModes yes/" /etc/ssh/sshd_config || echo "StrictModes yes" >> /etc/ssh/sshd_config
# 监听端口更改
if [ -e ${SSHPORT} ];then export SSHPORT=20480;fi
egrep -q "^\s*Port\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*Port\s+.+$/Port ${SSHPORT}/" /etc/ssh/sshd_config || echo "Port ${SSHPORT}" >> /etc/ssh/sshd_config
# 启用密钥登录,禁止密码登录
egrep -q "^(#)?\s*PubkeyAuthentication\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*PubkeyAuthentication\s+.+$/PubkeyAuthentication yes/" /etc/ssh/sshd_config || echo "PubkeyAuthentication yes" >> /etc/ssh/sshd_config
egrep -q "^\s*ChallengeResponseAuthentication\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*ChallengeResponseAuthentication\s+.+$/ChallengeResponseAuthentication no/" /etc/ssh/sshd_config || echo "ChallengeResponseAuthentication no" >> /etc/ssh/sshd_config
egrep -q "^\s*UsePAM\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*UsePAM\s+.+$/UsePAM yes/" /etc/ssh/sshd_config || echo "UsePAM yes" >> /etc/ssh/sshd_config
egrep -q "^(#)?\s*PasswordAuthentication\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*PasswordAuthentication\s+.+$/PasswordAuthentication no/" /etc/ssh/sshd_config || echo "PasswordAuthentication no" >> /etc/ssh/sshd_config
# 禁用X11转发以及端口转发
egrep -q "^\s*X11Forwarding\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*X11Forwarding\s+.+$/X11Forwarding no/" /etc/ssh/sshd_config || echo "X11Forwarding no" >> /etc/ssh/sshd_config
egrep -q "^\s*X11UseLocalhost\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*X11UseLocalhost\s+.+$/X11UseLocalhost yes/" /etc/ssh/sshd_config || echo "X11UseLocalhost yes" >> /etc/ssh/sshd_config
egrep -q "^\s*AllowTcpForwarding\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*AllowTcpForwarding\s+.+$/AllowTcpForwarding no/" /etc/ssh/sshd_config || echo "AllowTcpForwarding no" >> /etc/ssh/sshd_config
egrep -q "^\s*AllowAgentForwarding\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*AllowAgentForwarding\s+.+$/AllowAgentForwarding no/" /etc/ssh/sshd_config || echo "AllowAgentForwarding no" >> /etc/ssh/sshd_config
# 关闭禁用用户的 .rhosts 文件  ~/.ssh/.rhosts 来做为认证: 缺省IgnoreRhosts yes
egrep -q "^(#)?\s*IgnoreRhosts\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^(#)?\s*IgnoreRhosts\s+.+$/IgnoreRhosts yes/" /etc/ssh/sshd_config || echo "IgnoreRhosts yes" >> /etc/ssh/sshd_config
# 禁止root远程登录（推荐配置-根据需求配置）
egrep -q "^\s*PermitRootLogin\s+.+$" /etc/ssh/sshd_config && sed -ri "s/^\s*PermitRootLogin\s+.+$/PermitRootLogin no/" /etc/ssh/sshd_config || echo "PermitRootLogin no" >> /etc/ssh/sshd_config
# 登陆前后欢迎提示设置
egrep -q "^\s*(banner|Banner)\s+\W+.*$" /etc/ssh/sshd_config && sed -ri "s/^\s*(banner|Banner)\s+\W+.*$/Banner \/etc\/issue/" /etc/ssh/sshd_config || \
echo "Banner /etc/issue" >> /etc/ssh/sshd_config
loginfo "[-] 远程SSH登录前后提示警告Banner设置"
# SSH登录前警告Banner
tee /etc/issue <<'EOF'
****************** [ 安全登陆 (Security Login) ] *****************
Authorized only. All activity will be monitored and reported.By Security Center.

EOF
# SSH登录后提示Banner
sed -i '/^fi/a\\n\necho "\\e[1;37;41;5m################## 安全运维 (Security Operation) ####################\\e[0m"\necho "\\e[32mLogin success. Please execute the commands and operation data carefully.By\\e[0m"' /etc/update-motd.d/00-header
systemctl stop multipathd.socket && systemctl restart sshd


# (5) 用户远程登录失败次数与终端超时设置
  loginfo "[-] 用户远程连续登录失败5次锁定帐号5分钟包括root账号"
sed -ri "/^\s*auth\s+required\s+pam_tally2.so\s+.+(\s*#.*)?\s*$/d" /etc/pam.d/sshd
sed -ri '2a auth required pam_tally2.so deny=5 unlock_time=300 even_deny_root root_unlock_time=300' /etc/pam.d/sshd
# 宿主机控制台登陆(可选)
# sed -ri "/^\s*auth\s+required\s+pam_tally2.so\s+.+(\s*#.*)?\s*$/d" /etc/pam.d/login
# sed -ri '2a auth required pam_tally2.so deny=5 unlock_time=300 even_deny_root root_unlock_time=300' /etc/pam.d/login

  loginfo "[-] 设置登录超时时间为10分钟 "
egrep -q "^\s*(export|)\s*TMOUT\S\w+.*$" /etc/profile && sed -ri "s/^\s*(export|)\s*TMOUT.\S\w+.*$/export TMOUT=600\nreadonly TMOUT/" /etc/profile || echo -e "export TMOUT=600\nreadonly TMOUT" >> /etc/profile
egrep -q "^\s*.*ClientAliveInterval\s\w+.*$" /etc/ssh/sshd_config && sed -ri "s/^\s*.*ClientAliveInterval\s\w+.*$/ClientAliveInterval 600/" /etc/ssh/sshd_config || echo "ClientAliveInterval 600" >> /etc/ssh/sshd_config


# (5) 切换用户日志记录或者切换命令更改(可选)
  loginfo "[-] 切换用户日志记录和切换命令更改名称为SU "
egrep -q "^(\s*)SULOG_FILE\s+\S*(\s*#.*)?\s*$" /etc/login.defs && sed -ri "s/^(\s*)SULOG_FILE\s+\S*(\s*#.*)?\s*$/\SULOG_FILE  \/var\/log\/.history\/sulog/" /etc/login.defs || echo "SULOG_FILE  /var/log/.history/sulog" >> /etc/login.defs
egrep -q "^\s*SU_NAME\s+\S*(\s*#.*)?\s*$" /etc/login.defs && sed -ri "s/^(\s*)SU_NAME\s+\S*(\s*#.*)?\s*$/\SU_NAME  switch_user/" /etc/login.defs || echo "SU_NAME  switch_user" >> /etc/login.defs
mkdir -vp /var/log/.backups /usr/local/bin /var/log/.history
cp /usr/bin/su /var/.backups/su.bak
mv /usr/bin/su /usr/bin/SU
chmod 777 /var/log/.history

# (6) 用户终端执行的历史命令记录
loginfo "[-] 用户终端执行的历史命令记录 "
egrep -q "^HISTSIZE\W\w+.*$" /etc/profile && sed -ri "s/^HISTSIZE\W\w+.*$/HISTSIZE=101/" /etc/profile || echo "HISTSIZE=101" >> /etc/profile
# 方式1
tee /etc/profile.d/history-record.sh <<'EOF'
# 历史命令执行记录文件路径
LOGTIME=$(date +%Y%m%d-%H-%M-%S)
export HISTFILE="/var/log/.history/${USER}.${LOGTIME}.history"
if [ ! -f ${HISTFILE} ];then
  touch ${HISTFILE}
fi
chmod 600 /var/log/.history/${USER}.${LOGTIME}.history
# 历史命令执行文件大小记录设置
HISTFILESIZE=128
HISTTIMEFORMAT="%F_%T $(whoami)#$(who -u am i 2>/dev/null| awk '{print $NF}'|sed -e 's/[()]//g'):"
EOF



# (7) GRUB 安全设置 （需要手动设置请按照需求设置）
  loginfo "[-] 系统 GRUB 安全设置(防止物理接触从grub菜单中修改密码) "
# Grub 关键文件备份
sudo cp -a /etc/grub.d/00_header /var/log/.backups
sudo cp -a /etc/grub.d/10_linux /var/log/.backups
# 设置Grub菜单界面显示时间
sudo sed -i -e 's|GRUB_TIMEOUT_STYLE=hidden|#GRUB_TIMEOUT_STYLE=hidden|g' -e 's|GRUB_TIMEOUT=0|GRUB_TIMEOUT=3|g' /etc/default/grub
sudo sed -i -e 's|set timeout_style=${style}|#set timeout_style=${style}|g' -e 's|set timeout=${timeout}|set timeout=3|g' /etc/grub.d/00_header
# 创建认证密码 (此处密码: 自定义手动输入)
sudo grub-mkpasswd-pbkdf2
# Enter password:
# Reenter password:
# 设置认证用户以及password_pbkdf2认证
sudo tee -a /etc/grub.d/00_header <<'END'
sudo cat <<'EOF'
# GRUB Authentication
set superusers="grub"
password_pbkdf2 grub grub.pbkdf2.sha512.10000.5FD0269A1E1216B31ED1F127DF6E47D164D85E37E6187A48341F5665092CC752DB1527C0D928080A3440C0F46E7B7C749EC3582AF0F02951EB0FB01F9F8424D0.1466487161CC5866FCE719D95DD1D70FAF67B4F8601804DC74B3FE3C82506648942FF7073C8BFCFE268FBC0D545BED047A7763D0131E0ABBF8A4E0922C52EFD5
EOF
END

sudo sed -i '/echo "$title" | grub_quote/ { s/menuentry /menuentry --user=grub /;}' /etc/grub.d/10_linux
sudo sed -i '/echo "$os" | grub_quote/ { s/menuentry /menuentry --unrestricted /;}' /etc/grub.d/10_linux

# Ubuntu 方式更新GRUB从而生成boot启动文件。
sudo update-grub


# (8) 操作系统防火墙启用以及策略设置
#  loginfo "[-] 系统防火墙启用以及规则设置 "
# systemctl enable ufw.service && systemctl start ufw.service && ufw enable
# ufw allow proto tcp from 172.10.2.0/24 to any port 20480
# ufw allow proto tcp from 10.0.3.0/16 to any port 20480
# ufw allow proto tcp to any port 20597

# (9) 核心文件加上不可更改属性
  loginfo "[-] 核心文件加上不可更改属性 "
sudo mkdir -p /home/$username/.ssh && sudo touch /home/$username/.ssh/authorized_keys
sudo chown -R $username:application /home/$username/

sudo cat > /home/$username/.ssh/authorized_keys  << EOF
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCoZ6vHOWm8UfIZqEsy4NSmyFLRf9RORIqeuck3JcCujK+NmdNDbujUwSzP4utybeZxFQFrfI1D7ArzC0scZcArix7hzdeEChscVCWTBxnvqEG1ZqeLvTnAr7UrWw202jTamDcjKdeKlz7rXq6uhDJmj8r5VQLt+r67ajLjHLb+YrkrNQeGu3726ps+CEtOh9LNJRiOKDa0cjeU70dWFCkKi4uqIJoW/zjnFXpEoovvPSgEzmnBI2RwWOYCICAUoYTsOz3x56UcQ5wGxQhvoU1HukF8Psi9+7VAlUPtNjqNqpPKCnIXcTFEGxpl4PoVcSWj40xjiwC1bX8FnWUPisHYQNXorK7S7jCj/7jEnJOVrUPv9qPfUIWm91M7+9jqu6ixZZjKJep/P/hjmP8OE2MqqqN4UEoQB+tJ0RF89fBPQNCCdbwji5utG1gPBcwdVhKju9Aq3jFxyfwzXG/VNuhBU79JaEUZZmEUwGirw6gE88aO4Sn128XPmGnvw2/rf4E= root@test
EOF
sudo chattr +iae /root/.ssh/authorized_keys
sudo chattr +iae /home/$username/.ssh/authorized_keys


#(10) 敏感命令提醒
  loginfo "[-] 敏感命令提醒 "
sudo  touch /etc/profile.d/alias.sh
sudo  cat > /etc/profile.d/alias.sh << EOF
alias reboot='echo -e "\033[41;05m  危险!!!! -检查后-确认-【谨慎使用重启】 \033[0m "'
alias rm='echo -e "\033[41;05m  危险!!!! -检查后-确认-【谨慎使用删除】 \033[0m "'
[ -d /tmp/remove/ ] ||  mkdir   /tmp/remove/
EOF
sudo source /etc/profile

}

## 名称: os::InstallSoft
## 用途: 操作系统时间与时区同步配置
## 参数: 无
#osInstallSoft () {

#(1) 安装系统依赖及必要软件  
sudo apt install -y cpufrequtils unzip mesa-opencl-icd ocl-icd-opencl-dev nfs-common ntpdate mdadm xfsprogs  jq pkg-config  curl  supervisor rustc  smartmontools hdparm  lxcfs   sysstat  lm-sensors iotop nload nfs-kernel-server make  gcc  nvme-cli hwloc net-tools lrzsz
sudo apt install -y gcc make  cpufrequtils ntp bash-completion  libhwloc-dev mesa-opencl-icd ocl-icd-opencl-dev nfs-common ntpdate mdadm xfsprogs   bzr jq pkg-config curl supervisor rustc smartmontools hdparm lxcfs iperf3
sudo apt install -y supervisor numactl
sudo apt-get update  -y

##静默安装GPU
echo "#################静默安装GPU##########################"

sudo sed  -i '/^blacklist nouveau/d' /etc/modprobe.d/blacklist.conf
sudo sed  -i '/^options nouveau modeset=0/d' /etc/modprobe.d/blacklist.conf
sudo  echo  'blacklist nouveau' >>  /etc/modprobe.d/blacklist.conf
sudo echo  'options nouveau modeset=0' >> /etc/modprobe.d/blacklist.conf
#apt-get remove --purge nvidia-*
rmmod nouveau

###cp
#cd /root/
#curl -u  'ftpuser:LGvC2#%BsHPrRZEqns!GgCh'        ftp://10.0.1.1/onup/NVIDIA-Linux-x86_64-460.84.run    -O
#chmod +x NVIDIA-Linux-x86_64-460.84.run
#/root/NVIDIA-Linux-x86_64-460.84.run --accept-license --silent --no-nouveau-check --disable-nouveau --no-opengl-files
#
sudo apt-get purge nvidia*  -y
sudo add-apt-repository ppa:graphics-drivers -y
sudo apt-get update  -y
sudo apt-get install nvidia-driver-440 nvidia-settings nvidia-prime  -y




echo " "
sleep  3
####################################################

#禁用桌面
echo "#################禁用桌面##########################"
apt-get remove lightdm -y
sudo systemctl set-default multi-user.target
sudo systemctl disable multi-user.target
repaddline "^GRUB_CMDLINE_LINUX_DEFAULT" "GRUB_CMDLINE_LINUX_DEFAULT=\"text\"" /etc/default/grub
sudo apt-get remove gnome-shell  -y
sudo apt-get remove gnome
sudo apt-get autoremove -y
sudo apt-get purge gnome
sudo apt-get autoclean
sudo apt-get clean
sudo update-grub



#(2) 启动守护进程
sudo systemctl enable supervisor.service && sudo systemctl start supervisor.service

#}

echo "loginfo"
echo "logerr"
echo "logwarning"

echo "osNetwork"
#osNetwork
#echo "osTimedataZone"
#osTimedataZone
echo "osSecurity"
#osSecurity
echo "osInstallSoft"
#osInstallSoft

