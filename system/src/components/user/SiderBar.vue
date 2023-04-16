<script setup lang="ts">
import {
    House,
    Coin,
    Van,
    Clock,
} from "@element-plus/icons-vue";

import {useRouter} from "vue-router";

const router = useRouter();

// 退出登陆
const logout = () => {
    router.push({
        path: "/"
    })
}

</script>

<template>
    <el-menu
        active-text-color="#ffd04b"
        background-color="#545c64"
        class="el-menu-vertical-demo"
        default-active="2"
        text-color="#fff"
        router
    >
        <el-menu-item index="home">
            <template #title>
                <el-icon>
                    <House/>
                </el-icon>
                <span>个人中心</span>
            </template>
        </el-menu-item>
        <el-menu-item index="elec-power">
            <template #title>
                <el-icon>
                    <Van/>
                </el-icon>
                <span>电力中心</span>
            </template>
        </el-menu-item>

        <el-menu-item index="asset">
            <template #title>
                <el-icon>
                    <Coin/>
                </el-icon>
                <span>发布记录</span>
            </template>
        </el-menu-item>

        <el-menu-item index="transaction">
            <template #title>
                <el-icon>
                    <Clock/>
                </el-icon>
                <span>交易记录</span>
            </template>
        </el-menu-item>
        
        <el-button style="margin-left: 25px;" type="info" @click="logout">
            退出登陆
        </el-button>
    </el-menu>
</template>

<style scoped lang="less">
.el-menu-vertical-demo {
    height: 100%;
}
</style>