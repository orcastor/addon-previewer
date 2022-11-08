class GlobalStore {
  // token
  token: string;

  constructor() {
    this.token = localStorage.getItem('token')||'';
  }

  // setToken
  setToken(token: string) {
    this.token = token;
    localStorage.setItem('token', token);
  }
};

export const store:GlobalStore = new GlobalStore;
