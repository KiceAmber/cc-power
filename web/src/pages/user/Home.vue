<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessageBox, ElMessage } from "element-plus";
import request from "@/tools/request"
import qs from "qs";
import { useStatusStore } from "@/store/status"

const statusStore = useStatusStore();

const userInfo = reactive({
    username: "",
    role: "",
    balance: 0,
})

const formData = reactive({
    scope: "",
    price: "",
})

// 发布求购订单
onMounted(() => {
    loadUserInfo();
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
        console.log(res)
        userInfo.username = res.data.data.username;
        userInfo.role = res.data.data.role;
        userInfo.balance = res.data.data.balance;
        statusStore.setStatus(userInfo.role);
    }).catch(err => {
        console.log("Error is", err);
    })
}

// 提交求购订单
const submitPurchaseOrder = () => {
    request({
        url: "/publishBuyOrder",
        method: "POST",
        data: qs.stringify({
            scope: formData.scope,
            price: formData.price,
        }),
        headers: {
            // 添加类型
            "Content-Type": "application/x-www-form-urlencoded"
        }
    }).then(() => {
        // 发布成功
        ElMessageBox.alert(
            '求购订单发布成功',
            '提示消息'
        )
    }).catch(() => {
        ElMessage({
            showClose: false,
            message: "发布求购订单失败",
            type: "error"
        })
    })
}

// 用户充值
const balanceForm = reactive({
    balance: ""
})
let dialogVisible = ref(false)
const userTopUp = () => {
    console.log("balanceForm:", balanceForm.balance)
    request({
        url: "/userTopUp",
        method: "POST",
        data: qs.stringify({
            balance: balanceForm.balance,
        })
    }).then(res => {
        ElMessageBox.alert(
            '充值成功',
            '提示消息'
        )
    }).catch(err => {
        ElMessageBox.alert(
            "充值失败",
            "error"
        )
    })

    // 关闭弹窗
    dialogVisible.value = false
}
</script>

<template>
    <div class="container">
        <el-card class="box-card" style="width: 40%;">
            <div class="user">
                <img src="@/assets/user/images/avatar.png" />
                <div class="user-info">
                    <p class="name">{{ userInfo.username }}</p>
                    <p class="access">{{ userInfo.role }}</p>
                </div>
            </div>
        </el-card>

        <el-card class="balance">
            <div class="balance-num">
                余额：{{ userInfo.balance }}
            </div>
            <div class="recharge">
                <el-button type="success" @click="dialogVisible = true">充值</el-button>
                <el-dialog
                    v-model="dialogVisible"
                    title="充值余额"
                    width="30%"
                >
                <span style="display: block; font-size: 16px;margin-bottom: 10px;">请输入要充值的金额：</span>
                <el-form :model="balanceForm">
                    <el-form-item prop="balance">
                        <el-input type="text" v-model="balanceForm.balance"/>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="userTopUp">
                            充值
                        </el-button>
                    </el-form-item>
                </el-form>
                </el-dialog>
            </div>
        </el-card>
        <div class="watch-info">
            注意： <br>
                &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;在本系统中，购电用户无法发布电力数据在网络上，<br>
                供电商无法发布求购订单。
        </div>
    </div>  

    <el-card class="trade">
        <div class="post-power">发布求购电量信息</div>
        <el-form status-icon label-width="120px" class="demo-ruleForm">
            <el-form-item label="电量">
                <el-input v-model="formData.scope" autocomplete="off" placeholder="输入电量度数" />
            </el-form-item>
            <el-form-item label="价格">
                <el-input v-model="formData.price" autocomplete="off" placeholder="输入价格" />
            </el-form-item>
            <el-form-item>
                <el-button type="success" :disabled="statusStore.status === '供电商'" @click="submitPurchaseOrder">
                    上传提交 
                </el-button>
            </el-form-item>
        </el-form>
    </el-card>
</template>

<style scoped lang="less">
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
    width: 20%;
    font-size: 18px;
    margin-left: 20px;
    text-align: center;

    .balance-num {
        padding: 20px;
    }

}

.trade {
    margin-top: 20px;

    .post-power {
        font-size: 30px;
        text-align: center;
        margin-bottom: 20px;

    }
}

.watch-info {
    font-size: 20px;
    margin: 10px;
}
</style>