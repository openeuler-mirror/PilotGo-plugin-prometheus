<template>
  <div class="top">
    <span class="top-title">告警管理</span>
  </div>
  <div class="list shadow">
    <pm-table ref="ruleTableRef" :show-check="false" :show-search="false" :get-data="getHistoryAlerts"
      :get-all-data="getHistoryAlerts" @handleSelect="handleSelect" @handleRowclick="handleRowclick"
      @handleAllCheckHost="handleAllCheckHost">
      <el-table-column type="selection" :reserve-selection="true" width="50" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="ip" label="IP" width="180" />
      <el-table-column prop="departmentName" label="部门" width="180" />
      <el-table-column prop="alertName" label="告警名称" width="220" />
      <el-table-column prop="level" label="告警级别" width="180" />
      <el-table-column prop="alertTime" label="告警开始时间" width="200" />
      <el-table-column prop="alertEndTime" label="告警结束时间" width="200" />
      <el-table-column prop="handleState" label="处理状态" width="180" />
      <el-table-column label="告警状态">
        <template #default="{ row }">
          <el-tag type="success" v-if="row.alertEndTime === ''">活跃</el-tag>
          <el-tag type="info" v-else-if="row.handleState === '已完成'">已处理</el-tag>
          <el-tag type="primary" v-else>待处理</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="告警描述" :show-overflow-tooltip="true" width="300" />
    </pm-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onBeforeUnmount } from "vue";
import pmTable from "@/components/PmTable.vue";
import {
  getHistoryAlerts,
  getMetrics,
} from "@/api/prometheus";
import { ElMessage } from "element-plus";
import type { Alert } from "@/types/alert";
import { alertStore } from "@/store/alert";

const ruleTableRef: any = ref(null); // 历史告警列表
const dialog = ref(false);
onMounted(() => {
  getAllMetrics();
});
onBeforeUnmount(() => {
  alertStore().alert_state = "";
});
// 搜索配置规则
interface SearchItem {
  level: string;
  value: string;
}
const levels = ref<SearchItem[]>([]);
const states = ref<SearchItem[]>([]);
const getAllMetrics = () => {
  getMetrics().then((res) => {
    let all_level: string[];
    levels.value = [];
    let all_state: string[];
    states.value = [];
    if (res.data.code === 200) {
      all_level = res.data.data.alertLevel;
      console.log('all_level', all_level)
      all_level.forEach((item) => {
        levels.value.push({ level: item, value: item });
      });
      all_state = res.data.data.alertState;
      all_state.forEach((item) => {
        states.value.push({ level: item, value: item });
      });
    }
  });
};

// 刷新列表数据
const advanceRef: any = ref(null);

// 处理选中
const checkedIds = ref([] as number[]);
const unable_finish = ref([] as Alert[]);
const unable_confirm = ref([] as Alert[]);
// 针对全选状态下，翻页的选中数据处理
const handleAllCheckHost = (checkedRows: Alert[]) => {
  handleSelect(checkedRows, "");
};
const handleSelect = (
  rows: Alert[],
  _type: string,
  _state_change_rows?: Alert[]
) => {
  unable_confirm.value = [];
  unable_finish.value = [];
  checkedIds.value = [];
  if (rows && rows.length > 0) {
    rows.forEach((item: Alert) => {
      checkedIds.value.push(item.id);
      if (item.alertEndTime == "" || item.handleState === "已完成") {
        unable_finish.value.push(item);
      }
      if (["已确认", "已完成"].includes(item.handleState)) {
        unable_confirm.value.push(item);
      }
    });
  }
  ruleTableRef.value.changeCheckedCount(checkedIds.value.length);
};
// 取消选中的行
const handleRowclick = (rows: Alert[]) => {
  if (ruleTableRef.value) {
    rows.forEach((rowItem) => {
      ruleTableRef.value!.toggleRowSelection(rowItem, false);
    });
  }
};

// 监听概览页面跳转时携带的参数state
watch(
  () => alertStore().alert_state,
  (new_state, old_state) => {
    if (new_state) {
      nextTick(() => {
        setTimeout(() => {
          ruleTableRef.value!.handleSearch({ state: new_state, search: true });
        }, 400);
      });
    }
  },
  { immediate: true }
);
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
