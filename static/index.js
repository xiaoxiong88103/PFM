const BASE_URL = "";
var app = new Vue({
  el: "#app",
  data: {
    axios: new Request("http://127.0.0.1:8281"),
    activeNavArrow: 0,
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
      list: [
        {
          label: "一级 1",
          children: [
            {
              label: "二级 1-1",
            },
          ],
        },
      ],
      editShow: false,
      editData: {
        id: "",
        ip: "",
        comment: "",
      },
      editFormInputWidth: "120px",
    },
  },
  methods: {
    // 切换导航
    async navSelect(key, keyPath) {
      this.activeNavArrow = key;
      setTimeout(() => {
        this.init(this.activeNavArrow);
      }, 100);
    },
    // 代理操作
    handleClick(row) {
      console.log(row);
    },
    // 删除代理
    proxyDelete(row = {}) {
      let { id } = row;
      this.axios.proxyDeletePort({ id });
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
      for (const i in data.data) {
        if (Object.prototype.hasOwnProperty.call(data.data, i)) {
          const el = data.data[i];
          this.proxy.list.push(el);
        }
      }
      this.proxy.editData = {
        id: "",
        type: "",
        remote_ip: "",
        remote_port: "",
        local_port: "",
        comment: "",
      };
    },
    // 编辑代理
    proxyEditOpen(row = {}) {
      this.proxy.editShow = true;
      if (row.id !== undefined) {
        row.id++;
        row.id = String(row.id);
        row.local_port++;
        row.local_port = String(row.local_port);
        this.proxy.editData = row;
      }
    },

    // 编辑代理保存
    async proxyEditSave() {
      this.proxy.editShow = false;
      await this.axios.proxySetProt(this.proxy.editData);
      this.init(this.activeNavArrow);
    },
    // 退出编辑代理
    closeDialog(key = "") {
      setTimeout(() => {
        this[key].editShow = false;
      }, 100);
    },
    // 白名单点击
    whiteClick(data = {}) {
      console.log("点击了白名单item", data);
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
      if (row.id !== undefined) {
        row.id++;
        row.id = String(row.id);
        this.whiteList.editData = row;
      }
    },
    async whiteLoad(isRestful = 2) {
      console.log(isRestful);
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
            return { label: t, child: true };
          });
          this.whiteList.list.push(item);
        }
      }
      this.whiteList.editData = {
        id: "",
        ip: "",
        comment: "",
      };
    },
    async whiteEditSave() {
      this.whiteList.editShow = false;
      await this.axios.whiteAdd(this.whiteList.editData);
      this.init(this.activeNavArrow);
      console.log(this.whiteList.editData);
    },
    // 初始化
    async init(arrow = 1) {
      console.log(arrow);
      if (arrow == 1) {
        this.proxyLoad();
      }
      if (arrow == 2) {
        this.whiteLoad();
      }
    },
  },
  async mounted() {
    this.init(this.activeNavArrow);
  },
  watch: {},
});
