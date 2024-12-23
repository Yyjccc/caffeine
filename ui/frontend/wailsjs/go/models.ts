export namespace client {
	
	export class Connection {
	    localAddress: string;
	    remoteAddress: string;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new Connection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.localAddress = source["localAddress"];
	        this.remoteAddress = source["remoteAddress"];
	        this.state = source["state"];
	    }
	}
	export class NetworkInterface {
	    name: string;
	    ipAddress: string;
	    macAddress: string;
	    bytesReceived: number;
	    bytesSent: number;
	
	    static createFrom(source: any = {}) {
	        return new NetworkInterface(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.ipAddress = source["ipAddress"];
	        this.macAddress = source["macAddress"];
	        this.bytesReceived = source["bytesReceived"];
	        this.bytesSent = source["bytesSent"];
	    }
	}
	export class Port {
	    port: number;
	    protocol: string;
	    process: string;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new Port(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.protocol = source["protocol"];
	        this.process = source["process"];
	        this.state = source["state"];
	    }
	}
	export class ShellEntry {
	    ID: number;
	    Location: string;
	    ShellType: string;
	    IP: string;
	    CreateTime: string;
	    UpdateTime: string;
	    URL: string;
	    Note: string;
	    Password: string;
	    Encoding: string;
	    Status: number;
	
	    static createFrom(source: any = {}) {
	        return new ShellEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Location = source["Location"];
	        this.ShellType = source["ShellType"];
	        this.IP = source["IP"];
	        this.CreateTime = source["CreateTime"];
	        this.UpdateTime = source["UpdateTime"];
	        this.URL = source["URL"];
	        this.Note = source["Note"];
	        this.Password = source["Password"];
	        this.Encoding = source["Encoding"];
	        this.Status = source["Status"];
	    }
	}
	export class SystemMetric {
	    cpu: number;
	    memory: number;
	    time: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemMetric(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu = source["cpu"];
	        this.memory = source["memory"];
	        this.time = source["time"];
	    }
	}

}

export namespace core {
	
	export class JavaInfo {
	    runtimeName: string;
	    vmVersion: string;
	    vmName: string;
	    home: string;
	    libraryPath: string;
	    version: string;
	    classPath: string;
	    tempDir: string;
	    extDirs: string[];
	    securityPolicy: string;
	
	    static createFrom(source: any = {}) {
	        return new JavaInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.runtimeName = source["runtimeName"];
	        this.vmVersion = source["vmVersion"];
	        this.vmName = source["vmName"];
	        this.home = source["home"];
	        this.libraryPath = source["libraryPath"];
	        this.version = source["version"];
	        this.classPath = source["classPath"];
	        this.tempDir = source["tempDir"];
	        this.extDirs = source["extDirs"];
	        this.securityPolicy = source["securityPolicy"];
	    }
	}
	export class NetIfaceInfo {
	    name: string;
	    mac: string;
	    ipAddresses: string[];
	    mtu: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new NetIfaceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.mac = source["mac"];
	        this.ipAddresses = source["ipAddresses"];
	        this.mtu = source["mtu"];
	        this.status = source["status"];
	    }
	}
	export class OSInfo {
	    name: string;
	    version: string;
	    arch: string;
	    kernel: string;
	    distribution: string;
	    timezone: string;
	    language: string;
	    lastBoot: string;
	
	    static createFrom(source: any = {}) {
	        return new OSInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.arch = source["arch"];
	        this.kernel = source["kernel"];
	        this.distribution = source["distribution"];
	        this.timezone = source["timezone"];
	        this.language = source["language"];
	        this.lastBoot = source["lastBoot"];
	    }
	}
	export class ProcessInfo {
	    pid: number;
	    ppid: number;
	    startTime: string;
	    cmdLine: string;
	    workingDir: string;
	    owner: string;
	    memory: number;
	    openFiles: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.ppid = source["ppid"];
	        this.startTime = source["startTime"];
	        this.cmdLine = source["cmdLine"];
	        this.workingDir = source["workingDir"];
	        this.owner = source["owner"];
	        this.memory = source["memory"];
	        this.openFiles = source["openFiles"];
	    }
	}
	export class SecurityInfo {
	    selinuxEnabled: boolean;
	    appArmor: boolean;
	    capabilities: string[];
	    sudoersConfig: string;
	    sshKeys: string[];
	    cronJobs: string[];
	    services: string[];
	
	    static createFrom(source: any = {}) {
	        return new SecurityInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.selinuxEnabled = source["selinuxEnabled"];
	        this.appArmor = source["appArmor"];
	        this.capabilities = source["capabilities"];
	        this.sudoersConfig = source["sudoersConfig"];
	        this.sshKeys = source["sshKeys"];
	        this.cronJobs = source["cronJobs"];
	        this.services = source["services"];
	    }
	}
	export class SystemInfo {
	    ID: number;
	    fileRoot: string;
	    currentDir: string;
	    currentUser: string;
	    processArch: number;
	    tempDirectory: string;
	    hostname: string;
	    ipList: string[];
	    netIfaces: {[key: string]: NetIfaceInfo};
	    listeningPorts: string[];
	    os: OSInfo;
	    env: {[key: string]: any};
	    processInfo: ProcessInfo;
	    java: JavaInfo;
	    webRoot: string;
	    configPaths: string[];
	    logPaths: string[];
	    security: SecurityInfo;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.fileRoot = source["fileRoot"];
	        this.currentDir = source["currentDir"];
	        this.currentUser = source["currentUser"];
	        this.processArch = source["processArch"];
	        this.tempDirectory = source["tempDirectory"];
	        this.hostname = source["hostname"];
	        this.ipList = source["ipList"];
	        this.netIfaces = this.convertValues(source["netIfaces"], NetIfaceInfo, true);
	        this.listeningPorts = source["listeningPorts"];
	        this.os = this.convertValues(source["os"], OSInfo);
	        this.env = source["env"];
	        this.processInfo = this.convertValues(source["processInfo"], ProcessInfo);
	        this.java = this.convertValues(source["java"], JavaInfo);
	        this.webRoot = source["webRoot"];
	        this.configPaths = source["configPaths"];
	        this.logPaths = source["logPaths"];
	        this.security = this.convertValues(source["security"], SecurityInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace webshell {
	
	export class TerminalInfo {
	    id: number;
	    currentPath: string;
	    currentUser: string;
	    isWindows: boolean;
	    executePath: string;
	
	    static createFrom(source: any = {}) {
	        return new TerminalInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.currentPath = source["currentPath"];
	        this.currentUser = source["currentUser"];
	        this.isWindows = source["isWindows"];
	        this.executePath = source["executePath"];
	    }
	}
	export class WebClient {
	    ID: number;
	
	    static createFrom(source: any = {}) {
	        return new WebClient(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	    }
	}

}

