<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: wanghaohao <wanghaohao@kylinos.cn>
 * Date: Mon Oct 28 17:43:48 2024 +0800
-->
<template>
  <div class="top">
    <span class="top-title">告警规则配置</span>
  </div>
  <!-- <box-draggable></box-draggable> -->
  <div class="list shadow">
    <pm-table ref="ruleTableRef" :show-check="false" :show-search="false" :get-data="getRuleList">
      <template #search_bar>
        <div style="display: flex; align-items: center">
          <span>告警级别：</span>
          <MyAutoComplete ref="autocomplete" :all_data="levels" :placeholder="'请输入告警级别'" @change="handleSearch"
            @input="searchInputKey" />
        </div>
      </template>
      <template #button_bar>
        <el-button plain class="el-button1" type="primary" @click="handleAdd"
          @mousedown="(e: any) => e.preventDefault()">新增</el-button>
      </template>
      <el-table-column prop="id" label="ID" width="140" />
      <el-table-column prop="alertName" label="告警名称" :show-overflow-tooltip="true" width="220" />
      <el-table-column prop="alertTargets" label="告警主机" :show-overflow-tooltip="true" width="260">
        <template #default="{ row }">
          <span v-for="(item, index) in row.alertTargets">
            <span v-if="index + 1 < row.alertTargets.length">{{ item.ip }} <el-divider direction="vertical" /></span>
            <span v-else>{{ item.ip }}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="severity" label="告警级别" width="220" />
      <el-table-column prop="metrics" label="告警指标" width="260" />
      <el-table-column prop="threshold" label="告警阈值">
        <template #default="{ row }">
          <span v-if="['cpu使用率', '内存使用率', '磁盘容量'].includes(row.metrics)">{{ row.threshold + "%" }}</span>
          <span v-if="['网络流入', '网络流出'].includes(row.metrics)">{{
            row.threshold + "KB"
          }}</span>
          <span v-if="['TCP连接数'].includes(row.metrics)">{{
            row.threshold + "个"
          }}</span>
          <span v-if="['服务器宕机'].includes(row.metrics)">{{
            row.threshold
          }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="duration" label="持续触发时间" width="160">
        <template #default="{ row }">
          <span>{{ row.duration + "s" }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="desc" label="告警描述" :show-overflow-tooltip="true" width="280" />
      <el-table-column fixed="right" label="操作" width="200">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleUpdate(row)"
            @mousedown="(e: any) => e.preventDefault()">编辑</el-button>
          <el-popconfirm title="请确认是否要删除该配置规则？" confirm-button-text="确定" cancel-button-text="取消" icon-color="#f00"
            confirm-button-type="danger" @confirm="handleDelete(row)">
            <template #reference>
              <el-button link type="danger" size="small" @mousedown="(e: any) => e.preventDefault()">删除</el-button>
            </template></el-popconfirm>
        </template>
      </el-table-column>
    </pm-table>
  </div>
  <el-dialog :title="title" width="40%" v-model="showDialog" destroy-on-close :center="false">
    <addRule v-if="showDialog" @close="closeDialog" :row="selectedEditRow" :is-update="isUpdate"
      @update="handleRefresh" />
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import pmTable from "@/components/PmTable.vue";
import MyAutoComplete from "@/components/MyAutoComplete.vue";
import addRule from "./ruleForm/addRule.vue";
import { getRuleList, delConfigRule, getMetrics } from "@/api/prometheus";
import type { ConfigRule } from "@/types/rule";
import { ElMessage } from "element-plus";
import { useMacStore } from "@/store/mac";

const title = ref("");

const isUpdate = ref(false);
const showDialog = ref(false);
const ruleTableRef: any = ref(null);
const selectedEditRow = ref({} as ConfigRule);
onMounted(() => {
  getAllMetrics();
});

// 关闭弹窗
const closeDialog = () => {
  showDialog.value = false;
  title.value = "";
  selectedEditRow.value = {} as ConfigRule;
  isUpdate.value = false;
};

// 新增配置规则
const handleAdd = () => {
  title.value = "新增告警规则";
  showDialog.value = true;
  isUpdate.value = false;
};

// 编辑配置规则
const handleUpdate = (row: ConfigRule) => {
  title.value = "编辑告警规则";
  showDialog.value = true;
  isUpdate.value = true;
  selectedEditRow.value = row;
};

// 删除配置规则
const handleDelete = (row: ConfigRule) => {
  if (!row.id) return;
  delConfigRule({ id: row.id }).then((res) => {
    if (res.data.code === 200) {
      searchLevel.value = "";
      handleRefresh(true);
      ElMessage.success(res.data.msg);
    } else {
      ElMessage.error(res.data.msg);
    }
  });
};

// 搜索配置规则
const autocomplete: any = ref(null);
const searchLevel = ref("");
const levels = ref([] as { level: string; value: string }[]);
const getAllMetrics = () => {
  getMetrics().then((res) => {
    let all_level: string[];
    levels.value = [];
    if (res.data.code === 200) {
      all_level = res.data.data.ruleLevel;
      all_level.forEach((item) => {
        levels.value.push({ level: item, value: item });
      });
    }
  });
};

//
const handleSearch = (value: { value: string }, level_key: string) => {
  let search_word: string = "";
  search_word = level_key === value.value ? level_key : "";
  ruleTableRef.value!.handleSearch({ search: search_word });
};
const searchInputKey = (level_key: string) => {
  ruleTableRef.value!.handleSearch({ search: level_key });
};

// 刷新列表数据
const handleRefresh = (is_first_page: Boolean) => {
  getAllMetrics();
  searchLevel.value = "";
  if (is_first_page) {
    ruleTableRef.value!.handleRefresh();
  } else {
    ruleTableRef.value!.getTableData();
  }
};
</script>

<style scoped lang="scss">
.top {
  width: 97.4%;
  margin: 0 auto;
  height: 64px;
  display: flex;
  justify-content: space-between;
  align-items: center;

  &-title {
    font-size: 20px;
    color: #222;
    font-weight: 500;
    display: inline-block;
  }
}

.list {
  width: 98.4%;
  height: calc(100% - 64px);
  margin: 0 auto;
  background-color: #fff;
  padding: 0 20px;
}
</style>
