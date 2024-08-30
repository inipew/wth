import { createWebHistory, createRouter } from "vue-router";
import HomeView from "@/view/HomeView.vue";
import Edit from "../view/EditView.vue";

const routes = [
  {
    name: "Home",
    path: "/",
    component: HomeView,
  },
  {
    name: "EditItem",
    path: "/edit/:filepath",
    component: Edit,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
