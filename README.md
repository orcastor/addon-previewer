<p align="center">
  <a href="https://orcastor.github.io/doc/">
    <img src="https://orcastor.github.io/doc/logo.svg">
  </a>
</p>

<p align="center"><strong>OrcaS 在线预览插件</strong></p>


|类型|方案|格式|备注|
|-|-|-|-|
|PDF文档|[pdf.js](https://github.com/mozilla/pdf.js) / [vue-pdf-embed](https://github.com/hrynko/vue-pdf-embed)|pdf|支持有密码文件|
|Office文档|LibreOffice转换成pdf / [Office365](https://www.microsoft.com/en-us/microsoft-365/blog/2013/04/10/office-web-viewer-view-office-documents-in-a-browser/?eu=true) / [CKEditor](https://ckeditor.com/docs/ckeditor5/latest/examples/builds/document-editor.html) / [PPTist](https://github.com/pipipi-pikachu/PPTist) |docx、pptx、xlsx等|LibreOffice字体问题；CKEditor / PPTist支持编辑；有密码文档问题|
|CAD文档|工具转换成pdf|dwg、dwt等|
|文本文档|转pdf / 富文本编辑器|txt、json/toml/yml/config/xml等|enca/file检测编码防止中文乱码|
|代码文档|md代码染色 / WebIDE|cpp、c、h、java、py、go、php、js、html、css等|enca/file检测编码防止中文乱码|
|图片|vips / ImagicMagick / Graphics等转换缩略图|png、jpg、gif、bmp、ico等|gif需要处理多帧|
|视频|ffmpeg转换成HLS(m3u8)|mp4、wmv、mkv等|fps和码率；可以尝试h.265|
|压缩包|本身库支持|rar、zip、7z等|有密码的文档|
|图标|图片 / PE格式解析|desktop.ini、dmg、exe、*.app目录等|
