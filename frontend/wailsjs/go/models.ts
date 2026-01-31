export namespace models {
	
	export class Account {
	    id: number;
	    account: string;
	    password: string;
	    accountType: string;
	    isSold: boolean;
	    // Go type: time
	    soldAt?: any;
	    // Go type: time
	    expireAt?: any;
	    reminderSent: boolean;
	    notes: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.account = source["account"];
	        this.password = source["password"];
	        this.accountType = source["accountType"];
	        this.isSold = source["isSold"];
	        this.soldAt = this.convertValues(source["soldAt"], null);
	        this.expireAt = this.convertValues(source["expireAt"], null);
	        this.reminderSent = source["reminderSent"];
	        this.notes = source["notes"];
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
	export class AccountStats {
	    total: number;
	    plusCount: number;
	    businessCount: number;
	    freeCount: number;
	    soldCount: number;
	    expiredCount: number;
	    expiringIn7Days: number;
	
	    static createFrom(source: any = {}) {
	        return new AccountStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.plusCount = source["plusCount"];
	        this.businessCount = source["businessCount"];
	        this.freeCount = source["freeCount"];
	        this.soldCount = source["soldCount"];
	        this.expiredCount = source["expiredCount"];
	        this.expiringIn7Days = source["expiringIn7Days"];
	    }
	}
	export class AuditLog {
	    id: number;
	    // Go type: time
	    timestamp: any;
	    user: string;
	    action: string;
	    resourceType: string;
	    resourceId: number;
	    ipAddress: string;
	    details: string;
	    success: boolean;
	    errorMessage: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new AuditLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.user = source["user"];
	        this.action = source["action"];
	        this.resourceType = source["resourceType"];
	        this.resourceId = source["resourceId"];
	        this.ipAddress = source["ipAddress"];
	        this.details = source["details"];
	        this.success = source["success"];
	        this.errorMessage = source["errorMessage"];
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
	export class BatchImportResult {
	    success: number;
	    errors: string[];
	
	    static createFrom(source: any = {}) {
	        return new BatchImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.errors = source["errors"];
	    }
	}
	export class EmailConfig {
	    id: number;
	    smtpHost: string;
	    smtpPort: number;
	    senderEmail: string;
	    senderPassword: string;
	    recipientEmail: string;
	    isActive: boolean;
	    useRemoteServer: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new EmailConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.smtpHost = source["smtpHost"];
	        this.smtpPort = source["smtpPort"];
	        this.senderEmail = source["senderEmail"];
	        this.senderPassword = source["senderPassword"];
	        this.recipientEmail = source["recipientEmail"];
	        this.isActive = source["isActive"];
	        this.useRemoteServer = source["useRemoteServer"];
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
	export class EmailLog {
	    id: number;
	    subject: string;
	    content: string;
	    recipient: string;
	    status: string;
	    error: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new EmailLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.subject = source["subject"];
	        this.content = source["content"];
	        this.recipient = source["recipient"];
	        this.status = source["status"];
	        this.error = source["error"];
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
	export class EmailLogsResult {
	    logs: EmailLog[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new EmailLogsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.logs = this.convertValues(source["logs"], EmailLog);
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
	export class HostKey {
	    id: number;
	    host: string;
	    port: number;
	    keyType: string;
	    fingerprint: string;
	    publicKey: string;
	    // Go type: time
	    firstSeen: any;
	    // Go type: time
	    lastUsed: any;
	    trusted: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new HostKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.keyType = source["keyType"];
	        this.fingerprint = source["fingerprint"];
	        this.publicKey = source["publicKey"];
	        this.firstSeen = this.convertValues(source["firstSeen"], null);
	        this.lastUsed = this.convertValues(source["lastUsed"], null);
	        this.trusted = source["trusted"];
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
	export class PaginatedAccounts {
	    data: Account[];
	    total: number;
	    page: number;
	    pageSize: number;
	    totalPages: number;
	
	    static createFrom(source: any = {}) {
	        return new PaginatedAccounts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], Account);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.pageSize = source["pageSize"];
	        this.totalPages = source["totalPages"];
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
	export class PaginatedAuditLogs {
	    data: AuditLog[];
	    total: number;
	    page: number;
	    pageSize: number;
	    totalPages: number;
	
	    static createFrom(source: any = {}) {
	        return new PaginatedAuditLogs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], AuditLog);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.pageSize = source["pageSize"];
	        this.totalPages = source["totalPages"];
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
	export class SMTPProvider {
	    name: string;
	    host: string;
	    port: number;
	    helpText: string;
	
	    static createFrom(source: any = {}) {
	        return new SMTPProvider(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.helpText = source["helpText"];
	    }
	}
	export class ServerConfig {
	    id: number;
	    host: string;
	    port: number;
	    username: string;
	    password: string;
	    privateKey: string;
	    deployPath: string;
	    isActive: boolean;
	    // Go type: time
	    lastDeployedAt?: any;
	    serviceStatus: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ServerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.privateKey = source["privateKey"];
	        this.deployPath = source["deployPath"];
	        this.isActive = source["isActive"];
	        this.lastDeployedAt = this.convertValues(source["lastDeployedAt"], null);
	        this.serviceStatus = source["serviceStatus"];
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
	export class ServerInfo {
	    osName: string;
	    osVersion: string;
	    osPrettyName: string;
	    packageManager: string;
	    hasSystemd: boolean;
	    systemdVersion: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.osName = source["osName"];
	        this.osVersion = source["osVersion"];
	        this.osPrettyName = source["osPrettyName"];
	        this.packageManager = source["packageManager"];
	        this.hasSystemd = source["hasSystemd"];
	        this.systemdVersion = source["systemdVersion"];
	    }
	}
	export class SystemConfig {
	    id: number;
	    defaultValidityDays: number;
	    reminderDaysBefore: number;
	    copyFormat: string;
	    emailFormat: string;
	    accountTypes: string;
	    accountStatuses: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new SystemConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.defaultValidityDays = source["defaultValidityDays"];
	        this.reminderDaysBefore = source["reminderDaysBefore"];
	        this.copyFormat = source["copyFormat"];
	        this.emailFormat = source["emailFormat"];
	        this.accountTypes = source["accountTypes"];
	        this.accountStatuses = source["accountStatuses"];
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

}

