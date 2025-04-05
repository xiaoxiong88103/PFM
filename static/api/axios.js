const req = {};

const _conf = (path = "/", methods = "GET", data = {}, headers = {}) => {
  return {
    url: path,
    method: methods,
    headers: headers,
    ...data,
  };
};

req.axios = async (conf) => {
  try {
    let { data, status } = await axios(conf);
    console.log("请求实例：", status, data);
    return data;
  } catch (e) {
    console.log(e);
    window.ELEMENT.Notification({
      title: "错误",
      message: e.response.data.msg,
    });
  }
};

req.post = (path = "/", data = {}, headers = {}) => {
  let options = { data };
  return req.axios(_conf(path, "POST", options, headers));
};

req.get = (path = "/", data = {}) => {
  let options = { params: data };
  return req.axios(_conf(path, "GET", options));
};

// POST   /proxy/setPort
// POST   /proxy/deletePort
// GET    /proxy/getPort
// GET    /whiteList/reload
// POST   /whiteList/add
// GET    /whiteList/list
// GET    /whiteList/list/all
// POST   /whiteList/delete
// GET    /whiteList/number/status
// GET    /whiteList/number/clear
const apiPath = {
  proxySetProt: "/proxy/setPort",
  proxyDeletePort: "/proxy/deletePort",
  proxyProtList: "/proxy/getPort",
  whiteReload: "/whiteList/reload",
  whiteList: "/whiteList/list",
  whiteListAll: "/whiteList/list/all",
  whiteAdd: "/whiteList/add",
  whiteDelete: "/whiteList/delete",
  proxyStatusStop: "/proxy/stopPort",
  proxyStatusStart: "/proxy/restartPort",
  // 流量
  tcpTrafficStats: "/tcp/stats",
};

const StatusTrueEum = 1,
  StatusFalseEum = 2;

// 请求类
class Request {
  baseUrl = "";
  constructor(host = "") {
    this.baseUrl = host;
  }
  proxySetProt(
    data = {
      id: "", //获取ID索引用的
      type: "", // tcp or udp
      remote_ip: "", // 获取远程IP字段
      remote_port: "", //获取远程端口转发的要
      local_port: "", //本地的端口绑定
      comment: "", //备注接口负责备注作用
    }
  ) {
    return req.post(this.baseUrl + apiPath.proxySetProt, data, {});
  }
  proxyDeletePort(
    data = {
      id: "", // 删除端口转发的接口（直接删除，不使用队列）
    }
  ) {
    return req.post(this.baseUrl + apiPath.proxyDeletePort, data, {});
  }
  proxyProtList(data = {}) {
    return req.get(this.baseUrl + apiPath.proxyProtList, data, {});
  }
  whiteReload(data = {}) {
    return req.get(this.baseUrl + apiPath.whiteReload, data, {});
  }
  whiteList(data = { port: "" }) {
    return req.get(this.baseUrl + apiPath.whiteList, data, {});
  }
  whiteListAll() {
    return req.get(this.baseUrl + apiPath.whiteListAll, {}, {});
  }
  whiteAdd(
    data = {
      port: "", // ip地址
      ip: "", // 端口号
    }
  ) {
    return req.post(this.baseUrl + apiPath.whiteAdd, data, {});
  }
  whiteDelete(
    data = {
      port: "", // ip地址
      ip: "", // 端口号
    }
  ) {
    return req.post(this.baseUrl + apiPath.whiteDelete, data, {});
  }
  setProxyStatus(id = 0, thisStatus = 0) {
    console.log(id, thisStatus);
    switch (thisStatus) {
      case StatusFalseEum:
        return req.post(this.baseUrl + apiPath.proxyStatusStart, { id: id }, {});
        break;
      case StatusTrueEum:
        return req.post(this.baseUrl + apiPath.proxyStatusStop, { id: id }, {});
        break;
    }
  }
  tcpTrafficStats() {
    return req.get(this.baseUrl + apiPath.tcpTrafficStats, {}, {});
  }
}
