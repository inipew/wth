import Vue from "vue";
import Router from "vue-router";
import App from "../App.vue";
import EditFile from "../components/EditFile.vue"; // pastikan path benar

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: "/",
      name: "Home",
      component: App,
    },
    {
      path: "/edit",
      name: "EditFile",
      component: EditFile,
      props: (route) => ({
        fileName: route.query.fileName,
        filePath: route.query.filePath,
      }),
    },
  ],
});
