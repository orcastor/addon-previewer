<p align="center">
  <a href="https://orcastor.github.io/doc/">
    <img src="https://orcastor.github.io/doc/logo.svg">
  </a>
</p>

<p align="center"><strong>OrcaS 文档预览插件</strong></p>

<p align="center">—— 基于h5，纯离线本地版，效果/兼容性更好，可自行适配存储 ——</p>

|类型|进度|方案|格式|备注|
|-|-|-|-|-|
|PDF|✅|[vue-pdf-embed](https://github.com/hrynko/vue-pdf-embed)|pdf|支持有密码文件|
|文字|✅|[docxjs](https://www.npmjs.com/package/docx-preview)+[x2t](https://github.com/ONLYOFFICE/core)转换成docx格式|docx、doc、pages、wps等|\*字体问题；有密码文档问题|
|表格|✅|[luckyexcel](https://github.com/dream-num/Luckyexcel) + [x2t](https://github.com/ONLYOFFICE/core)转换成xlsx格式|xlsx、xls、csv、numbers、et等|\*字体问题；有密码文档问题|
|演示|✅|[x2t](https://github.com/ONLYOFFICE/core)转换成pdf格式|pptx、ppt、key等|\*字体问题；有密码文档问题|
|设计||[LibreCAD](https://github.com/LibreCAD/LibreCAD)转换成pdf/png/svg|dwg、dwt、dxf等|
|文本|✅|转docx / 富文本编辑器|txt、json/toml/yml/config/xml等|enca/file检测编码防止中文乱码|
|代码||md代码染色 / WebIDE|cpp、c、h、java、py、go、php、js、html、css等|enca/file检测编码防止中文乱码|
|图片|✅|vips / ImagicMagick / GraphicsMagick等转换缩略图|png、jpg、gif、bmp、ico等|gif需要处理多帧；GM支持OpenMP加速|
|视频||ffmpeg转换成HLS(m3u8)|mp4、wmv、mkv等|fps和码率；可以尝试h.265；支持GPU加速（OpenCL/Vulkan）|
|音频||[aPlayer](https://github.com/DIYgod/APlayer)/[cPlayer](https://github.com/MoePlayer/cPlayer)|mp3,wma,wav,ape,flac,ogg,aac等|支持匹配封面、lrc歌词文件|
|存档||写到临时目录解压|rar、zip、7z等|有密码的文档|
|备份||imobax & abe|iOS备份目录，安卓备份ab文件等|有密码的备份|
|图标||图片 / PE格式解析|desktop.ini、dmg、exe、apk、*.app目录等|

## 关于onlyoffice/x2t

#### 部署方法：

- 把[back/x2t](https://github.com/orcastor/addon-previewer/tree/main/back/x2t)下的`common`和对应平台（`linux_arm64`/`linux_x64`）目录下的文件都拷贝到`/opt/x2t`（换其他路径执行下面的脚本即可）下
- 执行`sh allfontsgen.sh`安装字体文件、生成字体列表

#### 注意事项：
- 如果转换插件在docker镜像中部署，需要配置`ORCAS_DOCKER_EXEC="docker exec -i <container_id>"`环境变量，其中`container_id`为容器的id值
- 打包后的文件和webapp的放置到一起：
`ln -s $(addon-previewer)/front/dist $(webapp)/dist/prvw`

#### 方案优势：

- 更好的性能

  - 纯C++开发，执行更快
  - 原生支持同时并发转换多个文件
  - 秒级启动，不需要常驻

- 更小的体积（完整支持仅215MB左右）
  - 二进制文件37MB
  - 系统字体文件下载39MB
  - 常见字体文件附带39MB
  - js文件100MB

- 更好的兼容性

  - Case举例： keynote创建的pptx格式，LibreOffice/OpenOffice转换pdf失败，onlyoffice没问题，并且字体更准确
  - 插件中仅移除了不常用的epub和fb2格式
  - [点我查看完整文档](https://api.onlyoffice.com/editors/conversionapi#text-matrix)
        <table>
          <tbody>
            <tr><th rowspan="2">文字文档</th><th colspan="18">输出格式</th></tr>
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

