<script setup lang="ts">

import { reactive, ref, onMounted } from "vue";
import { ElMessageBox, ElLoading, ElMessage, FormInstance } from "element-plus";
import request from "@/tools/request";
import { useStatusStore } from "@/store/status";
import qs from "qs";

onMounted(() => {
    queryAllOrder();
})

const statusStore = useStatusStore();
const tableData = reactive({ electricity: [] })

const formData = reactive({
    Scope: "",
    Price: "",
})

// 提交求购订单
const submitPurchaseOrder = () => {
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    request({
        url: "/publishBuyOrder",
        method: "POST",
        data: qs.stringify({
            scope: formData.Scope,
            price: formData.Price,
        }),
        headers: {
            // 添加类型
            "Content-Type": "application/x-www-form-urlencoded"
        }
    }).then(res => {
        loading.close();
        if (res.data.flag === true) {
            ElMessageBox.alert(
                '求购订单发布成功',
                '提示消息'
            )
        } else {
            ElMessageBox.alert(
                res.data.data,
                "提示消息"
            )
        }
    })
}

// 查询所有订单
const queryAllOrder = () => {
    request({
        url: "/queryAllOrder",
        method: "GET"
    }).then(res => {
        tableData.electricity = res.data.data;
    }).catch(err => {
        alert(err.message);
    })
}

// 上传提交电量售卖数据
const sellElectricity = () => {
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    request({
        url: "/sellElectricity",
        method: "POST",
        data: qs.stringify({
            scope: formData.Scope,
            price: formData.Price,
        }),
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        }
    }).then(res => {
        loading.close();
        if (res.data.flag === true) {
            ElMessageBox.alert(
                '电量数据发布成功',
                '提示消息'
            ).then(res => {
                queryAllOrder();
            })
        } else {
            ElMessageBox.alert(
                "电量数据发布失败",
                "提示消息"
            ).then(res => {
                queryAllOrder();
            })
        }
    })
}

</script>

<template>
    <div class="container">
        <el-card class="box-card" style="width: 100%; height: 500px">
            <div class="recent-data">🪙近期电量在售数据</div>
            <el-table :data="tableData.electricity" >
                <el-table-column prop="scope" label="数量" width="auto" />
                <el-table-column prop="price" label="价格" width="auto" />
                <el-table-column prop="date" label="发布日期" />
            </el-table>
        </el-card>

        <el-card class="trade" v-if="statusStore.status === '购电用户'" style="width: 100%; height:300px">
            <div class="post-power">🧾 发布求购电量订单</div>
            <el-form :model="formData" status-icon label-width="60px" class="demo-ruleForm">
                <el-form-item label="电量">
                    <el-input v-model="formData.Scope" autocomplete="off" placeholder="输入电量度数" />
                </el-form-item>
                <el-form-item label="价格">
                    <el-input v-model="formData.Price" autocomplete="off" placeholder="输入价格" />
                </el-form-item>
                <el-form-item>
                    <el-button type="success" @click="submitPurchaseOrder">上传提交</el-button>
                </el-form-item>
            </el-form>
    </el-card>

    <el-card class="trade" style="width: 100%; height:300px" v-else>
        <div class="post-power">🔋 发布电量资源订单</div>
            <el-form :model="formData" status-icon label-width="60px" class="demo-ruleForm">
                <el-form-item label="电量">
                    <el-input autocomplete="off" v-model="formData.Scope" placeholder="输入电量度数" />
                </el-form-item>
                <el-form-item label="价格">
                    <el-input autocomplete="off" v-model="formData.Price" placeholder="输入价格" />
                </el-form-item>
                <el-form-item>
                    <el-button type="success" @click="sellElectricity">上传提交</el-button>
                </el-form-item>
            </el-form>
    </el-card>
    </div>
</template>

<style scoped lang="less">

* {
    font-size: 20px;
}

.container {
    margin: 0;
    padding: 0;
}

.recent-data {
    margin-bottom: 30px;
    text-align: center;
    font-size: 30px;
}

.trade {
    margin-top: 30px;

    .post-power {
        font-size: 30px;
        text-align: center;
        margin-bottom: 30px;
    }

    .recommend {
        font-size: 20px;
        margin-bottom: 30px;
    }
}

.info-text {
    font-size: 22px;
    margin-bottom: 20px;
}
</style>