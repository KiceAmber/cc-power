<script setup lang="ts">

import { reactive, ref, onMounted } from "vue";
import { ElMessageBox, FormInstance, ElMessage } from "element-plus";
import request from "@/tools/request";
import { useStatusStore } from "@/store/status";
import qs from "qs";

onMounted(() => {
    request({
        url: "/queryAllOrder",
        method: "GET"
    }).then(res => {
        tableData.electricity = res.data.data;
    }).catch(err => {
        alert(err.message);
    })
})

const statusStore = useStatusStore();
const tableData = reactive({ electricity: [] })

let powerPrice = ref(123)
const ruleFormRef = ref<FormInstance>()
const formData = reactive({
    Scope: "",
    Price: "",
})

// 上传提交电量售卖数据
const sellElectricity = () => {
    request({
        url: "/sellElectricity",
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
        // 发布成功
        ElMessageBox.alert(
            '电量数据发布成功',
            '提示消息'
        )
    }).catch(res => {
        ElMessage({
            showClose: false,
            message: "发布电量数据失败",
            type: "error"
        })
    })
}

</script>

<template>
    <div class="container">
        <el-card class="box-card">
            <div class="recent-data">🪙近期电量在售数据</div>
            <el-table :data="tableData.electricity" style="width: 100%" height="350px">
                <el-table-column prop="scope" label="数量" width="auto" />
                <el-table-column prop="price" label="价格" width="auto" />
                <el-table-column prop="date" label="发布日期" />
            </el-table>
        </el-card>

        <el-card class="trade">
            <div class="post-power">🔋售卖电量资源</div>
            <div class="recommend" style="margin-top:5px;">推荐的电量价格范围：￥{{ powerPrice }}/度</div>
            <el-form ref="ruleFormRef" :model="formData" status-icon label-width="120px" class="demo-ruleForm">
                <el-form-item label="电量">
                    <el-input autocomplete="off" v-model="formData.Scope" placeholder="输入电量度数" />
                </el-form-item>
                <el-form-item label="价格">
                    <el-input autocomplete="off" v-model="formData.Price" placeholder="输入价格" />
                </el-form-item>
                <el-form-item>
                    <el-button type="success" @click="sellElectricity" :disabled="statusStore.status !== '供电商'">上传提交</el-button>
                </el-form-item>
            </el-form>
        </el-card>
    </div>
</template>

<style scoped lang="less">
.recent-data {
    margin-bottom: 30px;
    text-align: center;
    font-size: 30px;
}

.trade {
    margin-top: 30px;
    height: 350px;

    .post-power {
        font-size: 30px;
        text-align: center;
    }

    .recommend {
        font-size: 20px;
        margin-bottom: 30px;
    }
}
</style>