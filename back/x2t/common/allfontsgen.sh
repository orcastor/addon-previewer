#!/usr/bin/sh

# if there's 'Unable to locate package' errors, use command below to add source for apt
#echo "deb http://mirrors.aliyun.com/ubuntu/ vivid main restricted universe multiverse
#deb http://mirrors.aliyun.com/ubuntu/ vivid-security main restricted universe multiverse
#deb http://mirrors.aliyun.com/ubuntu/ vivid-updates main restricted universe multiverse
#deb http://mirrors.aliyun.com/ubuntu/ vivid-proposed main restricted universe multiverse
#deb http://mirrors.aliyun.com/ubuntu/ vivid-backports main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ vivid main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ vivid-security main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ vivid-updates main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ vivid-proposed main restricted universe multiverse
#deb-src http://mirrors.aliyun.com/ubuntu/ vivid-backports main restricted universe multiverse" >> /etc/apt/sources.list
#apt-get update

apt-get install fontconfig fonts-dejavu fonts-liberation fonts-crosextra-carlito fonts-takao-gothic fonts-opensymbol fonts-beng-extra fonts-gujr-extra fonts-telu-extra fonts-nanum fonts-noto fonts-wqy-zenhei fonts-arphic-ukai fonts-freefont-ttf fonts-takao-gothic -y --no-install-recommends
fc-cache -fv

#Start generate AllFonts.js, font thumbnails and font_selection.bin
echo Generating AllFonts.js, please wait...

"./allfontsgen"\
  --input="`pwd`/core-fonts"\
  --allfonts="./js/AllFonts.js"\
  --images="./js/sdkjs/common/Images"\
  --selection="./font_selection.bin"\
  --use-system="true"

echo Done
