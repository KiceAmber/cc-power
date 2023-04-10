import { createRouter, createWebHashHistory } from "vue-router";

import Login from "@/pages/LoginPage.vue";

import UserLayout from "@/layout/user/UserLayout.vue";
import UserHome from "@/pages/user/Home.vue";
import ElectricityCenter from "@/pages/user/ElectricityCenter.vue";
import UserAssetManage from "@/pages/user/AssetManage.vue";

import AdminLayout from "@/layout/admin/AdminLayout.vue";
import AdminHome from "@/pages/admin/HomePage.vue";
import AdminUserManage from "@/pages/admin/UserManagePage.vue";

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
        ]
    },
    // 管理员部分 界面
    {
        path: "/admin",
        component: AdminLayout,
        redirect: "/admin/home",
        children: [
            { path: "home", name: "admin-home", component: AdminHome },
            { path: "user-manage", name: "admin-user-manage", component: AdminUserManage }
        ]
    }
]

export const router = createRouter({
    history: createWebHashHistory(),
    routes: routes
})