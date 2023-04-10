<template>
    <div class="login-page">
        <el-card class="box-card">
            <div slot="header" class="clearfix" style="text-align: center">
                <span class="login-title">🥝欢迎使用</span>
            </div>
            <div class="login-form">
                <el-select v-model="form.role" placeholder="请选择身份" size="default">
                    <el-option v-for="item in options" :key="item.key" :label="item.label" :value="item.value" />
                </el-select>
                <el-form :model="form" :rules="loginRules" ref="loginForm">
                    <el-form-item prop="username">
                        <el-input type="text" v-model="form.username" auto-complete="off" placeholder="请输入用户名">
                            <template slot="prepend"><i style="font-size:20px" class="el-icon-user"></i></template>
                        </el-input>
                    </el-form-item>
                    <el-form-item prop="password">
                        <el-input type="password" v-model="form.password" auto-complete="off" placeholder="请输入密码" clearable>
                            <template slot="prepend"><i style="font-size:20px" class="el-icon-key"></i></template>
                        </el-input>
                    </el-form-item>

                    <div class="replace" v-if="!isRegister">
                        <el-form-item>
                            <el-button style="width:100%;" type="primary" @click="handleLogin">登录</el-button>
                        </el-form-item>
                    </div>
                    <div class="replace" v-else>
                        <el-form-item prop="rePassword">
                            <el-input type="password" v-model="form.re_password" auto-complete="off" placeholder="请确认密码"
                                clearable>
                                <template slot="prepend"><i style="font-size:20px" class="el-icon-key"></i></template>
                            </el-input>
                        </el-form-item>
                    </div>

                    <el-form-item>
                        <el-button style="width:100%;" type="warning" @click="handleRegister">
                            {{ isRegister ? "提交" : "注册" }}
                        </el-button>
                    </el-form-item>

                    <div v-if="isRegister">
                        <el-button style="width:100%;" type="info" @click="cancel">
                            取消
                        </el-button>
                    </div>
                </el-form>
            </div>
        </el-card>
    </div>
</template>

<script setup lang="ts">
import request from "@/tools/request"
import { ElMessage } from "element-plus";
import qs from "qs";
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
let isRegister = ref(false);

const form = reactive({
    role: "",
    username: "",
    password: "",
    re_password: ""
})

const loginRules = {
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    re_password: [{ required: true, message: '请确认密码', trigger: 'blur' }]
}

const handleLogin = () => {
    request({
        url: "/login",
        method: "POST",
        data: qs.stringify({
            username: form.username,
            password: form.password,
        }),
        headers: {
            // 请求头添加类型
            "Content-Type": "application/x-www-form-urlencoded"
        }
    }).then(res => {
        if (res.data.flag === true) {
            ElMessage({
                showClose: false,
                message: "登录成功",
                type: "success"
            })
            router.push({
                path: "/user"
            })
        }
    }).catch(err => {
        ElMessage({
            showClose: false,
            message: "登录失败",
            type: "error"
        })
    })
}

const handleRegister = () => {
    if (!isRegister.value) {
        isRegister.value = !isRegister.value;
        return
    } else {
        request({
            url: "/register",
            method: "POST",
            // 将数据转为 form 表单
            data: qs.stringify({
                role: form.role,
                username: form.username,
                password: form.password,
            }),
            headers: {
                // 请求头添加类型
                "Content-Type": "application/x-www-form-urlencoded"
            }
        }).then(res => {
            if (res.data.flag === true) {
                // 注册成功
                ElMessage({
                    showClose: false,
                    message: "注册成功",
                    type: "success"
                })
                isRegister.value = !isRegister.value;
            }
        }).catch(err => {
            ElMessage({
                showClose: false,
                message: "注册失败",
                type: "error",
            })
        })
    }
}

const cancel = () => {
    isRegister.value = !isRegister.value;
}

const options = [
    {
        key: "0",
        value: "购电用户",
        label: "购电用户"
    },
    {
        key: "1",
        value: "供电商",
        label: "供电商"
    }
]
</script>

<style scoped lang="less">
.login-page {
    background: url("@/assets/user/images/bg.jpg") no-repeat right;
    background-size: cover;
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
}

.login-title {
    font-size: 20px;
}

.box-card {
    background-color: rgba(185, 214, 243, 0.8);
    width: 400px;
}

.m-2 {
    width: 360px;
    margin-top: 15px;
    margin-bottom: 15px;
}
</style>

