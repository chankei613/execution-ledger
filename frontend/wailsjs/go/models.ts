export namespace api {
	
	export class EntryFilters {
	    AgentID: string;
	    Source: string;
	    Status: string;
	    Subject: string;
	    Query: string;
	    MinConfidence?: number;
	    MaxConfidence?: number;
	    // Go type: time
	    From?: any;
	    // Go type: time
	    To?: any;
	
	    static createFrom(source: any = {}) {
	        return new EntryFilters(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AgentID = source["AgentID"];
	        this.Source = source["Source"];
	        this.Status = source["Status"];
	        this.Subject = source["Subject"];
	        this.Query = source["Query"];
	        this.MinConfidence = source["MinConfidence"];
	        this.MaxConfidence = source["MaxConfidence"];
	        this.From = this.convertValues(source["From"], null);
	        this.To = this.convertValues(source["To"], null);
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
	export class IssueKeyResult {
	    id: string;
	    name: string;
	    api_key: string;
	
	    static createFrom(source: any = {}) {
	        return new IssueKeyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.api_key = source["api_key"];
	    }
	}
	export class SearchResult {
	    entries: db.LedgerEntry[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entries = this.convertValues(source["entries"], db.LedgerEntry);
	        this.total = source["total"];
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
	export class StatsResult {
	    total: number;
	    by_status: Record<string, number>;
	    avg_confidence: number;
	    low_confidence_rate: number;
	    by_agent: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new StatsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.by_status = source["by_status"];
	        this.avg_confidence = source["avg_confidence"];
	        this.low_confidence_rate = source["low_confidence_rate"];
	        this.by_agent = source["by_agent"];
	    }
	}

}

export namespace db {
	
	export class Action {
	    tool: string;
	    input_summary: string;
	    // Go type: time
	    timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new Action(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tool = source["tool"];
	        this.input_summary = source["input_summary"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
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
	export class AgentKey {
	    id: string;
	    name: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    revoked_at?: any;
	
	    static createFrom(source: any = {}) {
	        return new AgentKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.revoked_at = this.convertValues(source["revoked_at"], null);
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
	export class ConfidenceBreakdown {
	    task_understood: number;
	    execution_complete: number;
	    correctness: number;
	    side_effects_clean: number;
	
	    static createFrom(source: any = {}) {
	        return new ConfidenceBreakdown(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.task_understood = source["task_understood"];
	        this.execution_complete = source["execution_complete"];
	        this.correctness = source["correctness"];
	        this.side_effects_clean = source["side_effects_clean"];
	    }
	}
	export class CriterionResult {
	    description: string;
	    met: boolean;
	    evidence?: string;
	
	    static createFrom(source: any = {}) {
	        return new CriterionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.description = source["description"];
	        this.met = source["met"];
	        this.evidence = source["evidence"];
	    }
	}
	export class Decision {
	    description: string;
	    rationale: string;
	    alternatives_considered: string[];
	
	    static createFrom(source: any = {}) {
	        return new Decision(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.description = source["description"];
	        this.rationale = source["rationale"];
	        this.alternatives_considered = source["alternatives_considered"];
	    }
	}
	export class FollowUp {
	    description: string;
	    suggested_task?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new FollowUp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.description = source["description"];
	        this.suggested_task = source["suggested_task"];
	    }
	}
	export class Usage {
	    input_tokens: number;
	    output_tokens: number;
	    mcp_calls_by_server?: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new Usage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input_tokens = source["input_tokens"];
	        this.output_tokens = source["output_tokens"];
	        this.mcp_calls_by_server = source["mcp_calls_by_server"];
	    }
	}
	export class LedgerEntry {
	    id: string;
	    // Go type: time
	    received_at: any;
	    source: string;
	    agent_id: string;
	    subject: string;
	    status: string;
	    summary: string;
	    criteria_results: CriterionResult[];
	    outputs: Record<string, any>;
	    confidence_overall: number;
	    confidence_breakdown: ConfidenceBreakdown;
	    low_confidence_areas: string[];
	    decisions: Decision[];
	    actions_taken: Action[];
	    follow_up: FollowUp[];
	    usage: Usage;
	
	    static createFrom(source: any = {}) {
	        return new LedgerEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.received_at = this.convertValues(source["received_at"], null);
	        this.source = source["source"];
	        this.agent_id = source["agent_id"];
	        this.subject = source["subject"];
	        this.status = source["status"];
	        this.summary = source["summary"];
	        this.criteria_results = this.convertValues(source["criteria_results"], CriterionResult);
	        this.outputs = source["outputs"];
	        this.confidence_overall = source["confidence_overall"];
	        this.confidence_breakdown = this.convertValues(source["confidence_breakdown"], ConfidenceBreakdown);
	        this.low_confidence_areas = source["low_confidence_areas"];
	        this.decisions = this.convertValues(source["decisions"], Decision);
	        this.actions_taken = this.convertValues(source["actions_taken"], Action);
	        this.follow_up = this.convertValues(source["follow_up"], FollowUp);
	        this.usage = this.convertValues(source["usage"], Usage);
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

