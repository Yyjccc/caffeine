package core

type SystemInfo struct {
	FileRoot      string            `json:"fileRoot"`
	CurrentDir    string            `json:"currentDir"`
	CurrentUser   string            `json:"currentUser"`
	ProcessArch   int               `json:"processArch"` //操作系统位数
	TempDirectory string            `json:"tempDirectory"`
	IpList        []string          `json:"ipList"`
	Os            OSInfo            `json:"os"`
	Env           map[string]string `json:"env"`
}

// 操作系统信息
type OSInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
}

type JavaInfo struct {
	RuntimeName string //运行环境名称
	VmVersion   string
	VmName      string
	Home        string
	LibraryInfo string
	Version     string
}

//系统信息
//FileRoot : /;
//CurrentDir : /usr/local/nacos/nacos/bin/
//CurrentUser : root
//ProcessArch : 64
//TempDirectory : /tmp/
//RealFile : /tmp/tomcat-docbase.8848.4493972116701073386/
//OsInfo : os.name: Linux os.version: 4.19.91-28.2.an8.x86_64 os.arch: amd64
//IPList : [192.168.122.1,fe80:0:0:0:250:56ff:fe86:fc5d%ens192,49.123.1.92,0:0:0:0:0:0:0:1%lo,127.0.0.1]
//java.runtime.name : OpenJDK Runtime Environment
//java.protocol.handler.pkgs : org.springframework.boot.loader
//sun.boot.library.path : /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/amd64
//java.vm.version : 25.422-b05
//java.vm.vendor : Red Hat, Inc.
//java.vendor.url : https://www.redhat.com/
//path.separator : :
//java.vm.name : OpenJDK 64-Bit Server VM
//file.encoding.pkg : sun.io
//nacos.mode : stand alone
//user.country : CN
//sun.java.launcher : SUN_STANDARD
//sun.os.patch.level : unknown
//nacos.home : /usr/local/nacos/nacos
//PID : 10301
//java.vm.specification.name : Java Virtual Machine Specification
//user.dir : /usr/local/nacos/nacos/bin
//java.runtime.version : 1.8.0_422-b05
//java.awt.graphicsenv : sun.awt.X11GraphicsEnvironment
//java.endorsed.dirs : /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/endorsed
//os.arch : amd64
//CONSOLE_LOG_CHARSET : UTF-8
//java.io.tmpdir : /tmp
//line.separator :
//
//java.vm.specification.vendor : Oracle Corporation
//os.name : Linux
//nacos.local.ip : 49.123.1.92
//FILE_LOG_CHARSET : UTF-8
//sun.jnu.encoding : UTF-8
//spring.beaninfo.ignore : true
//java.library.path : /usr/java/packages/lib/amd64:/usr/lib64:/lib64:/lib:/usr/lib
//nacos.member.list :
//java.specification.name : Java Platform API Specification
//java.class.version : 52.0
//sun.management.compiler : HotSpot 64-Bit Tiered Compilers
//os.version : 4.19.91-28.2.an8.x86_64
//user.home : /root
//catalina.useNaming : false
//user.timezone : Asia/Shanghai
//java.awt.printerjob : sun.print.PSPrinterJob
//file.encoding : UTF-8
//java.specification.version : 1.8
//catalina.home : /usr/local/nacos/nacos/bin
//java.class.path : /usr/local/nacos/nacos/target/nacos-server.jar
//user.name : root
//loader.path : /usr/local/nacos/nacos/plugins,/usr/local/nacos/nacos/plugins/health,/usr/local/nacos/nacos/plugins/cmdb,/usr/local/nacos/nacos/plugins/selector
//java.vm.specification.version : 1.8
//sun.java.command : /usr/local/nacos/nacos/target/nacos-server.jar --spring.config.additional-location=file:/usr/local/nacos/nacos/conf/ --logging.config=/usr/local/nacos/nacos/conf/nacos-logback.xml --server.max-http-header-size=524288 nacos.nacos
//java.home : /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre
//sun.arch.data.model : 64
//user.language : zh
//java.specification.vendor : Oracle Corporation
//awt.toolkit : sun.awt.X11.XToolkit
//java.vm.info : mixed mode
//java.version : 1.8.0_422
//java.ext.dirs : /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/jre/lib/ext:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/ext
//sun.boot.class.path : /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/resources.jar:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/rt.jar:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/sunrsasign.jar:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/jsse.jar:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/jce.jar:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/charsets.jar:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/jfr.jar:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/classes
//java.awt.headless : true
//java.vendor : Red Hat, Inc.
//catalina.base : /usr/local/nacos/nacos/bin
//java.specification.maintenance.version : 5
//com.zaxxer.hikari.pool_number : 1
//nacos.standalone : true
//file.separator : /
//java.vendor.url.bug : https://access.redhat.com/support/cases/
//sun.io.unicode.encoding : UnicodeLittle
//sun.cpu.endian : little
//nacos.function.mode : All
//sun.cpu.isalist :
//PATH : /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin:/root/bin:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/bin
//JAVA : /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/bin/java
//HISTCONTROL : ignoredups
//XDG_DATA_DIRS : /root/.local/share/flatpak/exports/share:/var/lib/flatpak/exports/share:/usr/local/share:/usr/share
//HISTSIZE : 1000
//MODE : standalone
//JAVA_HOME : /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre
//BASH_FUNC_which%% : () {  ( alias;
//eval ${which_declare} ) | /usr/bin/which --tty-only --read-alias --read-functions --show-tilde --show-dot $@
//}
//TERM : xterm
//FUNCTION_MODE : all
//DBUS_SESSION_BUS_ADDRESS : unix:abstract=/tmp/dbus-pYxgDSWnmH,guid=61ef1f761c73d164016d069267171077
//LANG : zh_CN.UTF-8
//XDG_SESSION_ID : 7
//DISPLAY : localhost:10.0
//MAIL : /var/spool/mail/root
//SERVER : nacos-server
//LOGNAME : root
//which_declare : declare -f
//PWD : /usr/local/nacos/nacos/bin
//_ : /usr/bin/nohup
//BASE_DIR : /usr/local/nacos/nacos
//CUSTOM_SEARCH_LOCATIONS : file:/usr/local/nacos/nacos/conf/
//LESSOPEN : ||/usr/bin/lesspipe.sh %s
//SHELL : /bin/bash
//EMBEDDED_STORAGE :
//GDK_BACKEND : x11
//SSH_TTY : /dev/pts/0
//SSH_CLIENT : 10.102.176.170 48038 22
//OLDPWD : /usr/local/nacos/nacos/conf
//USER : root
//CLASSPATH : .:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/jre/lib:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib:/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.422.b05-2.0.2.an8.x86_64/jre/lib/tools.jar
//SSH_ASKPASS : /usr/libexec/openssh/gnome-ssh-askpass
//MEMBER_LIST :
//SSH_CONNECTION : 10.102.176.170 48038 49.123.1.92 22
//HOSTNAME : AnolisOS86
//XDG_RUNTIME_DIR : /run/user/0
//LS_COLORS : rs=0:di=01;34:ln=01;36:mh=00:pi=40;33:so=01;35:do=01;35:bd=40;33;01:cd=40;33;01:or=40;31;01:mi=01;05;37;41:su=37;41:sg=30;43:ca=30;41:tw=30;42:ow=34;42:st=37;44:ex=01;32:*.tar=01;31:*.tgz=01;31:*.arc=01;31:*.arj=01;31:*.taz=01;31:*.lha=01;31:*.lz4=01;31:*.lzh=01;31:*.lzma=01;31:*.tlz=01;31:*.txz=01;31:*.tzo=01;31:*.t7z=01;31:*.zip=01;31:*.z=01;31:*.dz=01;31:*.gz=01;31:*.lrz=01;31:*.lz=01;31:*.lzo=01;31:*.xz=01;31:*.zst=01;31:*.tzst=01;31:*.bz2=01;31:*.bz=01;31:*.tbz=01;31:*.tbz2=01;31:*.tz=01;31:*.deb=01;31:*.rpm=01;31:*.jar=01;31:*.war=01;31:*.ear=01;31:*.sar=01;31:*.rar=01;31:*.alz=01;31:*.ace=01;31:*.zoo=01;31:*.cpio=01;31:*.7z=01;31:*.rz=01;31:*.cab=01;31:*.wim=01;31:*.swm=01;31:*.dwm=01;31:*.esd=01;31:*.jpg=01;35:*.jpeg=01;35:*.mjpg=01;35:*.mjpeg=01;35:*.gif=01;35:*.bmp=01;35:*.pbm=01;35:*.pgm=01;35:*.ppm=01;35:*.tga=01;35:*.xbm=01;35:*.xpm=01;35:*.tif=01;35:*.tiff=01;35:*.png=01;35:*.svg=01;35:*.svgz=01;35:*.mng=01;35:*.pcx=01;35:*.mov=01;35:*.mpg=01;35:*.mpeg=01;35:*.m2v=01;35:*.mkv=01;35:*.webm=01;35:*.ogm=01;35:*.mp4=01;35:*.m4v=01;35:*.mp4v=01;35:*.vob=01;35:*.qt=01;35:*.nuv=01;35:*.wmv=01;35:*.asf=01;35:*.rm=01;35:*.rmvb=01;35:*.flc=01;35:*.avi=01;35:*.fli=01;35:*.flv=01;35:*.gl=01;35:*.dl=01;35:*.xcf=01;35:*.xwd=01;35:*.yuv=01;35:*.cgm=01;35:*.emf=01;35:*.ogv=01;35:*.ogx=01;35:*.aac=01;36:*.au=01;36:*.flac=01;36:*.m4a=01;36:*.mid=01;36:*.midi=01;36:*.mka=01;36:*.mp3=01;36:*.mpc=01;36:*.ogg=01;36:*.ra=01;36:*.wav=01;36:*.oga=01;36:*.opus=01;36:*.spx=01;36:*.xspf=01;36:
//SHLVL : 2
//HOME : /root
