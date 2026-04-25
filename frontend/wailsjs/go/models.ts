export namespace analyzer {
	
	export class DetectionResult {
	    stack: string[];
	    confidence: number;
	    indicators: string[];
	
	    static createFrom(source: any = {}) {
	        return new DetectionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stack = source["stack"];
	        this.confidence = source["confidence"];
	        this.indicators = source["indicators"];
	    }
	}

}

export namespace storage {
	
	export class Rule {
	    ID: number;
	    TemplateID: number;
	    Category: string;
	    Title: string;
	    Content: string;
	    Priority: number;
	    Enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Rule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.TemplateID = source["TemplateID"];
	        this.Category = source["Category"];
	        this.Title = source["Title"];
	        this.Content = source["Content"];
	        this.Priority = source["Priority"];
	        this.Enabled = source["Enabled"];
	    }
	}
	export class Template {
	    ID: number;
	    Name: string;
	    Description: string;
	    Stack: string[];
	    Rules: Rule[];
	    IsBuiltIn: boolean;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Template(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.Stack = source["Stack"];
	        this.Rules = this.convertValues(source["Rules"], Rule);
	        this.IsBuiltIn = source["IsBuiltIn"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
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

export namespace tokenizer {
	
	export class Warning {
	    Model: string;
	    Tokens: number;
	    WarnAt: number;
	    Max: number;
	    Exceeds: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Warning(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Model = source["Model"];
	        this.Tokens = source["Tokens"];
	        this.WarnAt = source["WarnAt"];
	        this.Max = source["Max"];
	        this.Exceeds = source["Exceeds"];
	    }
	}

}

