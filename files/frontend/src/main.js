import { createApp } from "vue";
import App from "./App.vue";
import router from "./router/router";
import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";

// createApp(App).use(Toast).mount("#app");

const toastOption = {
  transition: "Vue-Toastification__bounce",
  maxToasts: 5,
  newestOnTop: true,
};
const app = createApp(App);
app.use(Toast, toastOption);
app.use(router);
app.mount("#app");
