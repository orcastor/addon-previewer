<p align="center">
  <a href="https://orcastor.github.io/doc/">
    <img src="https://orcastor.github.io/doc/logo.svg">
  </a>
</p>

<p align="center"><strong>OrcaS 文档预览插件</strong></p>

<p align="center">—— 基于h5，纯离线本地版，效果/兼容性更好，可自行适配存储 ——</p>

|类型|进度|方案|格式|备注|
|-|-|-|-|-|
|PDF文档|✅|[pdf.js](https://github.com/mozilla/pdf.js) / [vue-pdf-embed](https://github.com/hrynko/vue-pdf-embed)|pdf|支持有密码文件|
|Word文档|✅|[docxjs](https://www.npmjs.com/package/docx-preview) + [x2t](https://github.com/ONLYOFFICE/core)转换成docx格式|docx、doc、pages等|\*字体问题；有密码文档问题|
|Excel文档|✅|[luckyexcel](https://github.com/dream-num/Luckyexcel) + [x2t](https://github.com/ONLYOFFICE/core)转换成xlsx格式|xlsx、xls、csv、numbers等|\*字体问题；有密码文档问题|
|PowerPoint文档|✅|[x2t](https://github.com/ONLYOFFICE/core)转换成pdf格式|pptx、ppt、key等|\*字体问题；有密码文档问题|
|CAD文档||工具转换成pdf|dwg、dwt等|
|文本文档|✅|转docx / 富文本编辑器|txt、json/toml/yml/config/xml等|enca/file检测编码防止中文乱码|
|代码文档||md代码染色 / WebIDE|cpp、c、h、java、py、go、php、js、html、css等|enca/file检测编码防止中文乱码|
|图片|✅|vips / ImagicMagick / GraphicsMagick等转换缩略图|png、jpg、gif、bmp、ico等|gif需要处理多帧；GM支持OpenMP加速|
|视频||ffmpeg转换成HLS(m3u8)|mp4、wmv、mkv等|fps和码率；可以尝试h.265；支持GPU加速（OpenCL/Vulkan）|
|音频||[aPlayer](https://github.com/DIYgod/APlayer)/[cPlayer](https://github.com/MoePlayer/cPlayer)|mp3,wma,wav,ape,flac,ogg,aac等|支持匹配封面、lrc歌词文件|
|压缩包||写到临时目录解压|rar、zip、7z等|有密码的文档|
|备份||imobax & abe|iOS备份目录，安卓备份ab文件等|有密码的备份|
|图标||图片 / PE格式解析|desktop.ini、dmg、exe、apk、*.app目录等|

### 部署方法：

- 把[back/x2t](https://github.com/orcastor/addon-previewer/tree/main/back/x2t)下的`common`和对应平台（`linux_arm64`/`linux_x64`）目录下的文件都拷贝到`/opt/x2t`下
- 执行`sh allfontsgen.sh`安装字体文件、生成字体列表

其他：
- 如果使用docker，需要配置`ORCAS_DOCKER_EXEC="docker exec -i <container_id>"`环境变量，其中`container_id`为容器的id值
- 打包后的文件和webapp的放置到一起：
`ln -s $(addon-previewer)/front/dist $(webapp)/dist/prvw`
