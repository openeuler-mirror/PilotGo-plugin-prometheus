<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: wanghaohao <wanghaohao@kylinos.cn>
 * Date: Wed Oct 30 14:47:16 2024 +0800
-->
<template>
  <div style="width: 700px; margin: 0 auto;">
    <el-form ref="searchRef" :model="form" label-width="auto" style="max-width: 600px">
      <el-form-item label="IP地址">
        <el-input v-model="form.ip" clearable />
      </el-form-item>
      <el-form-item label="告警名称">
        <el-input v-model="form.alertName" clearable />
      </el-form-item>
      <el-form-item label="告警级别">
        <el-select v-model="form.level" size="" placeholder="请选择告警级别" clearable style="width:240px;">
          <el-option v-for="item in levels" :label="item.level" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="告警开始时间">
        <el-config-provider :locale="zhCn">
          <el-date-picker v-model="startTime" @change="handleStart" type="datetimerange" start-placeholder="开始时间"
            end-placeholder="结束时间" /></el-config-provider>
      </el-form-item>
      <el-form-item label="告警结束时间">
        <el-config-provider :locale="zhCn">
          <el-date-picker v-model="endTime" @change="handleEnd" type="datetimerange" start-placeholder="开始时间"
            end-placeholder="结束时间" /></el-config-provider>
      </el-form-item>
      <el-form-item label="处理状态">
        <el-input v-model="form.handleState" clearable />
      </el-form-item>
      <el-form-item label="告警状态">
        <el-select v-model="form.state" placeholder="请选择告警状态" clearable style="width:240px;">
          <el-option v-for="item in states" :label="item.level" :value="item.value" />
        </el-select>
      </el-form-item>
    </el-form>
    <div style="text-align: right;">
      <el-button @click="onCancle">取消</el-button>
      <el-button type="primary" @click="onReset">重置</el-button>
      <el-button type="primary" @click="onSubmit">搜索</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, watch, nextTick } from 'vue'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs';
import { getMetrics } from '@/api/prometheus';
import { alertStore } from '@/store/alert';

onMounted(() => {
  getAllMetrics();
})

const form = reactive({
  ip: '',
  state: '',
  alertName: '',
  handleState: '',
  level: '',
  alertStart: '',
  alertEnd: '',
  search: true,
})

// 搜索配置规则
interface SearchItem {
  level: string;
  value: string
}
const levels = ref<SearchItem[]>([]);
const states = ref<SearchItem[]>([]);
const getAllMetrics = () => {
  getMetrics().then(res => {
    let all_level: string[];
    levels.value = [];
    let all_state: string[];
    states.value = [];
    if (res.data.code === 200) {
      all_level = res.data.data.alertLevel;
      all_level.forEach(item => {
        levels.value.push({ level: item, value: item });
      })
      all_state = res.data.data.alertState;
      all_state.forEach(item => {
        states.value.push({ level: item, value: item });
      })
    }
  })
}

const emit = defineEmits(['search', 'cancle'])

// 时间
const startTime = ref();
const endTime = ref();
const handleStart = (value: any) => {
  if (!value) return;
  let start = new Date(value[0]).getTime() + '';
  let end = new Date(value[1]).getTime() + '';
  // 时间转化成json字符串
  form.alertStart = JSON.stringify({ start, end });
}
const handleEnd = (value: any) => {
  if (!value) return;
  let start = new Date(value[0]).getTime() + '';
  let end = new Date(value[1]).getTime() + '';
  // 时间转化成json字符串
  form.alertEnd = JSON.stringify({ start, end });
}

// 取消搜索
const onCancle = () => {
  emit('cancle')
}

// 确认搜索
const onSubmit = () => {
  emit('search', form);
}

// 重置
const searchRef: any = ref(null)
const onReset = () => {
  form.ip = '';
  form.alertEnd = '';
  form.alertStart = '';
  form.alertName = '';
  form.handleState = '';
  form.level = '';
  // form.search = false;
  form.state = '';
  startTime.value = '';
  endTime.value = '';
  alertStore().$reset();
}

// 监听概览页面跳转时携带的参数state
watch(() => alertStore().alert_state, (new_state) => {
  form.state = new_state;

}, { immediate: true })


defineExpose({
  onReset
})
</script>

<style scoped></style>