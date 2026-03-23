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
	    sourceType: string;
	    workflowRunId: string;
	    topicId: string;
	    qualityScore: number;
	    readTime: number;
	    wordCount: number;
	    publishTaskId: string;
	    // Go type: time
	    publishedAt?: any;
	
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
	        this.sourceType = source["sourceType"];
	        this.workflowRunId = source["workflowRunId"];
	        this.topicId = source["topicId"];
	        this.qualityScore = source["qualityScore"];
	        this.readTime = source["readTime"];
	        this.wordCount = source["wordCount"];
	        this.publishTaskId = source["publishTaskId"];
	        this.publishedAt = this.convertValues(source["publishedAt"], null);
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
	export class HotTrend {
	    id: string;
	    platform: string;
	    title: string;
	    url: string;
	    hotRank: number;
	    hotValue: number;
	    category: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new HotTrend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.platform = source["platform"];
	        this.title = source["title"];
	        this.url = source["url"];
	        this.hotRank = source["hotRank"];
	        this.hotValue = source["hotValue"];
	        this.category = source["category"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
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
	
	export class Topic {
	    id: string;
	    title: string;
	    category: string;
	    source: string;
	    sourceUrl: string;
	    score: number;
	    hotScore: number;
	    compScore: number;
	    fitScore: number;
	    keywords: string[];
	    summary: string;
	    references: string[];
	    angles: string[];
	    status: string;
	    workflowRunId: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Topic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.category = source["category"];
	        this.source = source["source"];
	        this.sourceUrl = source["sourceUrl"];
	        this.score = source["score"];
	        this.hotScore = source["hotScore"];
	        this.compScore = source["compScore"];
	        this.fitScore = source["fitScore"];
	        this.keywords = source["keywords"];
	        this.summary = source["summary"];
	        this.references = source["references"];
	        this.angles = source["angles"];
	        this.status = source["status"];
	        this.workflowRunId = source["workflowRunId"];
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
	export class WorkflowStep {
	    id: string;
	    name: string;
	    type: string;
	    config: Record<string, any>;
	    nextStep: string;
	    onError: string;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowStep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.config = source["config"];
	        this.nextStep = source["nextStep"];
	        this.onError = source["onError"];
	    }
	}
	export class WorkflowTrigger {
	    type: string;
	    cronExpr: string;
	    rssUrl: string;
	    webhookUrl: string;
	    config: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowTrigger(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.cronExpr = source["cronExpr"];
	        this.rssUrl = source["rssUrl"];
	        this.webhookUrl = source["webhookUrl"];
	        this.config = source["config"];
	    }
	}
	export class Workflow {
	    id: string;
	    name: string;
	    description: string;
	    trigger: WorkflowTrigger;
	    steps: WorkflowStep[];
	    autoPublish: boolean;
	    needReview: boolean;
	    status: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Workflow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.trigger = this.convertValues(source["trigger"], WorkflowTrigger);
	        this.steps = this.convertValues(source["steps"], WorkflowStep);
	        this.autoPublish = source["autoPublish"];
	        this.needReview = source["needReview"];
	        this.status = source["status"];
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
	export class WorkflowRunStep {
	    stepId: string;
	    status: string;
	    input: any;
	    output: any;
	    // Go type: time
	    startedAt: any;
	    // Go type: time
	    completedAt?: any;
	    error: string;
	    retryCount: number;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowRunStep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stepId = source["stepId"];
	        this.status = source["status"];
	        this.input = source["input"];
	        this.output = source["output"];
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.completedAt = this.convertValues(source["completedAt"], null);
	        this.error = source["error"];
	        this.retryCount = source["retryCount"];
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
	export class WorkflowRun {
	    id: string;
	    workflowId: string;
	    status: string;
	    currentStep: string;
	    input: Record<string, any>;
	    output: Record<string, any>;
	    steps: WorkflowRunStep[];
	    // Go type: time
	    startedAt: any;
	    // Go type: time
	    completedAt?: any;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowRun(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.workflowId = source["workflowId"];
	        this.status = source["status"];
	        this.currentStep = source["currentStep"];
	        this.input = source["input"];
	        this.output = source["output"];
	        this.steps = this.convertValues(source["steps"], WorkflowRunStep);
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.completedAt = this.convertValues(source["completedAt"], null);
	        this.error = source["error"];
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

