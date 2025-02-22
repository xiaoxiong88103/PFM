const BASE_URL = "";
var app = new Vue({
  el: "#app",
  data: {
    axios: new Request("http://localhost:8281"),
    activeNavArrow: 1,
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
          id: "1",
          ip: "127.0.0.1",
          comment: "白名单1",
        },
        {
          id: "2",
          ip: "192.168.0.1",
          comment: "白名单2",
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
    proxyDelete(row = {}) {
      let { id } = row;
      this.axios.proxyDeletePort({ id });
    },
    proxyEditOpen(row = {}) {
      this.proxy.editShow = true;
      row.id++;
      row.id = String(row.id);
      row.local_port++;
      row.local_port = String(row.local_port);
      this.proxy.editData = row;
    },
    async proxyEditSave() {
      this.proxy.editShow = false;
      await this.axios.proxySetProt(this.proxy.editData);
      this.init(this.activeNavArrow);
    },
    proxyCloseDialog(key = "") {
      setTimeout(() => {
        this[key].editShow = false;
      }, 100);
    },
    // 初始化
    async init(arrow = 1) {
      if (arrow == 1) {
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
      }
      if (arrow == 3) {
        let data = await this.axios.proxyProtList();
        let list = [];
        for (const i in data.data) {
          if (Object.prototype.hasOwnProperty.call(data.data, i)) {
            const el = data.data[i];
            list.push(el);
          }
        }
        let protList = await Promise.all(
          list.map(async (t) => {
            return {
              ...t,
              data: await this.axios.whiteList({ port: t.local_port }),
            };
          })
        );
        console.log(
          protList.map((t) => {
            return { ...t, data: t?.data?.data || null };
          })
        );
      }
    },
  },
  async mounted() {
    this.init(this.activeNavArrow);
  },
  watch: {},
});
