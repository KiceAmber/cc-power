<script setup lang="ts">

</script>

<template>
    <el-container>
        <el-header>
            头部栏
        </el-header>
        <el-container>
            <el-aside width="auto">
                侧边栏
            </el-aside>
            <el-main>
                <router-view/>
            </el-main>
        </el-container>
    </el-container>
</template>

<style scoped lang="less">

</style>