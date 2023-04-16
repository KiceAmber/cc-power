<script setup lang="ts">
import SiderBar from "@/components/user/SiderBar.vue";

</script>

<template>
    <el-container>
        <el-container>
            <el-aside width="auto">
                <SiderBar></SiderBar>
            </el-aside>
            <el-main>
                <router-view />
            </el-main>
        </el-container>
    </el-container>
</template>

<style scoped lang="less">
.el-container {
    height: 100vh;
}

.el-header {
    width: auto;
    height: 60px;
    padding: 0;
}

.el-aside {
    width: 150px;
}

.el-main {
    background-color: #E9EEF3;
}
</style>