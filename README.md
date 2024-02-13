<p align="center">
  <a href="https://orcastor.github.io/doc/">
    <img src="https://orcastor.github.io/doc/logo.svg">
  </a>
</p>

<p align="center"><strong>OrcaS 文档预览插件</strong></p>

<p align="center">—— 基于h5，纯离线本地版，效果/兼容性更好，可自行适配存储 ——</p>

|类型|进度|展示|兼容转换方案|支持格式|注意事项|
|-|-|-|-|-|-|
|pdf|✅|[vue-pdf-embed](https://github.com/hrynko/vue-pdf-embed)|-|pdf|支持有密码文件|
|文字|✅|[docxjs](https://www.npmjs.com/package/docx-preview)|[x2t](https://github.com/ONLYOFFICE/core)转换成docx|docx、doc、wps等|\* 字体问题；有密码文档问题|
|表格|✅|[luckyexcel](https://github.com/dream-num/Luckyexcel) |[x2t](https://github.com/ONLYOFFICE/core)转换成xlsx|xlsx、xls、csv、et等|\* 字体问题；有密码文档问题|
|演示|✅|同pdf|[x2t](https://github.com/ONLYOFFICE/core)转换成pdf|pptx、ppt、dps等|\* 字体问题；有密码文档问题|
|iWork|✅|h5原生|[iwork2html](https://github.com/orcastor/iwork-converter)转换成html|pages、numbers、key等|\* 临时方案|
|文本|✅|同docx|[x2t](https://github.com/ONLYOFFICE/core)转换成docx|txt、json/toml/yml/config/xml等|enca/file检测编码防止中文乱码|
|代码||md代码染色 / WebIDE|-|cpp、c、h、java、py、go、php、js、html、css等|enca/file检测编码防止中文乱码|
|设计|✅|同图片|[cad2x](https://github.com/orcastor/cad2x-converter)转换成png|dwg、dwt、dxf等|\* 字体问题；代码页问题|
|图片|✅|[el-image](https://element.eleme.cn/#/zh-CN/component/image)|[vips](https://github.com/libvips/libvips)转换缩略图|png、jpg、gif、bmp、ico、icns、eps、psd等；文档需要先获取首页截图|gif需要处理多帧；vips支持SIMD/OpenMP加速|
|视频|⌛|[vue-plyr](https://github.com/redxtech/vue-plyr)|[ffmpeg](https://github.com/FFmpeg/FFmpeg)转换成HLS(m3u8)|mp4、wmv、mkv等|fps和码率；可以尝试h.265；支持GPU加速（OpenCL/Vulkan）|
|音频||[aPlayer](https://github.com/DIYgod/APlayer)/[cPlayer](https://github.com/MoePlayer/cPlayer)|-|mp3,wma,wav,ape,flac,ogg,aac等|支持匹配封面、lrc歌词文件|
|存档|⌛||写到临时目录解压|rar、zip、7z、dmg、iso等|有密码的文档|
|备份|||imobax & abe|iOS备份目录，安卓备份ab文件等|有密码的备份|
|图标|⌛|h5原生|[f2ico](https://github.com/orcastor/f2ico)转换ico|exe、apk、autorun.inf、desktop.ini、\*.desktop(\*.AppImage/\*.run)、\*.app目录等|

## 注意事项：
- 如果转换插件在docker镜像中部署，需要配置`ORCAS_DOCKER_EXEC="docker exec -i <container_id>"`环境变量，其中`container_id`为容器的id值
- 打包后的文件和webapp的放置到一起：
`ln -s $(addon-previewer)/front/dist $(webapp)/dist/prvw`

## 关于x2t

#### 部署方法：

- 把[back/x2t](https://github.com/orcastor/addon-previewer/tree/main/back/x2t)下的`common`和对应平台（`linux_arm64`/`linux_x64`）目录下的文件都拷贝到`/opt/x2t`下
- 执行`sh allfontsgen.sh`安装字体文件、生成字体列表

#### 方案优势：

- 更好的性能

  - 纯C++开发，无依赖
  - 原生支持同时并发转换多个文件
  - 秒级启动，不需要常驻

- 更小的体积
  - x2t完整功能支持仅88.8MB：可执行文件31.5MB + 配置文件57.3MB
    - 对比LibreOffice（6.4.7.2）需要401MB：jvm运行环境179MB + 包222MB（无字体）
  - 字体文件166MB左右：常见字体82.8MB（见core-fonts） + 下载82.9MB（见allfontsgen.sh）
    <details><summary>移除了非中文字体，可按需补回</summary>
      <table>
        <tbody>
          <tr><th>类别</th><th>名称</th><th>大小</th><th>说明</th></tr>
          <tr><td>韩文</td><td>nanum</td><td>34.3MB</td><td>韩国的Naver公司开发</td></tr>
          <tr><td>藏文</td><td>TibetanMachineUni</td><td>4.3MB</td><td>用于显示藏文的Unicode字体</td></tr>
          <tr><td>高棉文</td><td>khmeros</td><td>520KB</td><td>支持柬埔寨语言的文字显示</td></tr>
          <tr><td>缅甸文</td><td>padauk</td><td>1.76MB</td><td>支持缅甸语言的文字显示</td></tr>
          <tr><td>孟加拉文</td><td>beng-extra</td><td>678KB</td><td>支持孟加拉语的文字显示</td></tr>
          <tr><td>阿拉伯文</td><td>kacst</td><td>1.06MB</td><td>KACST（沙特阿拉伯国王阿卜杜拉兹国王科学技术城）开发，旨在支持阿拉伯语的文字显示</td></tr>
          <tr><td>拉丁、希腊文等</td><td>asana</td><td>760KB</td><td>由 SIL International（国际文字系统学会）开发，支持拉丁字母、希腊字母和西里尔字母等</td></tr>
          <tr><td>埃塞俄比亚文</td><td>abyssinica</td><td>596KB</td><td>支持埃塞俄比亚文字（比如阿姆哈拉语和提格雷语等）的显示</td></tr>
          <tr><td rowspan="3">日文</td><td>takao-gothic</td><td>6MB</td><td>支持日文文字显示，包括Takao明朝（Takao Mincho）和Takao高桥（Takao Gothic）</td></tr>
          <tr><td>takao-pgothic</td><td>6MB</td><td>支持日文文字显示，无衬线字体，包括Takao P明朝（Takao P Mincho）和Takao Pゴシック（Takao P Gothic）</td></tr>
          <tr><td>mona</td><td>2.8MB</td><td>支持日文文字显示</td></tr>
          <tr><td rowspan="4">印度文</td><td>samyak</td><td>324KB</td><td>支持印度语言（包括印地语、马拉地语等）的文字显示</td></tr>
          <tr><td>lohit</td><td>1.3MB</td><td>支持印度语言（印地语、孟加拉语、古吉拉特语、卡纳达语等）的文字显示</td></tr>
          <tr><td>gujr-extra</td><td>394KB</td><td>支持古吉拉特语言（Gujarati）（印度的一种官方语言，主要在古吉拉特邦使用）的文字显示</td></tr>
          <tr><td>telu-extra</td><td>430KB</td><td>支持泰卢固语（Telugu）（印度的一种官方语言，主要在特伦甘纳邦和安得拉邦使用）的文字显示</td></tr>
          <tr><th>总计</th><td></td><th>61.22MB</th><td></td></tr>
        </tbody>
      </table>
    </details>

- 更好的兼容性

  - 从[OnlyOffice](https://github.com/onlyoffice/core)最新版本v7.5.1修改
  - 举例： 用MacOSX的keynote创建的pptx，OnlyOffice转换没问题，并且展示准确
    - 对比Microsoft Office/WPS打开失败、LibreOffice/OpenOffice转换pdf格式丢失（老版本也打开失败）
  - 插件移除了不常用的djvu、epub、fb2格式
    <details><summary>展开查看支持矩阵</summary>

      - [官方文档](https://api.onlyoffice.com/editors/conversionapi#text-matrix)
      <table>
        <tbody>
          <tr><th rowspan="2">文字文档</th><th colspan="16">输出格式</th></tr>
          <tr><td>bmp</td><td>docm</td><td>docx</td><td>docxf</td><td>dotm</td><td>dotx</td><td>gif</td><td>html</td><td>jpg</td><td>odt</td><td>ott</td><td>pdf</td><td>pdfa</td><td>png</td><td>rtf</td><td>txt</td></tr>
          <tr><td>doc</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>docm</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>docx</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>docxf</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>dot</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>dotm</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>dotx</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>fodt</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>htm</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>html</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>mht</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>mhtml</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>odt</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>ott</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>oxps</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>pdf</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>rtf</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td></tr>
          <tr><td>stw</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>sxw</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>txt</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td></tr>
          <tr><td>wps</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>wpt</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>xml</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>xps</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><th rowspan="2">表格文档</th><th colspan="13">输出格式</th></tr>
          <tr><td>bmp</td><td>csv</td><td>gif</td><td>jpg</td><td>ods</td><td>ots</td><td>pdf</td><td>pdfa</td><td>png</td><td>xlsm</td><td>xlsx</td><td>xltm</td><td>xltx</td></tr>
          <tr><td>csv</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>et</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>ett</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>fods</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>ods</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>ots</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>sxc</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>xls</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>xlsb</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>xlsm</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>xlsx</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>xlt</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>xltm</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td></tr>
          <tr><td>xltx</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td></tr>
          <tr><td>xml</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><th rowspan="2">演示文档</th><th colspan="14">输出格式</th></tr>
          <tr><td>bmp</td><td>gif</td><td>jpg</td><td>odp</td><td>otp</td><td>pdf</td><td>pdfa</td><td>png</td><td>potm</td><td>potx</td><td>ppsm</td><td>ppsx</td><td>pptm</td><td>pptx</td></tr>
          <tr><td>dps</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>dpt</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>fodp</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>odp</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>otp</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>pot</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>potm</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>potx</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>pps</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>ppsm</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>ppsx</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>ppt</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
          <tr><td>pptm</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td><td>🆗</td></tr>
          <tr><td>pptx</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>=</td></tr>
          <tr><td>sxi</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td><td>🆗</td></tr>
        </tbody>
      </table>
    </details>

## 关于cad2x

#### 部署方法：

- 把[back/cad2x](https://github.com/orcastor/addon-previewer/tree/main/back/cad2x)下的`common`和对应平台（`linux_arm64`/`linux_x64`）目录下的`cad2x`文件都拷贝到`/opt/cad2x`下

#### 方案优势：

- 更好的兼容性

  - 从[LibreCAD](https://github.com/LibreCAD/LibreCAD)最新代码（commit-id：0601535）裁剪
  - 修复编码错误导致乱码问题
  - 支持自动识别横向纵向
  - 支持复用系统和x2t的ttf字体

- 更小的体积

  - 去除了GUI部分（界面渲染和跨平台）
  - 裁剪了大量无用依赖库和代码
  - 最终二进制文件仅2.9MB