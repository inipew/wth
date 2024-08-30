const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  transpileDependencies: true,
  // devServer: {
  //   port: 5678, // Mengatur port ke 5678
  //   host: "157.230.247.64", // Atur host untuk menerima koneksi dari semua alamat
  //   allowedHosts: "all",
  //   proxy: {
  //     "/api": {
  //       target: "http://157.230.247.64:4567", // Mengatur proxy ke backend
  //       changeOrigin: true,
  //       pathRewrite: { "^/api": "" },
  //     },
  //   },
  // },
  publicPath: "/",
});
