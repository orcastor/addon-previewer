import axios, {
  AxiosInstance,
  AxiosError,
  AxiosRequestConfig,
  AxiosResponse,
} from "axios";
import { ResultData } from "@/api/interface";
import { store } from "@/store";

var callback:({}) => void = ((_o) => {});

export function setCallback(cb:({}) => void) {
    callback = cb;
}

// * 请求枚举配置
/**
 * @description：请求配置
 */
 export enum ResultEnum {
  SUCCESS = 200,
  ERROR = 500,
  OVERDUE = 599,
  TIMEOUT = 10000,
  TYPE = "success",
}

/**
 * @description：请求方法
 */
export enum RequestEnum {
  GET = "GET",
  POST = "POST",
  PATCH = "PATCH",
  PUT = "PUT",
  DELETE = "DELETE",
}

/**
 * @description：常用的contentTyp类型
 */
export enum ContentTypeEnum {
  // json
  JSON = "application/json;charset=UTF-8",
  // text
  TEXT = "text/plain;charset=UTF-8",
  // form-data 一般配合qs
  FORM_URLENCODED = "application/x-www-form-urlencoded;charset=UTF-8",
  // form-data 上传
  FORM_DATA = "multipart/form-data;charset=UTF-8",
}

/**
 * @description: 校验网络请求状态码
 * @param {Number} status
 * @return void
 */
const checkStatus = (status: number): void => {
  switch (status) {
    case 400:
      callback({ type: 'warning', message: "请求失败！请您稍后重试" });
      break;
    case 401:
      callback({ type: 'warning', message: "登录失效！请您重新登录" });
      break;
    case 403:
      callback({ type: 'warning', message: "当前账号无权限访问！" });
      break;
    case 404:
      callback({ type: 'warning', message: "你所访问的资源不存在！" });
      break;
    case 405:
      callback({ type: 'warning', message: "请求方式错误！请您稍后重试" });
      break;
    case 408:
      callback({ type: 'warning', message: "请求超时！请您稍后重试" });
      break;
    case 500:
      callback({ type: 'warning', message: "服务异常！" });
      break;
    case 502:
      callback({ type: 'warning', message: "网关错误！" });
      break;
    case 503:
      callback({ type: 'warning', message: "服务不可用！" });
      break;
    case 504:
      callback({ type: 'warning', message: "网关超时！" });
      break;
  }
};

/**
 * pinia 错误使用说明示例
 * https://github.com/vuejs/pinia/discussions/971
 * https://github.com/vuejs/pinia/discussions/664#discussioncomment-1329898
 * https://pinia.vuejs.org/core-concepts/outside-component-usage.html#single-page-applications
 */
const config = {
  // 默认地址请求地址，可在 .env 开头文件中修改
  baseURL: import.meta.env.VITE_API_URL as string,
  // 设置超时时间（10s）
  timeout: ResultEnum.TIMEOUT as number,
  // 跨域时候允许携带凭证
  withCredentials: true,
};

class RequestHttp {
  service: AxiosInstance;
  public constructor(config: AxiosRequestConfig) {
    // 实例化axios
    this.service = axios.create(config);

    /**
     * @description 请求拦截器
     * 客户端发送请求 -> [请求拦截器] -> 服务器
     * token校验(JWT) : 接受服务器返回的token,存储到vuex/pinia/本地储存当中
     */
    this.service.interceptors.request.use(
      (config: AxiosRequestConfig) => {
        const token: string = store.token;
        return {
          ...config,
          headers: { ...config.headers, "Authorization": token },
        };
      },
      (error: AxiosError) => {
        return Promise.reject(error);
      }
    );

    /**
     * @description 响应拦截器
     *  服务器换返回信息 -> [拦截统一处理] -> 客户端JS获取到信息
     */
    this.service.interceptors.response.use(
      (response: AxiosResponse) => {
        const { data } = response;
        // * 登陆失效（code == 599）
        if (data.code == ResultEnum.OVERDUE) {
          callback({ type: 'warning', message: data.msg });
          store.setToken("");
          return Promise.reject(data);
        }
        // * 全局错误信息拦截（防止下载文件得时候返回数据流，没有code，直接报错）
        if (data.code && data.code !== ResultEnum.SUCCESS) {
          callback({ type: 'warning', message: data.msg });
          return Promise.reject(data);
        }
        // * 成功请求（在页面上除非特殊情况，否则不用处理失败逻辑）
        return data;
      },
      async (error: AxiosError) => {
        const { response } = error;
        // 请求超时单独判断，因为请求超时没有 response
        if (error.message.indexOf("timeout") !== -1)
          callback({ type: 'warning', message: "请求超时！请您稍后重试" });
        // 根据响应的错误状态码，做不同的处理
        if (response) checkStatus(response.status);
        // 服务器结果都没有返回(可能服务器错误可能客户端断网)，断网处理:可以跳转到断网页面
        // if (!window.navigator.onLine) router.replace({ path: "/500" });
        return Promise.reject(error);
      }
    );
  }

  // * 常用请求方法封装
  get<T>(url: string, params?: object, _object = {}): Promise<ResultData<T>> {
    return this.service.get(url, { params, ..._object });
  }
  post<T>(url: string, params?: object, _object = {}): Promise<ResultData<T>> {
    return this.service.post(url, params, _object);
  }
  put<T>(url: string, params?: object, _object = {}): Promise<ResultData<T>> {
    return this.service.put(url, params, _object);
  }
  delete<T>(url: string, params?: any, _object = {}): Promise<ResultData<T>> {
    return this.service.delete(url, { params, ..._object });
  }
}

export default new RequestHttp(config);
