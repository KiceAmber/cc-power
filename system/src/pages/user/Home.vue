<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessageBox, ElMessage, ElLoading } from "element-plus";
import request from "@/tools/request"
import qs from "qs";
import { useStatusStore } from "@/store/status"
import * as echarts from "echarts";

const statusStore = useStatusStore();

const userInfo = reactive({
    username: "",
    role: "",
    balance: 0,
})


onMounted(() => {
    loadUserInfo();
    loadEchartNum();
    loadEchartMoney();
})

// 加载用户信息
const loadUserInfo = () => {
    request({
        url: "/queryUserBaseInfo",
        method: "GET",
        headers: {
            // 请求头添加类型
            "Content-Type": "application/x-www-form-urlencoded"
        }
    }).then(res => {
        userInfo.username = res.data.data.username;
        userInfo.role = res.data.data.role;
        userInfo.balance = res.data.data.balance;
        statusStore.setStatus(userInfo.role);
    })
}

// 用户充值
let dialogVisible = ref(false)
const balanceForm = reactive({
    balance: ""
})

const userTopUp = () => {
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    request({
        url: "/userTopUp",
        method: "POST",
        data: qs.stringify({
            balance: balanceForm.balance,
        })
    }).then(res => {
        loading.close();
        if (res.data.flag === true) {        
            ElMessageBox.alert(
                '充值成功',
                '提示消息'
            ).then(() => {
                loadUserInfo();
            })
        } else {
            ElMessageBox.alert(
                "充值失败, " + res.data.data,
                "error"
            )
        }
    })
    // 关闭弹窗
    dialogVisible.value = false
}

// ===============图表部分===============
const loadEchartNum = () => {
    // echars 图表
    const chart = echarts.init(document.getElementById("transaction-num") as HTMLElement);
    const option = {
        title: {
            text: '本年度近4个月交易订单笔数',
            left: 'center'
        },
        tooltip: {
            trigger: 'item'
        },
        legend: {
            orient: 'vertical',
            left: 'left'
        },
        series: [
            {
            name: '数据来自',
            type: 'pie',
            radius: '50%',
            data: [
                { value: 104, name: "1月" },
                { value: 126, name: "2月" },
                { value: 112, name: "3月" },
                { value: 89, name: "4月" },
            ],
            emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }
        ]
    };
    chart.setOption(option);
}

// 初始化订单金额图表
const loadEchartMoney = () => {
    // echars 图表
    const chart = echarts.init(document.getElementById("transaction-money") as HTMLElement);
    const option = {
        xAxis: {
            type: 'category',
            data: ["1月", "2月", "3月", "4月"],
        },
        yAxis: {
            type: 'value'
        },
        series: [
            {
                data: [8900, 7860, 6720, 14420],
                type: 'bar'
            }
        ]
    };
    chart.setOption(option);
}
</script>

<template>
    <div class="container">
        <el-card class="box-card" style="width: 20%;">
            <div class="user">
                <img src="@/assets/user/images/avatar.png" />
                <div class="user-info">
                    <p class="name">{{ userInfo.username }}</p>
                    <p class="access">{{ userInfo.role }}</p>
                </div>
            </div>
        </el-card>

        <el-card class="balance" style="">
            <div class="balance-num">
                余额：{{ userInfo.balance }}
            </div>
            <div class="recharge">
                <el-button type="success" @click="dialogVisible = true">充值</el-button>
                <el-dialog
                    v-model="dialogVisible"
                    title="充值余额"
                    width="40%"
                >
                <span style="display: block; font-size: 16px;margin-bottom: 10px;">请输入要充值的金额：</span>
                <el-form :model="balanceForm">
                    <el-form-item prop="balance">
                        <el-input type="text" v-model="balanceForm.balance"/>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="userTopUp" >
                            确认
                        </el-button>
                    </el-form-item>
                </el-form>
                </el-dialog>
            </div>
        </el-card>
    </div>  

    <div class="echart-datas">
        <el-card style="margin-top: 20px; width: 40%; height: 500px;" >
            <div id="transaction-num" style="width: 600px; height: 500px;"></div>
        </el-card>

        <el-card style="margin-top: 20px; width: 40%; height: 500px; margin-left: 2%;">
            <strong><div style="text-align: center;">本年度近4个月交易金额</div></strong>
            <div id="transaction-money" style="width: 600px; height: 500px;"></div>
        </el-card>
    </div>
</template>

<style scoped lang="less">
* {
    font-size: 20px;
}
.container {
    display: flex;
}

.box-card {
    display: flex;
    .user {
        display: flex;
        align-items: center;
        width: 300px;

        img {
            margin-right: 40px;
            width: 150px;
            height: 150px;
            border-radius: 50%;
        }

        .user-info {
            .name {
                font-size: 32px;
                margin-bottom: 10px;
            }

            .access {
                color: #999999;
            }
        }
    }
}

.balance {
    width: 19%;
    font-size: 18px;
    margin-left: 1%;
    text-align: center;

    .balance-num {
        padding: 20px;
        font-size: 20px;
    }
}

.trade {
    margin-top: 20px;
    width: 82%;
    .post-power {
        font-size: 30px;
        text-align: center;
        margin-bottom: 20px;
    }
}

.echart-datas {
    display: flex;
    flex: left;
}

.month-transaction {
    text-align: center;
    width: 40%;
    margin-left: 2%;
    background-color: rgb(176, 235, 38);
}

.temp {
    text-align: center;
    margin-left: 1%;
    width: 19%;
}
</style>