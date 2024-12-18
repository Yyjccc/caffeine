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
	export class WebShellItem {
	    ID: number;
	    Location: string;
	    ShellType: string;
	    IP: string;
	    CreateTime: string;
	    UpdateTime: string;
	    URL: string;
	    Note: string;
	
	    static createFrom(source: any = {}) {
	        return new WebShellItem(source);
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
	    }
	}

}

export namespace core {
	
	export class OSInfo {
	    name: string;
	    version: string;
	    arch: string;
	
	    static createFrom(source: any = {}) {
	        return new OSInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.arch = source["arch"];
	    }
	}
	export class SystemInfo {
	    fileRoot: string;
	    currentDir: string;
	    currentUser: string;
	    processArch: number;
	    tempDirectory: string;
	    ipList: string[];
	    os: OSInfo;
	    env: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileRoot = source["fileRoot"];
	        this.currentDir = source["currentDir"];
	        this.currentUser = source["currentUser"];
	        this.processArch = source["processArch"];
	        this.tempDirectory = source["tempDirectory"];
	        this.ipList = source["ipList"];
	        this.os = this.convertValues(source["os"], OSInfo);
	        this.env = source["env"];
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

