#!/usr/bin/sh

# if there's an 'Unable to locate package' error, change /etc/apt/sources.list and update
#apt-get update

apt-get install fontconfig fonts-dejavu fonts-liberation fonts-crosextra-carlito fonts-opensymbol fonts-noto-core fonts-freefont-ttf fonts-droid-fallback fonts-open-sans fonts-symbola ttf-mscorefonts-installer -y --no-install-recommends
fc-cache -fv

#Start generate AllFonts.js, font thumbnails and font_selection.bin
echo Generating AllFonts.js, please wait...

cp -f "./libgraphics.so" "./libgraphics.so.bak"

"./upx" -d "./libgraphics.so.bak" -o "./libgraphics.so" -f > /dev/null 2>&1

"./allfontsgen"\
  --input="`pwd`/core-fonts"\
  --allfonts="./js/AllFonts.js"\
  --selection="./font_selection.bin"\
  --use-system="true"

#--images="./js/sdkjs/common/Images"

mv -f "./libgraphics.so.bak" "./libgraphics.so"

echo Done.
