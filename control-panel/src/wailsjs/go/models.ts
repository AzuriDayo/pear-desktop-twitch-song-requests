export namespace main {
	
	export class HistoryItem {
	    id: number;
	    video_id: string;
	    twitch_username: string;
	    requested_at: string;
	    is_ninja: boolean;
	    song_title: string;
	    artist_name: string;
	    image_url: string;
	
	    static createFrom(source: any = {}) {
	        return new HistoryItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.video_id = source["video_id"];
	        this.twitch_username = source["twitch_username"];
	        this.requested_at = source["requested_at"];
	        this.is_ninja = source["is_ninja"];
	        this.song_title = source["song_title"];
	        this.artist_name = source["artist_name"];
	        this.image_url = source["image_url"];
	    }
	}
	export class HistoryResponse {
	    max_results: number;
	    items: HistoryItem[];
	
	    static createFrom(source: any = {}) {
	        return new HistoryResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.max_results = source["max_results"];
	        this.items = this.convertValues(source["items"], HistoryItem);
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
	export class SelectableReward {
	    id: string;
	    name: string;
	    cost: number;
	
	    static createFrom(source: any = {}) {
	        return new SelectableReward(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.cost = source["cost"];
	    }
	}
	export class SettingsResponse {
	    reward_id: string;
	    cmd_permissions: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new SettingsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.reward_id = source["reward_id"];
	        this.cmd_permissions = source["cmd_permissions"];
	    }
	}
	export class SongQueueItem {
	    requested_by: string;
	    song: songrequests.SongResult;
	    is_ninja: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SongQueueItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.requested_by = source["requested_by"];
	        this.song = this.convertValues(source["song"], songrequests.SongResult);
	        this.is_ninja = source["is_ninja"];
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
	export class TwitchInfoPayload {
	    type: string;
	    stream_online: boolean;
	    reward_id: string;
	    login: string;
	    login_bot: string;
	    expiry_date: string;
	    expiry_date_bot: string;
	    refresh_expiry_date: string;
	    refresh_expiry_date_bot: string;
	
	    static createFrom(source: any = {}) {
	        return new TwitchInfoPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.stream_online = source["stream_online"];
	        this.reward_id = source["reward_id"];
	        this.login = source["login"];
	        this.login_bot = source["login_bot"];
	        this.expiry_date = source["expiry_date"];
	        this.expiry_date_bot = source["expiry_date_bot"];
	        this.refresh_expiry_date = source["refresh_expiry_date"];
	        this.refresh_expiry_date_bot = source["refresh_expiry_date_bot"];
	    }
	}

}

export namespace songrequests {
	
	export class SongResult {
	    title: string;
	    artist: string;
	    videoId: string;
	    imageUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new SongResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.artist = source["artist"];
	        this.videoId = source["videoId"];
	        this.imageUrl = source["imageUrl"];
	    }
	}

}

