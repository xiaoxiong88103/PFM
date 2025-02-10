const BASE_URL = "";
var app = new Vue({
  el: "#app",
  data: {
    activeNavArrow: 1,
    // 代理配置
    proxy: {
      list: [
        {
          id: "1",
          type: "tcp",
          remote_ip: "127.0.0.1",
          remote_port: "3309",
          local_port: "2201",
          comment: "mysql",
        },
        {
          id: "2",
          type: "tcp",
          remote_ip: "127.0.0.1",
          remote_port: "3309",
          local_port: "22011",
          comment: "mysql",
        },
      ],
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
    navSelect(key, keyPath) {
      // console.log(key, keyPath);
      this.activeNavArrow = key;
      JSON.stringify
    },
    // 代理操作
    handleClick(row) {
      console.log(row);
    },
    proxyDelete(row = {}) {
      this.$notify({
        title: "成功",
        message: "移除成功",
        type: "success",
      });
    },
    proxyEdit(row = {}) {
      this.proxy.editShow = true;
      console.log(row);
    },
    proxyCloseDialog() {
      setTimeout(() => {
        this.proxy.editShow = false;
      }, 100);
    },
  },
  mounted() {
    console.log(this.axios);
  },
  watch: {},
});
