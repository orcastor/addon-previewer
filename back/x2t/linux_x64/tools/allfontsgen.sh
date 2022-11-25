#!/bin/sh

#Start generate AllFonts.js, font thumbnails and font_selection.bin
echo -n Generating AllFonts.js, please wait...

"./allfontsgen"\
  --input="../core-fonts"\
  --allfonts="../js/AllFonts.js"\
  --images="../js/sdkjs/common/Images"\
  --selection="../font_selection.bin"\
  --use-system="true"