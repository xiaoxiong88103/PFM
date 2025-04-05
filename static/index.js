let BASE_URL = "";
var app = new Vue({
  el: "#app",
  data: {
    axios: new Request("http://" + window.location.hostname + ":8281"),
    // axios: new Request("http://aws.xiaoxiong.host:8281"),
    activeNavArrow: "0",
    // 代理配置
    proxy: {
      list: [],
      editShow: false,
      editData: {
        id: "",
        type: "",
        remote_ip: "",
        remote_port: "",
        local_port: "",
        comment: "",
      },
      editFormInputWidth: "120px",
    },
    // 白名单配置
    whiteList: {
      list: [],
      editShow: false,
      editData: {
        port: "",
        ip: "",
      },
      editFormInputWidth: "120px",
    },
    // 转发流量统计
    tcpProxyTrafficStats: {},
    
  },
  methods: {
    // 切换导航
    async navSelect(key, keyPath) {
      this.activeNavArrow = key;
      setTimeout(() => {
        this.init(this.activeNavArrow);
      }, 100);
    },
    // 删除代理
    proxyDelete(row = {}) {
      let { id } = row;
      this.axios.proxyDeletePort({ id });
      setTimeout(() => {
        this.init(this.activeNavArrow);
      }, 300);
    },
    async proxyLoad(isRestful = 2) {
      if (isRestful === 1) {
        window.ELEMENT.Notification({
          title: "提示",
          message: "刷新中...",
        });
      }
      let data = await this.axios.proxyProtList();
      this.proxy.list = [];
      let list = [];
      for (const i in data.data) {
        if (Object.prototype.hasOwnProperty.call(data.data, i)) {
          const el = data.data[i];
          list.push(el);
        }
      }
      console.log(list.sort((a, b) => Number(b.id) - Number(a.id)));
      this.proxy.list = list;
      this.proxy.editData = {
        id: "",
        type: "",
        remote_ip: "",
        remote_port: "",
        local_port: "",
        comment: "",
      };
      // 刷新流量实时数据
      this.loadTcpTrafficStats();
    },
    // 编辑代理
    proxyEditOpen(row = {}, isClone = false) {
      row = JSON.parse(JSON.stringify(row));
      if (row.id !== undefined && isClone) {
        let len = this.proxy.list.length;
        len++;
        row.id = String(len);
        row.local_port++;
        row.local_port = String(row.local_port);
      }
      this.proxy.editData = row;
      this.proxy.editShow = true;
    },
    editProxyStatus(row = { status: 1 }) {
      let { id, status } = row;
      this.axios.setProxyStatus(id, status);
      setTimeout(() => {
        this.init(this.activeNavArrow);
      }, 300);
      window.ELEMENT.Notification({
        title: "提示",
        type: "success",
        message: "操作完成",
      });
    },
    // 编辑代理保存
    async proxyEditSave() {
      this.proxy.editShow = false;
      await this.axios.proxySetProt(this.proxy.editData);
      setTimeout(() => {
        this.init(this.activeNavArrow);
      }, 300);
    },
    // 退出编辑代理
    closeDialog(key = "") {
      setTimeout(() => {
        this[key].editShow = false;
      }, 100);
    },
    // 白名单点击
    async whiteRemove(data = {}) {
      let res = await this.axios.whiteDelete({
        ip: data.label,
        port: data.port,
      });
      window.ELEMENT.Notification({
        title: "提示",
        message: res.msg,
      });
      setTimeout(() => {
        this.init(this.activeNavArrow);
      }, 300);
    },
    // 重载白名单
    async whiteReload() {
      let data = await this.axios.whiteReload();
      window.ELEMENT.Notification({
        title: "提示",
        message: data.msg,
      });
    },
    // 编辑白名单
    whiteEditOpen(row = {}) {
      this.whiteList.editShow = true;
      if (row.label !== undefined) {
        this.whiteList.editData = {
          ip: row.label,
          port: row.port,
        };
      }
    },
    // 加载白名单数据
    async whiteLoad(isRestful = 2) {
      if (isRestful === 1) {
        window.ELEMENT.Notification({
          title: "提示",
          message: "刷新中...",
        });
      }
      let data = await this.axios.whiteListAll();
      this.whiteList.list = [];
      for (const key in data.data) {
        if (Object.prototype.hasOwnProperty.call(data.data, key)) {
          const it = data.data[key];
          let item = {};
          item.label = key;
          item.child = false;
          item.children = it.map((t) => {
            return { label: t, child: true, port: key };
          });
          this.whiteList.list.push(item);
        }
      }
      this.whiteList.editData = {
        port: "",
        ip: "",
      };
    },
    // 白名单编辑保存
    async whiteEditSave() {
      this.whiteList.editShow = false;
      await this.axios.whiteAdd(this.whiteList.editData);
      setTimeout(() => {
        this.init(this.activeNavArrow);
      }, 300);
    },
    // 加载tcp流量状态
    async loadTcpTrafficStats() {
      let data = await this.axios.tcpTrafficStats();
      this.tcpProxyTrafficStats = data;
    },
    // 初始化
    async init(arrow = 1) {
      if (arrow == 1) {
        this.proxyLoad();
      }
      if (arrow == 2) {
        this.whiteLoad();
      }
      // 初始化代理流量抓取代理数据
      this.loadTcpTrafficStats();
      setInterval(() => {
        this.loadTcpTrafficStats();
      }, 1000 * 30);
    },
  },
  async mounted() {
    this.init(this.activeNavArrow);
  },
  watch: {},
});
