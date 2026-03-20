export namespace model {
	
	export class OutlineNode {
	    id: string;
	    level: number;
	    title: string;
	    content: string;
	    parentId: string;
	    status: string;
	    wordCount: number;
	    targetWords: number;
	
	    static createFrom(source: any = {}) {
	        return new OutlineNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.level = source["level"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.parentId = source["parentId"];
	        this.status = source["status"];
	        this.wordCount = source["wordCount"];
	        this.targetWords = source["targetWords"];
	    }
	}
	export class Article {
	    id: string;
	    title: string;
	    content: string;
	    outline: OutlineNode[];
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Article(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.outline = this.convertValues(source["outline"], OutlineNode);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.status = source["status"];
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
	export class Config {
	    id: string;
	    apiBaseUrl: string;
	    apiKey: string;
	    model: string;
	    temperature: number;
	    styleTags: string[];
	    audience: string;
	    persona: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.apiBaseUrl = source["apiBaseUrl"];
	        this.apiKey = source["apiKey"];
	        this.model = source["model"];
	        this.temperature = source["temperature"];
	        this.styleTags = source["styleTags"];
	        this.audience = source["audience"];
	        this.persona = source["persona"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class Material {
	    id: string;
	    type: string;
	    title: string;
	    content: string;
	    tags: string[];
	    source: string;
	    // Go type: time
	    createdAt: any;
	    usageCount: number;
	
	    static createFrom(source: any = {}) {
	        return new Material(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.tags = source["tags"];
	        this.source = source["source"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.usageCount = source["usageCount"];
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

