#!/usr/bin/sh

apt-get install fontconfig ttf-mscorefonts-installer fonts-dejavu fonts-liberation fonts-crosextra-carlito fonts-takao-gothic fonts-opensymbol fonts-beng-extra fonts-gujr-extra fonts-telu-extra -y --no-install-recommends
fc-cache -fv

#Start generate AllFonts.js, font thumbnails and font_selection.bin
echo -n Generating AllFonts.js, please wait...

"./allfontsgen"\
  --input="`pwd`/core-fonts"\
  --allfonts="./js/AllFonts.js"\
  --images="./js/sdkjs/common/Images"\
  --selection="./font_selection.bin"\
  --use-system="true"
