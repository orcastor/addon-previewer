// * 请求响应参数(不包含data)
export interface Result {
  code: string;
  msg: string;
}

// * 请求响应参数(包含data)
export interface ResultData<T = any> extends Result {
  data?: T;
}

export namespace Object {
  export interface ListOption {
    w?: string;
    d?: string;
    t?: number;
    c: number;
    o?: string;
    b?: number;
  }
  export interface ReqList {
    p: number;
    b: number;
    o: ListOption;
  }
  export interface ObjectInfo {
    i: number; // id
    p: number; // pid
    m: number; // mtime
    d: number; // did
    t: number; // type
    n: string; // name
    s: number; //size
    e: any; // ext
  }
  export interface ResList {
    o?: ObjectInfo[];
  }
  export interface ReqGet {
    b: number;
    i: number;
  }
  export interface ResGet {
    o?: ObjectInfo;
  }
}
