<script setup lang="ts">
import request from '@/tools/request';
import { ElMessage } from 'element-plus';
import { onMounted, reactive, ref } from 'vue';
import qs from "qs";

onMounted(() => {
    loadUserTableData();
    loadAssetTraceBackTable();
})

// 加载用户表数据
const loadUserTableData = () => {
    request({
        url: "/queryAllUser",
        method: "POST",
        headers: {
            // 请求头添加类型
            "Content-Type": "application/x-www-form-urlencoded"
        }
    }).then(res => {
        userTable.userData = res.data.data
    }).catch(err => {
        ElMessage({
            showClose: false,
            message: "查询用户失败",
            type: "error",
        })
    })
}

// 加载资产溯源表数据
const loadAssetTraceBackTable = () => {
    request({
        url: "/displayAllOrderWithTraceBack",
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        }
    }).then(res => {
        assetTraceBackTable.assetData = res.data.data
    }).catch(err => {
        ElMessage({
            showClose: false,
            message: "查询用户失败: " + err,
            type: "error",
        })
    })
}

let dialogTableVisible = ref(false)

const userTable = reactive({userData:[]})

const assetTraceBackTable = reactive({assetData: []})

const traceBackTable = reactive({traceBackData: []})

// 溯源按钮操作
const traceBack = (id: any, date: any) => {
    request({
        url: "/traceBackOrder",
        method: "POST",
        data: qs.stringify({
            asset_id: id,
            date: date,
        }),
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        }       
    }).then(res => {
        traceBackTable.traceBackData = res.data.data
    }).catch(err => {
        ElMessage({
            showClose: false,
            message: "查询用户失败:" + err,
            type: "error",
        })
    })
    dialogTableVisible.value = true
}

</script>

<template>
    <div class="container">
        <h1>管理员页面</h1>
        <el-tabs type="border-card" class="nav-tab">
            <el-tab-pane>
                <template #label>
                    <span class="custom-tabs-label">
                        <el-icon><user /></el-icon>
                        <span>&nbsp;&nbsp;用户管理</span>
                    </span>
                </template>
                
                <!-- 用户表 -->
                <el-table :data="userTable.userData" style="width: 100%">
                    <el-table-column prop="id" label="用户编号" width="auto" />
                    <el-table-column prop="username" label="用户名称" width="auto" />
                    <el-table-column prop="role" label="用户身份" width="auto" />                   
                </el-table>

            </el-tab-pane>
            <el-tab-pane>
                <template #label>
                    <span class="custom-tabs-label">
                        <el-icon><View /></el-icon>
                        <span>&nbsp;&nbsp;溯源查找</span>
                    </span>
                </template>
                    <!-- 资产溯源表 -->
                    <el-table :data="assetTraceBackTable.assetData" style="width: 100%">
                    <el-table-column prop="asset_id" label="资产编号" width="auto" />
                    <el-table-column prop="scope" label="交易量" width="auto" />
                    <el-table-column prop="price" label="金额" width="auto" />
                    <el-table-column prop="date" label="交易日期" width="auto" />
                    <el-table-column fixed="right" label="操作" width="auto">
                        <template #default="scope">
                            <el-button type="primary" size="large" @click="traceBack(scope.row.asset_id, scope.row.date)">
                                溯源
                            </el-button>
                        </template>
                    </el-table-column>                    
                </el-table>
            </el-tab-pane>
        </el-tabs>
    </div>

    <el-dialog v-model="dialogTableVisible" title="交易溯源记录">
        <el-table :data="traceBackTable.traceBackData">
            <el-table-column property="producer_id" label="供电商id" width="auto" />
            <el-table-column property="producer_name" label="供电商名称" width="auto" />
            <el-table-column property="producer_created_at" label="加入网络日期" width="auto" />
            <el-table-column property="consumer_id" label="购电方id" width="auto" />
            <el-table-column property="consumer_name" label="购电方名称" width="auto" />
            <el-table-column property="consumer_created_at" label="加入网络日期" width="auto" />
        </el-table>
    </el-dialog>

</template>

<style scoped lang="less">  

* {
    font-size: 20px;
}

.container {
    height: 100vh;
    background-color: rgb(255, 255, 255);
    
    h1 {
        display: block;
        margin-top: 20px;
        text-align: center;
        color: rgb(14, 22, 2);
        display: block;
        margin-bottom: 20px;
    }
    .nav-tab {
        margin: 0 40px;
        font-size: 20px;
        color: black;
    }
}

</style>