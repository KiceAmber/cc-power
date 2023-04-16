import { createRouter, createWebHashHistory } from "vue-router";

import Login from "@/pages/LoginPage.vue";

import UserLayout from "@/layout/user/UserLayout.vue";
import UserHome from "@/pages/user/Home.vue";
import ElectricityCenter from "@/pages/user/ElectricityCenter.vue";
import UserAssetManage from "@/pages/user/AssetManage.vue";
import Transaction from "@/pages/user/TransAction.vue";

import AdminHome from "@/pages/admin/AdminHome.vue";

const routes = [
    { path: "/", redirect: "/login" },
    { path: "/login", name: "login", component: Login },
    // 用户部分 界面
    {
        path: "/user",
        component: UserLayout,
        redirect: "/user/home",
        children: [
            { path: "home", name: "user-home", component: UserHome },
            { path: "elec-power", name: "user-elec-power", component: ElectricityCenter },
            { path: "asset", name: "user-asset", component: UserAssetManage },
            { path: "transaction", name: "user-transaction", component: Transaction },
        ]
    },
    // 管理员部分 界面
    {
        path: "/admin",
        component: AdminHome,
    }
]

export const router = createRouter({
    history: createWebHashHistory(),
    routes: routes
})