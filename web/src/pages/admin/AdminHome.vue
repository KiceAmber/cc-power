<script setup lang="ts">
import { reactive } from 'vue';

const userTable = reactive({
    userData:[

    ]
})
</script>

<template>
    <div class="container">
        <center><h1>管理员页面</h1></center>
        <el-tabs type="border-card" class="nav-tab">
            <el-tab-pane>
                <template #label>
                    <span class="custom-tabs-label">
                        <el-icon><user /></el-icon>
                        <span>&nbsp;&nbsp;用户管理</span>
                    </span>
                </template>
                
                <!-- 用户表内容 -->
                <el-table :data="userTable.userData" style="width: 100%">
                    <el-table-column prop="id" label="用户编号" width="auto" />
                    <el-table-column prop="name" label="用户名称" width="auto" />
                    <el-table-column prop="role" label="用户身份" width="auto" />
                    <el-table-column prop="count" label="发布的交易量" width="auto" />
                </el-table>

            </el-tab-pane>
            <el-tab-pane>
                <template #label>
                    <span class="custom-tabs-label">
                        <el-icon><View /></el-icon>
                        <span>&nbsp;&nbsp;溯源查找</span>
                    </span>
                </template>
            </el-tab-pane>
        </el-tabs>
    </div>
    
</template>

<style scoped lang="less">  
.container {
    height: 100vh;
    background-color: rgb(255, 255, 255);
    
    .nav-tab {
        margin: 0 40px;
        color: black;
    }
}

h1 {
    display: block;
    margin-bottom: 20px;
}
</style>