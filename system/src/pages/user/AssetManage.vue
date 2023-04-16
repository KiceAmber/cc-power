<script setup lang="ts">

import { reactive, onMounted, ref } from 'vue';
import {useStatusStore} from "@/store/status"
import request from "@/tools/request";
import { ElMessageBox, ElMessage, ElLoading } from 'element-plus';
import qs from 'qs';

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
    })
}

const detailTableData = reactive({detail:[]})
let dialogTableVisible = ref(false)

const watchDetail = (id: any) => {
    dialogTableVisible.value = true
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    if (statusStore.status === "供电商" ) {
        request({
            url: "/findAssetRecord",
            method: "POST",
            data: qs.stringify({
                asset_id: id,
            })
        }).then(res => {
            loading.close();
            detailTableData.detail = res.data.data
        })
    } else {
        request({
            url: "/findOrderRecord",
            method: "POST",
            data: qs.stringify({
                order_id: id,
            })
        }).then(res => {
            loading.close();
            detailTableData.detail = res.data.data
        })
    }
}

// 撤销交易订单
const cancelOrder = (id: any) => {
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    if (statusStore.status === "供电商") {
        request({
            url: "/sellUserCancelOrder",
            method: "POST",
            data: qs.stringify({
                assert_id: id,
            }),
        }).then(res => {    
            loading.close();
            if (res.data.flag === true) {
                ElMessageBox.alert(
                    '撤销成功',
                    '提示消息',
                ).then(() => {
                    loadElectricityData();
                })
            } else {
                ElMessageBox.alert(
                    "撤销失败，该交易可能已完成或撤销，请刷新页面",
                    "提示信息"
                )
            } 
        })
    } else {
        request({
            url: "/buyUserCancelOrder",
            method: "POST",
            data: qs.stringify({
                order_id: id,
            }),
        }).then(res => {    
            loading.close();
            if (res.data.flag === true) {
                ElMessageBox.alert(
                    '撤销成功',
                    '提示消息'
                ).then(() => {
                    loadElectricityData();
                })
            } else {
                ElMessageBox.alert(
                    "撤销失败，可能该交易已撤销，请刷新页面",
                    "提示信息"
                )
            }
        })
    }
}

// 购电用户支付订单
const payOrder = (id: any) => {
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    request({
        url: "/payOrder",
        method: "POST",
        data: qs.stringify({
            order_id: id,
        })
    }).then(res => {
        loading.close();
        if (res.data.flag === true) {
            ElMessageBox.alert(
                res.data.data,
                '提示消息'
        ).then(() => {
            loadElectricityData();
        })
        } else {
            ElMessageBox.alert(
                res.data.data,
                "提示信息"
            )
        }
    })
}

</script>

<template>
    <el-card style="text-align: center;" v-if="statusStore.status === '供电商'">
        <div style="font-size: 30px; margin-bottom: 10px;">🗒已发布的交易记录(供电商)</div>
        <el-table :data="tableData.electricityData" style="width: 100%;height: 80vh; font-size: 20px;">
            <el-table-column prop="id" label="订单编号" width="auto"/>
            <el-table-column prop="scope" label="数量" width="auto" />
            <el-table-column prop="price" label="价格" width="auto" />
            <el-table-column prop="date" label="发布日期" />
            <el-table-column prop="status" label="交易状态" />
            <el-table-column prop="operation" label="操作">
                <template #default="scope">
                    <el-button  type="danger" size="large" 
                            @click="cancelOrder(scope.row.id)">撤销
                    </el-button>
                    <el-button  type="primary" size="large"
							@click="watchDetail(scope.row.id)">详情
                    </el-button>
                    <el-dialog v-model="dialogTableVisible" title="供电商交易记录数据">
                        <el-table :data="detailTableData.detail">
                            <el-table-column property="scope" label="交易量" width="auto" />
                            <el-table-column property="price" label="金额" width="auto" />
                            <el-table-column property="date" label="交易日期" width="auto" />
                        </el-table>
                    </el-dialog>
                </template>
            </el-table-column>
        </el-table>
    </el-card>

    <el-card style="text-align: center;" v-else>
        <div style="font-size: 30px; margin-bottom: 20px;">📦求购的交易记录(购电用户)</div>
        <el-table :data="tableData.electricityData" style="width: 100%;height: 700px; font-size: 20px;">
            <el-table-column prop="id" label="订单编号" width="auto"/>
            <el-table-column prop="scope" label="数量" width="auto" />
            <el-table-column prop="price" label="价格" width="auto" />
            <el-table-column prop="date" label="发布日期" />
            <el-table-column prop="status" label="交易状态" />
            <el-table-column prop="operation" label="操作">
                <template #default="scope">
                    <el-button type="success" size="large"
							@click="payOrder(scope.row.id)">支付
                    </el-button>

                    <el-button type="danger" size="large"
							@click="cancelOrder(scope.row.id)">撤销
                    </el-button>
                    <el-button type="primary" size="large"
							@click="watchDetail(scope.row.id)">详情
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
    </el-card>

    <el-dialog v-model="dialogTableVisible" title="购电用户交易记录数据">
        <el-table :data="detailTableData.detail">
            <el-table-column property="scope" label="交易量" width="200" />
            <el-table-column property="price" label="金额" width="200" />
            <el-table-column property="date" label="交易日期"/>
        </el-table>
    </el-dialog>
</template>

<style scoped lang="less">

* {
    font-size: 20px;
}

.info-text {
    font-size: 22px;
    margin: 0px 0px;
}
</style>
