import { Object } from "@/api/interface/index";
import { API } from "@/api/config/config";

import http from "@/api";

export const listApi = (params: Object.ReqList) => {
  return http.post<Object.ResList>(API + `/list`, params);
};

export const getApi = (params: Object.ReqGet) => {
  return http.post<Object.ResGet>(API + `/get`, params);
};
