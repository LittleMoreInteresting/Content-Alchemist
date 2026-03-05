export namespace models {
	
	export class AIConfig {
	    baseUrl: string;
	    token: string;
	    temperature: number;
	    model: string;
	
	    static createFrom(source: any = {}) {
	        return new AIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.baseUrl = source["baseUrl"];
	        this.token = source["token"];
	        this.temperature = source["temperature"];
	        this.model = source["model"];
	    }
	}
	export class Article {
	    id: number;
	    uuid: string;
	    filePath: string;
	    title: string;
	    summary: string;
	    tags: string[];
	    wordCount: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    // Go type: time
	    lastOpenedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Article(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.uuid = source["uuid"];
	        this.filePath = source["filePath"];
	        this.title = source["title"];
	        this.summary = source["summary"];
	        this.tags = source["tags"];
	        this.wordCount = source["wordCount"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.lastOpenedAt = this.convertValues(source["lastOpenedAt"], null);
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
	export class GenerateOutlineResult {
	    titles: string[];
	    outline: string;
	
	    static createFrom(source: any = {}) {
	        return new GenerateOutlineResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.titles = source["titles"];
	        this.outline = source["outline"];
	    }
	}
	export class ReadArticleResponse {
	    article?: Article;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new ReadArticleResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.article = this.convertValues(source["article"], Article);
	        this.content = source["content"];
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

