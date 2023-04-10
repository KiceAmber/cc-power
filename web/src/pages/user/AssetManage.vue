<script setup lang="ts">

import { reactive, onMounted } from 'vue';
import {useStatusStore} from "@/store/status"
import request from "@/tools/request";
import { ElMessageBox } from 'element-plus';

const statusStore = useStatusStore();
onMounted(() => {
    loadElectricityData()
})

// 供电商：显示自己发布的电量数据
// 购电方：显示自己发布的电量求购订单
const tableData = reactive({electricityData: []})

const loadElectricityData = () => {
    let queryUrl = "";
    if (statusStore.status === "购电用户") {
        queryUrl = "/querySelfBuyOrder";
    } else {
        queryUrl = "/queryUserSelfElectricity";
    }
    request({
        url: queryUrl,
        method: "POST"
    }).then(res => {
        tableData.electricityData = res.data.data;
    }).catch(err => {
        alert(err.message);
    })

}

const deleteData = () => {}

const watchDetail = () => {
    ElMessageBox.alert(
        `<table style="width: 100%;height: 220px; font-size: 15px;">
            <tr>
                <th>交易量</th>    
                <th>金额</th>    
                <th>时间</th>    
            </tr>
            <tr>
                <td>120</td>
                <td>1211</td>
                <td>2002-02-12 13:01:39</td>
            </tr>
            <tr>
                <td>120</td>
                <td>2400</td>
                <td>2002-02-12 13:01:39</td>
            </tr>
            <tr>
                <td>120</td>
                <td>2400</td>
                <td>2002-02-12 13:01:39</td>
            </tr>
        </table>`,
        "交易细节",
        {
            dangerouslyUseHTMLString: true
        }
    )
}

const cancelOrder = () => {}

</script>

<template>
    <el-card style="text-align: center;" v-if="statusStore.status === '供电商'">
        <div style="font-size: 30px; margin-bottom: 10px;">🗒已发布的交易记录(供电商)</div>
        <el-table :data="tableData.electricityData" style="width: 100%;height: 700px; font-size: 20px;">
            <el-table-column prop="id" label="订单编号" width="auto"/>
            <el-table-column prop="scope" label="数量" width="auto" />
            <el-table-column prop="price" label="价格" width="auto" />
            <el-table-column prop="date" label="发布日期" />
            <el-table-column prop="status" label="交易状态" />
            <el-table-column prop="operation" label="操作">
                <template #default>
                    <el-button link type="primary" size="default" @click="deleteData">下架</el-button>
                </template>
            </el-table-column>
        </el-table>
    </el-card>

    <el-card style="text-align: center; margin-top: 20px;" v-else>
        <div style="font-size: 30px; margin-bottom: 10px;">📦求购的交易记录(购电方)</div>
        <el-table :data="tableData.electricityData" style="width: 100%;height: 700px; font-size: 20px;">
            <el-table-column prop="id" label="订单编号" width="auto"/>
            <el-table-column prop="scope" label="数量" width="auto" />
            <el-table-column prop="price" label="价格" width="auto" />
            <el-table-column prop="date" label="发布日期" />
            <el-table-column prop="status" label="交易状态" />
            <el-table-column prop="operation" label="操作">
                <template #default>
                    <el-button link type="primary" size="default"
							@click="deleteData">支付
                    </el-button>

                    <el-button link type="primary" size="default"
							@click="cancelOrder">撤销
                    </el-button>
                    <el-button link type="primary" size="default"
							@click="watchDetail">详情
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
    </el-card>
</template>

<style scoped lang="less"></style>
