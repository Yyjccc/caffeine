export namespace client {
	
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

