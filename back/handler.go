package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/orcastor/iwork-converter/iwork2html"
	"github.com/orcastor/orcas/core"
	"github.com/orcastor/orcas/rpc/util"
	"github.com/orcastor/orcas/sdk"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/elog"
	"github.com/mholt/archiver/v3"
)

var ORCAS_DOCKER_EXEC = os.Getenv("ORCAS_DOCKER_EXEC")
var ORCAS_CACHE = os.Getenv("ORCAS_CACHE")

func init() {
	if ORCAS_CACHE == "" {
		ORCAS_CACHE = "/tmp/orcas_cache"
	}
	os.MkdirAll(ORCAS_CACHE, 0766)
}

var hanlder = core.NewLocalHandler()

var cadTypes = map[string]bool{
	"dwg": true,
	"dxf": true,
}

var iworkTypes = map[string]bool{
	"pages":   true,
	"numbers": true,
	"key":     true,
}

var docConvTypes = map[string]string{
	"dwg": "png",
	"dxf": "png",

	"pages":   "html",
	"numbers": "html",
	"key":     "html",

	"dsp":  "pdf",
	"ppt":  "pdf",
	"pptx": "pdf",

	"et":  "xlsx",
	"csv": "xlsx",
	"xls": "xlsx",

	"wps":    "docx",
	"doc":    "docx",
	"txt":    "docx",
	"json":   "docx",
	"toml":   "docx",
	"yaml":   "docx",
	"xml":    "docx",
	"config": "docx",
}

func get(ctx *gin.Context) {
	bktID, _ := strconv.ParseInt(ctx.Query("b"), 10, 64)
	id, _ := strconv.ParseInt(ctx.Query("i"), 10, 64)
	langID, _ := strconv.ParseInt(ctx.Query("l"), 10, 64) // 预留的语言id

	from := strings.ToLower(ctx.Query("t")) // from无法注入，不在白名单会直接返回
	to := docConvTypes[from]
	if to == "" {
		// 不需要转换格式，那就直接写到http
		if err := writeTo(ctx, bktID, id, ctx.Writer, true); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
		}
		return
	}

	fromPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.%s", id, from)) // id无法注入，强制转成数字
	toPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.%s", id, to))

	// 先看转换后文件生成了没有
	if st, err := os.Stat(toPath); err == nil && st.Size() > 0 {
		goto READ_TO_FILE
	}

	// 下载到临时文件
	if err := download(ctx, bktID, id, fromPath); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	// 转换格式
	if cadTypes[from] {
		if err := cad2xConv(fromPath, toPath); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
			return
		}
	} else if iworkTypes[from] {
		if err := iwork2htmlConv(fromPath, toPath); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
			return
		}
	} else {
		if err := x2tConv(fromPath, toPath); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
			return
		}
	}

	// 删除临时文件
	os.Remove(fromPath)

READ_TO_FILE:
	ft, err := os.Open(toPath)
	if err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}
	defer ft.Close()

	if _, err = io.Copy(ctx.Writer, bufio.NewReader(ft)); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}
}

func langID2CodePage(langID int64) string {
	switch langID {
	case 1054:
		return "ANSI_874" // 泰语
	case 1041:
		return "ANSI_932" // 日语
	case 2052:
		return "ANSI_936" // 简体中文
	case 1042:
		return "ANSI_949" // 韩语
	case 1028:
		return "ANSI_950" // 繁体中文
	case 1029:
		return "ANSI_1250" // 中欧: 保加利亚语
	case 1049, 1062, 1087, 1104, 1121, 2072, 2073:
		return "ANSI_1251" // 西里尔语: 俄语, 白俄罗斯语, 马其顿语, 塞尔维亚语, 斯洛文尼亚语, 乌兹别克语, 波斯尼亚语
	case 1032:
		return "ANSI_1253" // 希腊语
	case 1055:
		return "ANSI_1254" // 土耳其语
	case 1037:
		return "ANSI_1255" // 希伯来语
	case 1025:
		return "ANSI_1256" // 阿拉伯语
	case 1066:
		return "ANSI_1257" // 波罗的海语
	case 1068:
		return "ANSI_1258" // 越南语
	}
	return "ANSI_1252" // 西欧: 英语
}

func cad2xConv(fromPath, toPath string, langID int64) error {
	var out bytes.Buffer
	cmds := append(strings.Split(ORCAS_DOCKER_EXEC, " "),
		"/opt/cad2x/cad2pdf",
		fromPath,
		"-o", toPath,
		"-ac",
		"-e", langID2CodePage(langID),
		"-f", "simsun",
		"-l", "/opt/x2t/core-fonts")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		elog.Errorf("cad2xConv error: %+v", out.String())
	}
	return err
}

func iwork2htmlConv(fromPath, toPath string) error {
	if err := iwork2html.Convert(fromPath, toPath); err != nil {
		elog.Errorf("iwork2htmlConv error: %+v", err)
		return err
	}
	return nil
}

func x2tConv(fromPath, toPath string) error {
	var out bytes.Buffer
	cmds := append(strings.Split(ORCAS_DOCKER_EXEC, " "),
		"/opt/x2t/x2t",
		fromPath,
		toPath)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		elog.Errorf("x2tConv error: %+v", out.String())
	}
	return err
}

func ico(ctx *gin.Context) {
	bktID, _ := strconv.ParseInt(ctx.Query("b"), 10, 64)
	id, _ := strconv.ParseInt(ctx.Query("i"), 10, 64)
	langID, _ := strconv.ParseInt(ctx.Query("l"), 10, 64) // 预留的语言id

	fileName := strings.ToLower(ctx.Query("n"))

	switch filepath.Ext(fileName)[1:] {
	// 文件
	// *.dmg、*.exe、*.apk
	case "dmg":
		/*
			在 macOS 的 DMG（Disk Image）文件中，图标文件通常存放在 .DS_Store 文件中以及 .VolumeIcon.icns 文件中。

			.DS_Store 文件：.DS_Store 文件是 macOS 用来存储文件夹的自定义显示属性的隐藏文件。它可能包含与文件夹关联的图标信息。您可以将自定义图标应用于文件夹，并将其存储在 .DS_Store 文件中。在创建 DMG 文件时，您可以使用命令行工具或可视化工具来自定义文件夹的图标，这样在挂载 DMG 时文件夹会显示相应的图标。

			.VolumeIcon.icns 文件：.VolumeIcon.icns 文件是存储在 DMG 文件中的磁盘图标文件。您可以通过创建一个包含所需图标的 .icns 文件，并将其命名为 .VolumeIcon.icns，然后将其添加到 DMG 文件中。这样，当用户挂载 DMG 文件时，磁盘会显示指定的图标。

			.DS_Store 是 macOS 操作系统中用于存储目录自定义属性和图标布局信息的隐藏文件。https://metacpan.org/dist/Mac-Finder-DSStore/view/DSStoreFormat.pod
		*/
	case "exe":
		/*
			在 Windows 中，当匹配一个 EXE 文件的图标时，通常会选择其中的一个资源，这个资源通常是包含在 PE 文件中的一组图标资源中的一个。选择的资源不一定是具有最小 ID 的资源，而是根据一些规则进行选择。
			具体来说，Windows 会根据以下几个因素来选择图标资源：
			图标大小：Windows 会根据显示图标的大小来选择最合适的资源。这意味着如果需要的是小尺寸图标，Windows 会选择包含小尺寸图标的资源。
			分辨率：如果图标资源包含不同分辨率的图标，Windows 会选择与显示器分辨率最匹配的图标。
			语言：Windows 还会考虑资源中的语言信息，优先选择与当前系统语言匹配的资源。
			图标类型：有时候 PE 文件中可能包含多个类型的图标资源，如 16x16 和 32x32 大小的图标，Windows 会根据需要选择适合的类型。
			综上所述，选择的图标资源并不一定是具有最小 ID 的资源，而是根据上述因素来决定的。因此，在制作 PE 文件时，确保包含了适合不同显示情况的图标资源是很重要的。
		*/
	case "apk":
		/*
			在 Android 应用程序中，应用程序的图标通常是在 AndroidManifest.xml 文件中定义的。图标资源通常存储在 res 目录下的 drawable 目录中，并在 AndroidManifest.xml 文件中通过 <application> 元素的 android:icon 属性引用。
			例如：

			<application
				android:icon="@drawable/ic_launcher"
				android:label="@string/app_name">
				...
			</application>

			在上面的示例中，@drawable/ic_launcher 引用了应用程序的图标资源，该图标资源存储在 drawable 目录下，并命名为 ic_launcher。

			APK 文件实际上是一种 ZIP 压缩文件格式，其中包含了 Android 应用程序的所有资源和代码。APK 文件的目录组织通常如下所示：
			myapp.apk
			│
			├── AndroidManifest.xml
			├── META-INF/
			│   └── CERT.RSA
			│   └── CERT.SF
			│   └── MANIFEST.MF
			├── res/
			│   ├── drawable/
			│   │   └── icon.png (应用程序图标)
			│   │   └── ...
			│   ├── layout/
			│   │   └── activity_main.xml
			│   │   └── ...
			│   ├── values/
			│   │   └── strings.xml
			│   │   └── ...
			└── classes.dex (编译后的 Java 字节码文件)
			在这个目录结构中，res/drawable/ 目录通常包含了应用程序的图标文件，通常以 .png 或 .jpg 格式存储。应用程序的图标文件通常命名为 icon.png 或 ic_launcher.png，但实际名称可以根据开发人员的喜好而有所不同。

			对于 Android 应用程序的图标文件，通常可以省略 .png 后缀，只使用图标的名称即可。例如，如果图标文件名为 icon.png，则在配置文件中引用时可以简单地写为 icon，而不需要指定完整的文件名。Android 系统会自动查找匹配的图标文件。
		*/

	// 配置文件
	// autorun.inf、desktop.ini、*.desktop(*.AppImage/*.run)
	case "ini":
		/*
			在 Windows 操作系统中，desktop.ini 文件用于自定义文件夹的外观和行为。您可以在文件夹中创建 desktop.ini 文件，并在其中指定如何显示该文件夹的图标。

			要在 desktop.ini 文件中定义图标，可以使用 IconFile 和 IconIndex 字段。下面是一个示例 desktop.ini 文件的基本结构：

			[.ShellClassInfo]
			IconFile=path\to\icon.ico
			IconIndex=0

			IconFile 字段指定要用作文件夹图标的图标文件的路径。这可以是包含图标的 .ico 文件，也可以是 .exe 或 .dll 文件，其中包含一个或多个图标资源。
			IconIndex 字段指定要在 IconFile 中使用的图标的索引。如果 IconFile 是 .ico 文件，则索引从0开始，表示图标在文件中的位置。如果 IconFile 是 .exe 或 .dll 文件，则索引表示图标资源的标识符。
			完成后，您可以将 desktop.ini 文件放置在所需文件夹中，并在 Windows 资源管理器中刷新文件夹，以查看所指定的图标。
		*/
	case "inf":
		/*
			在 Windows 系统中，autorun.inf 文件用于自定义 CD、DVD 或 USB 驱动器上的自动运行功能。您可以在 autorun.inf 文件中定义要显示的图标。以下是如何定义图标的方法：

			使用 Icon 指令：
			在 autorun.inf 文件中添加 Icon 指令，并指定要显示的图标文件的路径。图标文件可以是 .ico 格式的图标文件。

			示例：

			[AutoRun]
			Icon=path\to\icon.ico

			在这个示例中，Icon 指令指定了要显示的图标文件的路径。

			使用 DefaultIcon 指令：
			另一种定义图标的方法是使用 DefaultIcon 指令。与 Icon 指令类似，DefaultIcon 指令也用于指定要显示的图标文件的路径。

			示例：

			[AutoRun]
			DefaultIcon=path\to\icon.ico

			与 Icon 指令不同的是，DefaultIcon 指令可以同时用于指定文件和文件夹的图标。

			在这两种方法中，path\to\icon.ico 是要显示的图标文件的路径。

			完成后，将 autorun.inf 文件与您的可移动媒体（如 CD、DVD 或 USB 驱动器）一起放置，并在 Windows 系统中插入该媒体，系统会根据 autorun.inf 文件中的设置自动运行，并显示所指定的图标。
		*/
	case "desktop":
		/*
			创建包含图标和其他资源的 .desktop 文件来为 .run 文件指定图标。然后，您可以将 .run 文件与 .desktop 文件一起分发，并通过 .desktop 文件来启动 .run 文件，并在系统中显示指定的图标。

			以下是一个示例 .desktop 文件的基本结构：

			[Desktop Entry]
			Version=1.0
			Type=Application
			Name=YourApp
			Icon=/path/to/your/icon.png
			Exec=/path/to/your/run/file.run
			Terminal=false

			您需要将 Icon 字段设置为指向您要在系统中显示的图标文件的路径，并将 Exec 字段设置为指向您的 .run 文件的路径。然后，您可以将 .desktop 文件放置在系统的应用程序启动器中，用户可以通过单击该图标来运行 .run 文件，并显示指定的图标。
		*/

	// 目录
	// *.app
	case "app":
		/*
			在 macOS 中，应用程序的图标通常包含在应用程序的 .app 包中，并通过特定的文件 Info.plist 来定义。Info.plist 文件包含了有关应用程序的元数据信息，包括应用程序的名称、版本、图标等。

			要定义应用程序的图标，可以在应用程序的 Info.plist 文件中添加 CFBundleIconFile 或 CFBundleIconName 键。具体取决于图标是存储在 .app 包的根目录中的 .icns 文件，还是存储在 .app 包的资源目录中的文件。

			以下是两种常见的情况：

			如果图标存储在 .app 包的根目录中的 .icns 文件中，则可以使用 CFBundleIconFile 键指定图标文件的名称，例如：

			<key>CFBundleIconFile</key>
			<string>AppIcon.icns</string>

			如果图标存储在 .app 包的资源目录中，则可以使用 CFBundleIconName 键指定图标文件的名称，例如：

			<key>CFBundleIconName</key>
			<string>AppIcon</string>

			无论是使用 CFBundleIconFile 还是 CFBundleIconName，都要确保图标文件的名称和位置正确，并且应该位于 .app 包的根目录或资源目录中。
			如果您使用 CFBundleIconName 键来指定应用程序的图标，您可以省略 .icns 后缀。

			请注意，.icns 文件是 macOS 中常见的图标文件格式，它可以包含多个图标尺寸和位深度的图像。.icns 文件是 macOS 中用于存储应用程序图标的标准格式。
		*/
	}
}

var thumbSupport = map[string]bool{
	"csv":  true,
	"bmp":  true,
	"raw":  true,
	"jpg":  true,
	"jpeg": true,
	"jpe":  true,
	"jfif": true,
	"png":  true,
	"gif":  true,
	"tif":  true,
	"tiff": true,
	"webp": true,

	"mat":  true,
	"pbm":  true,
	"pgm":  true,
	"ppm":  true,
	"pfm":  true,
	"pnm":  true,
	"fits": true,
	"fit":  true,
	"fts":  true,
	"exr":  true,
	"hdr":  true,
	"v":    true,
	"vips": true,
}

var outTypes = map[string]bool{
	"jpg":  true,
	"png":  true,
	"gif":  true,
	"webp": true,
}

func thumbnail(ctx *gin.Context) {
	bktID, _ := strconv.ParseInt(ctx.Query("b"), 10, 64)
	id, _ := strconv.ParseInt(ctx.Query("i"), 10, 64)

	w, _ := strconv.ParseInt(ctx.Query("w"), 10, 64)
	h, _ := strconv.ParseInt(ctx.Query("h"), 10, 64)

	from := strings.ToLower(ctx.Query("t")) // from无法注入，不在白名单会直接返回
	if !thumbSupport[from] {
		// 不需要转换格式，那就直接写到http
		if err := writeTo(ctx, bktID, id, ctx.Writer, true); err != nil {
			util.AbortResponse(ctx, 100, err.Error())
		}
		return
	}

	// TODO：如果是文档格式，先用x2t获取文档缩略图

	// TODO：如果涉及隐私文件，返回不支持获取缩略图

	to := strings.ToLower(ctx.Query("nt")) // to无法注入，不在白名单会直接返回
	if !outTypes[to] {
		util.AbortResponse(ctx, 400, "not supported format")
		return
	}

	fromPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d.%s", id, from)) // id无法注入，强制转成数字
	toPath := filepath.Join(ORCAS_CACHE, fmt.Sprintf("%d_%dx%d.%s", id, w, h, to))

	// 先看转换后文件生成了没有
	if st, err := os.Stat(toPath); err == nil && st.Size() > 0 {
		goto READ_TO_FILE
	}

	// 下载到临时文件
	if err := download(ctx, bktID, id, fromPath); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	// 转换缩略图
	if err := vipsConv(fromPath, toPath, w, h); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}

	// 删除临时文件
	os.Remove(fromPath)

READ_TO_FILE:
	ft, err := os.Open(toPath)
	if err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}
	defer ft.Close()

	if _, err = io.Copy(ctx.Writer, bufio.NewReader(ft)); err != nil {
		util.AbortResponse(ctx, 100, err.Error())
		return
	}
}

func vipsConv(fromPath, toPath string, w, h int64) error {
	var out bytes.Buffer
	cmds := append(strings.Split(ORCAS_DOCKER_EXEC, " "), "/opt/vips/vipsthumbnail", fromPath,
		"--size", fmt.Sprintf("%dx%d", w, h),
		"--smartcrop", "attention",
		"-o", toPath+"[keep=none]")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		elog.Errorf("vipsConv error: %+v", out.String())
	}
	return err
}

func download(ctx *gin.Context, bktID, id int64, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	defer b.Flush()

	return writeTo(ctx, bktID, id, b, false)
}

func writeTo(ctx *gin.Context, bktID, id int64, writer io.Writer, direct bool) error {
	o, err := hanlder.Get(ctx.Request.Context(), bktID, []int64{id})
	if err != nil {
		return err
	}

	if len(o) <= 0 || o[0].Type != core.OBJ_TYPE_FILE {
		return errors.New("not file")
	}

	dataID := o[0].DataID
	// 如果不是首版本
	if dataID == 0 {
		os, _, _, err := hanlder.List(ctx.Request.Context(), bktID, id, core.ListOptions{
			Type:  core.OBJ_TYPE_VERSION,
			Count: 1,
			Order: "-id",
		})
		if err != nil {
			return err
		}
		if len(os) < 1 {
			return errors.New("no version")
		}
		dataID = os[0].DataID
	}

	var d *core.DataInfo
	if dataID == core.EmptyDataID {
		d = core.EmptyDataInfo()
	} else {
		d, err = hanlder.GetDataInfo(ctx.Request.Context(), bktID, dataID)
		if err != nil {
			return err
		}
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(o[0].Name)))

	acceptEncoding := ctx.Request.Header["Accept-Encoding"]

	var decmpr archiver.Decompressor
	if d.Kind&core.DATA_CMPR_MASK != 0 {
		if d.Kind&core.DATA_CMPR_SNAPPY != 0 {
			decmpr = &archiver.Snappy{}
		} else if d.Kind&core.DATA_CMPR_ZSTD != 0 {
			decmpr = &archiver.Zstd{}
		} else if d.Kind&core.DATA_CMPR_GZIP != 0 {
			if direct && len(acceptEncoding) > 0 && strings.Contains(acceptEncoding[0], "gzip") {
				// 如果浏览器支持gzip，直接返回原始数据
				decmpr = &sdk.DummyArchiver{}
				ctx.Header("Content-Encoding", "gzip")
			} else {
				decmpr = &archiver.Gz{}
			}
		} else if d.Kind&core.DATA_CMPR_BR != 0 {
			if direct && len(acceptEncoding) > 0 && strings.Contains(acceptEncoding[0], "br") {
				// 如果浏览器支持br，直接返回原始数据
				decmpr = &sdk.DummyArchiver{}
				ctx.Header("Content-Encoding", "br")
			} else {
				decmpr = &archiver.Brotli{}
			}
		}
	} else {
		decmpr = &sdk.DummyArchiver{}
	}

	return decmpr.Decompress(sdk.NewDataReader(ctx.Request.Context(), hanlder, bktID, d, ""), writer)
}
