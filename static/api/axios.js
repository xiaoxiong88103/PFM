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
    console.log(e.response.data.msg);
    console.log(
      window.ELEMENT.Notification({
        title: "错误",
        message: e.response.data.msg,
      })
    );
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

// apis
// [GIN-debug] POST   /proxy/setPort            --> PFM/func/Proxy.SetPortForward (4 handlers)
// [GIN-debug] POST   /proxy/deletePort         --> PFM/func/Proxy.DeletePortForward (4 handlers)
// [GIN-debug] GET    /proxy/get_port           --> PFM/func/Proxy.ListPortForwards (4 handlers)
// [GIN-debug] GET    /whiteList/reload         --> PFM/route.Proxy_Route.func1 (4 handlers)
// [GIN-debug] POST   /whiteList/add            --> PFM/func/WhiteList.AddWhiteListHandler (4 handlers)
// [GIN-debug] GET    /whiteList/list           --> PFM/func/WhiteList.ViewWhiteListHandler (4 handlers)
// [GIN-debug] POST   /whiteList/delete         --> PFM/func/WhiteList.DeleteWhiteListHandler (4 handlers)
const apiPath = {
  proxySetProt: "/proxy/setPort",
  proxyDeletePort: "/proxy/deletePort",
  proxyProtList: "/proxy/get_port",
  whiteReload: "/whiteList/reload",
  whiteList: "/whiteList/list",
  whiteListAll: "/whiteList/list/all",
  whiteAdd: "/whiteList/add",
  whiteDelete: "/whiteList/delete",
};

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
}
