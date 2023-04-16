<template>
    <h1>注意：可查看用户历史的交易记录</h1>
    <br>
    <div class="container">
        <!-- 用户表 -->
        <el-table :data="tableData.electricity"  style="width: 100%">
            <el-table-column prop="scope" label="交易电量" width="auto" />
            <el-table-column prop="cash" label="交易金额" width="auto" />        
            <el-table-column 
                prop="date" 
                label="交易时间" 
                sortable
                column-key="date"
                width="auto" 
                :filters="[
                    { text: '2001-01-09', value: '2001-01-09' },
                    { text: '2016-05-02', value: '2016-05-02' },
                    { text: '2016-05-03', value: '2016-05-03' },
                    { text: '2016-05-04', value: '2016-05-04' },
                ]"
                />                 
        </el-table>
    </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive} from "vue"
import { TableColumnCtx, ElLoading, ElMessageBox, rowProps } from 'element-plus'
import request from "@/tools/request";

onMounted(() => {
  loadData();
})

const tableData = reactive({ electricity: [] })

const loadData = () => {
  const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
  })
  request({
        url: "/queryAllRecords",
        method: "POST",
        headers: {
            // 添加类型
            "Content-Type": "application/x-www-form-urlencoded"
        }
    }).then(res => {
        loading.close();
        if (res.data.flag === true) {
            tableData.electricity = res.data.data
        } else {
            ElMessageBox.alert(
                res.data.data,
                "提示消息"
            )
        }
    }).catch(err =>  {
        console.log(err)
    })
}

</script>

<style lang="less">
* {
    font-size: 20px;
}

.select {
    margin-bottom: 15px;
}

</style>