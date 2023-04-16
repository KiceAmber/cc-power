import axios from 'axios'

const request = axios.create({
    baseURL: 'http://localhost:3000/', // 后端数据基本地址
    timeout: 50000, // 请求超时设置
})

export default request