export namespace main {
	
	export class Gateway {
	    gateway_address: string;
	    gateway_metric: number;
	
	    static createFrom(source: any = {}) {
	        return new Gateway(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gateway_address = source["gateway_address"];
	        this.gateway_metric = source["gateway_metric"];
	    }
	}
	export class Ip {
	    ip_address: string;
	    subnet_mask: string;
	    cidr: number;
	
	    static createFrom(source: any = {}) {
	        return new Ip(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip_address = source["ip_address"];
	        this.subnet_mask = source["subnet_mask"];
	        this.cidr = source["cidr"];
	    }
	}
	export class Interface {
	    interface_name: string;
	    interface_metric: number;
	    interface_type: string;
	    connected: boolean;
	    disabled: boolean;
	    ip_is_dhcp: boolean;
	    ips: Ip[];
	    gateways: Gateway[];
	    dns_is_dhcp: boolean;
	    dns_servers: string[];
	
	    static createFrom(source: any = {}) {
	        return new Interface(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.interface_name = source["interface_name"];
	        this.interface_metric = source["interface_metric"];
	        this.interface_type = source["interface_type"];
	        this.connected = source["connected"];
	        this.disabled = source["disabled"];
	        this.ip_is_dhcp = source["ip_is_dhcp"];
	        this.ips = this.convertValues(source["ips"], Ip);
	        this.gateways = this.convertValues(source["gateways"], Gateway);
	        this.dns_is_dhcp = source["dns_is_dhcp"];
	        this.dns_servers = source["dns_servers"];
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
	export class InterfaceConfig {
	    interface_name: string;
	    interface_metric: number;
	    ip_is_dhcp: boolean;
	    ips: Ip[];
	    gateways: Gateway[];
	    dns_is_dhcp: boolean;
	    dns_servers: string[];
	
	    static createFrom(source: any = {}) {
	        return new InterfaceConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.interface_name = source["interface_name"];
	        this.interface_metric = source["interface_metric"];
	        this.ip_is_dhcp = source["ip_is_dhcp"];
	        this.ips = this.convertValues(source["ips"], Ip);
	        this.gateways = this.convertValues(source["gateways"], Gateway);
	        this.dns_is_dhcp = source["dns_is_dhcp"];
	        this.dns_servers = source["dns_servers"];
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
	export class InterfaceState {
	    interface_name: string;
	    interface_type: string;
	    connected: boolean;
	    disabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new InterfaceState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.interface_name = source["interface_name"];
	        this.interface_type = source["interface_type"];
	        this.connected = source["connected"];
	        this.disabled = source["disabled"];
	    }
	}

}

