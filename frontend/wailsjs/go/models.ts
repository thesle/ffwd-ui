export namespace models {
	
	export class FileInfo {
	    path: string;
	    size: number;
	    duration: number;
	    format: string;
	    codec: string;
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.size = source["size"];
	        this.duration = source["duration"];
	        this.format = source["format"];
	        this.codec = source["codec"];
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}
	export class MountPoint {
	    path: string;
	    total: number;
	    available: number;
	    used: number;
	
	    static createFrom(source: any = {}) {
	        return new MountPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.total = source["total"];
	        this.available = source["available"];
	        this.used = source["used"];
	    }
	}

}

