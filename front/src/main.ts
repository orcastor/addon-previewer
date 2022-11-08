import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

import { setCallback } from "@/api/index";

import 'element-plus/es/components/message/style/css';
import { ElMessage } from 'element-plus';
setCallback(ElMessage);

const app = createApp(App);

app.mount('#app');
