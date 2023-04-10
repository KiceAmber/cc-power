<script setup lang="ts">
import {onMounted, onUnmounted, reactive, ref} from "vue";
import {FormInstance} from "element-plus";
import * as echarts from "echarts";

let powerPrice = ref(123)

const tableData = [
    {
        count: '200',
        price: '3',
        uploadTime: '2023-01-02'
    },
    {
        count: '200',
        price: '4',
        uploadTime: '2023-02-01'
    },
    {
        count: '180',
        price: '5',
        uploadTime: '2023-02-05'
    },
    {
        count: '90',
        price: '6',
        uploadTime: '2023-02-03'
    },
]

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive({
    organization: '',
    selfPower: '',
    selfPrice: '',
})

let submitStatus = ref(true)

// echart 图设置
let ecahrt = echarts;

onMounted(() => {
    initChart();
})

onUnmounted(() => {
    ecahrt.dispose;
})

const initChart = () => {
    const elem = document.querySelector(".myEcharts") as HTMLElement;
    let chart = ecahrt.init(elem);
    // 柱状图
    const option = {
        xAxis: {
            type: 'category',
            data: [
                '1月',
                '2月',
                '3月',
                '4月',
                '5月',
                '6月',
                '7月',
                '8月',
                '9月',
                '10月',
                '11月',
                '12月'
            ]
        },
        yAxis: {
            type: 'value'
        },
        series: [
            {
                data: [327, 288, 271, 164, 171, 182, 190, 180, 218, 217, 290, 330],
                type: 'bar'
            }
        ]
    };

    chart.setOption(option);
    window.onresize = function () {
        //自适应大小
        chart.resize();
    };
}

</script>

<template>
    <div class="container">
        <div>
            <el-card style="width: 750px">
                <div class="post-power">🔋售卖电量资源</div>
                <div class="recommend" style="margin-top:5px;">推荐的电量价格范围：￥{{ powerPrice }}/度</div>
                <el-form
                    ref="ruleFormRef"
                    :model="ruleForm"
                    status-icon
                    label-width="120px"
                    class="demo-ruleForm"
                >
                    <el-form-item label="电量">
                        <el-input autocomplete="off" v-model="ruleForm.selfPower" placeholder="输入电量度数"/>
                    </el-form-item>
                    <el-form-item label="价格">
                        <el-input autocomplete="off" v-model="ruleForm.selfPrice" placeholder="输入价格"/>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="success" @click="">上传提交</el-button>
                    </el-form-item>
                </el-form>
            </el-card>

            <el-card class="box-card" style="width: 750px">
                <div class="recent-data">🪙近期电量在售数据</div>
                <el-table :data="tableData" style="width: 100%">
                    <el-table-column prop="count" label="数量" width="auto" sortable/>
                    <el-table-column prop="price" label="价格" width="auto"/>
                    <el-table-column prop="uploadTime" label="发布日期"/>
                </el-table>
            </el-card>
        </div>
        <div style="width:930px;">
            <el-card style="margin-left: 30px;">
                <div style="font-size: 20px; text-align:center">📈年交易数据图</div>
                <div class="myEcharts" style="width: 910px;height:565px;"></div>
            </el-card>
        </div>
    </div>
</template>

<style scoped lang="less">
.post-power {
    font-size: 30px;
    text-align: center;
}

.recommend {
    font-size: 20px;
    margin-bottom: 30px;
}

.recent-data {
    margin-bottom: 30px;
    text-align: center;
    font-size: 30px;
}

.box-card {
    margin-top: 30px;

    el-table {
        fontsize: 30px;
    }
}

.container {
    display: flex;
}
</style>