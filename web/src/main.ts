import {createApp} from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElIcons from '@element-plus/icons-vue'
import axios from "@/tools/request";
import {createPinia} from "pinia"
import {router} from "@/router";
import * as echarts from "echarts"; // 引入 echarts

const app = createApp(App)
const pinia = createPinia()

// 使用 ElIcons 图标
for (const name in ElIcons) {
    app.component(name, (ElIcons as any)[name])
}

app.use(router)
app.use(ElementPlus)
app.use(pinia)
app.config.globalProperties.$axios = axios;
app.config.globalProperties.$echarts = echarts;

app.mount('#app')
