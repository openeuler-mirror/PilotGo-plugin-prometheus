<template>
  <div class="pm_table">
    <!-- table工具条 -->
    <el-row class="pm_table_header">
      <div class="pm_table_header_operation">
        <!-- 全选 -->
        <el-dropdown v-if="props.showCheck" split-button class="operation-select-checkbox" @command="handleCheckChange"
          @click.stop="handleCheckClick($event, !checkedState)">
          <el-checkbox v-model="checkedState" @change.native="handleCheckState" />
          <span class="operation-select-spanText" v-show="checkedState">选择个数（{{ checkedCount }}个）</span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="不选择">不选择（0个）</el-dropdown-item>
              <el-dropdown-item command="选择当前页">{{
                "选择当前页（" + currentNum + "个）"
              }}</el-dropdown-item>
              <el-dropdown-item command="选择所有">{{
                "选择所有（" + page.total + "个）"
              }}</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <!-- 模糊搜索 -->
        <div class="operation-select-input">
          <el-autocomplete v-if="props.showSearch" v-model="keyWord" :fetch-suggestions="querySearch"
            popper-class="my-autocomplete" placeholder="请输入内容" clearable @clear="handleIpSearch" @select="handleIpSearch"
            @keyup.enter="handleInputIP">
            <template #prefix>
              <el-icon class="el-input__icon" @click="handleIpSearch({ ip: keyWord })">
                <Search />
              </el-icon>
            </template>
            <template #suffix>
              <el-icon class="el-input__icon" @click="handleArrowClick">
                <ArrowDown />
              </el-icon>
            </template>
            <template #default="{ item }">
              <div class="value">{{ item.ip }}</div>
            </template>
          </el-autocomplete>
          <slot name="search_bar"></slot>
        </div>
      </div>
      <div class="pm_table_header_button">
        <slot name="button_bar"></slot>
      </div>
    </el-row>
    <!-- 列表 -->
    <div class="pm_table_content" ref="tableBox">
      <el-table ref="multipleTableRef" :data="tableData" :row-key="(row: any) => row.id" class="table"
        @select="handleRowSelectionChange" @selection-change="handleSelectinChange" v-loading="loading">
        <slot></slot>
        <template #append>
          <slot name="append"></slot>
        </template>
        <template #empty>
          <el-empty description="暂无数据"></el-empty>
        </template>
      </el-table>
    </div>
    <!-- 分页 -->
    <div class="pm_table_page">
      <el-config-provider :locale="zhCn">
        <el-pagination v-model:current-page="page.currentPage" v-model:page-size="page.pageSize" popper-class="pagePopper"
          :page-sizes="[20, 25, 50, 75, 100]" :small="page.small" :background="page.background"
          layout="total, sizes, prev, pager, next, jumper" :total="page.total" @size-change="getTableData(true)"
          @current-change="getTableData(true)" />
      </el-config-provider>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeMount, nextTick } from "vue";
import zhCn from "element-plus/dist/locale/zh-cn.mjs";
import { ElTable } from "element-plus";
import { alertStore } from "@/store/alert";
const props = defineProps({
  showCheck: {
    type: Boolean,
    required: false,
    default: true,
  },
  showSearch: {
    type: Boolean,
    required: false,
    default: true,
  },
  getData: {
    type: Function,
    required: true,
  },
  getAllData: {
    type: Function,
    required: false,
  },
  searchFunc: {
    type: Function,
    required: false,
  },
  searchParams: {
    type: Object,
    required: false,
    default: {},
  },
});
const emit = defineEmits([
  "handleSelect",
  "handleRowclick",
  "handleAllCheckHost",
]);
const loading = ref(false);
const tableData = ref([] as Idata[]);

const multipleTableRef = ref<InstanceType<typeof ElTable>>();
const currentNum = ref(0); //复选框当前页数量
const keyWord = ref(""); // 当前机器的ip
const checkedState = ref(false);

const all_hosts = ref([] as Idata[]);
// isCheckAll：解决element组件分页问题导致的勾选全选第二页开始不生效的问题
const isCheckAll = ref(false);
// 需求：选择当前页保留之前的选中数据
const isCheckCurrent = ref(false);

const page = reactive({
  total: 0,
  currentPage: 1,
  pageSize: 20,
  small: true,
  background: true,
});
onMounted(() => {
  if (alertStore().alert_state !== "") {
  } else {
    getTableData();
  }
});
onBeforeMount(() => {
  keyWord.value = "";
});

const search_params = ref({});
interface Idata {
  id: number;
  [propName: string]: any;
}

// 获取全部不分页数据
const getAllData = () => {
  if (!props.getAllData) return;
  props.getAllData!({ paged: "false", ...search_params.value }).then(
    (res: Idata) => {
      if (res.data && res.data.code === 200) {
        all_hosts.value = res.data.data;
        page.total = res.data.total;
      }
    }
  );
};

// 刷新整个表格
const handleRefresh = () => {
  loading.value = true;
  handleCheckChange("不选择");
  page.currentPage = 1;
  keyWord.value = "";
  search_params.value = { search: "" };
  getAllData();
  props.getData!({ page: 1, size: page.pageSize }).then((res: any) => {
    if (res.data && res.data.code === 200) {
      loading.value = false;
      tableData.value = res.data.data;
      currentNum.value = res.data.data.length;
    } else {
      loading.value = false;
      tableData.value = [];
      currentNum.value = 0;
      page.total = 0;
    }
  });
};

// 获取表格数据
const getTableData = (pageChange?: boolean) => {
  if (!pageChange) {
    // 保持分页时复选框的状态
    checkedState.value = false;
  }
  loading.value = true;
  props.getData!({
    page: page.currentPage,
    size: page.pageSize,
    ...search_params.value,
  }).then((res: any) => {
    if (res.data && res.data.code === 200) {
      loading.value = false;
      tableData.value = res.data.data;
      currentNum.value = res.data.data.length;
      page.total = res.data.total;

      if (isCheckAll.value) {
        toggleSelection();
        // 针对全选状态下的翻页数据渲染
        let enableHost = [] as Array<Idata>;
        tableData.value.forEach((item) => {
          checkedRows.value.forEach((citem) => {
            if (citem.id === item.id) {
              enableHost.push(citem);
            }
          });
        });
        toggleSelection(enableHost);
      } else if (isCheckCurrent.value) {
        toggleSelection();
        // 针对选择当前页状态下的翻页数据渲染
        let enableHost = [] as Array<Idata>;
        tableData.value.forEach((item) => {
          checkedCurrentRows.value.forEach((citem) => {
            if (citem.id === item.id) {
              enableHost.push(citem);
            }
          });
        });
        toggleSelection(enableHost);
      }
    } else {
      loading.value = false;
      tableData.value = [];
      currentNum.value = 0;
      page.total = 0;
    }
  });
  // }
  if (props.getAllData) {
    getAllData();
  }
};

// 针对告警模块的搜索事件
const handleSearch = (params: Object) => {
  page.currentPage = 1;
  handleCheckChange("不选择");
  tableData.value = [];
  all_hosts.value = [];
  search_params.value = params;
  loading.value = true;
  setTimeout(() => {
    // 因为立即监听跳转携带的search参数，设置延时避免与页面一进来的时候获取table数据冲突
    getTableData();
  }, 400);
};

// 模糊搜索
const querySearch = (queryString: string, cb: Function) => {
  const results = queryString
    ? all_hosts.value.filter(createFilter(queryString))
    : all_hosts.value;
  // call callback function to return suggestion objects
  cb(results);
};
const createFilter = (queryString: string) => {
  return (host: any) => {
    return host.ip.toLowerCase().indexOf(queryString.toLowerCase()) === 0;
  };
};

// 点击下拉图标事件
const handleArrowClick = (ev: Event) => {
  // console.log('arrow', ev)
  getAllData();
};

// 输入回车事件
const handleInputIP = () => {
  handleIpSearch();
};
// 选中ip主机事件
const handleIpSearch = (item?: { ip: string }) => {
  if (item) {
    keyWord.value = item.ip;
  }
  search_params.value = {
    search: keyWord.value,
  };
  loading.value = true;
  currentNum.value = 0;
  page.currentPage = 1;
  handleCheckChange("不选择");
  getTableData();
};

// 表格被选择数据发生变化
const handleSelectinChange = (val: Idata[]) => {
  // if (val.length === 1 && val[0].agentStatus === '断开') {
  //   multipleTableRef.value!.clearSelection();
  //   console.log('判断需要关闭多选框', val)
  //   checkedState.value = false;
  // }
};

const checkedRows = ref([] as Array<Idata>); // 全选状态下的数据
const checkedCurrentRows = ref([] as Array<Idata>); // 当前页选中下的数据
const checkedCount = ref(0); // 选择的数量
// 用户点击某一行的复选框,val：当前选中的值
const handleRowSelectionChange = (val: Array<Idata>, row: Idata) => {
  if (isCheckAll.value) {
    // 针对全选状态下的选择一个修改左上角
    if (val.indexOf(row) >= 0) {
      // 选中
      checkedRows.value.push(row);
    } else {
      // 勾掉
      checkedRows.value = checkedRows.value.filter(
        (item) => item.id !== row.id
      );
    }
    emit("handleSelect", checkedRows.value, "全选", tableData.value);
    // emit('handleAllCheckHost', checkedRows.value);
    // checkedCount.value = checkedRows.value.length;
  } else if (isCheckCurrent.value) {
    // 针对全选状态下的选择一个修改左上角
    if (val.indexOf(row) >= 0) {
      // 选中
      checkedCurrentRows.value.push(row);
    } else {
      // 勾掉
      checkedCurrentRows.value = checkedCurrentRows.value.filter(
        (item) => item.id !== row.id
      );
    }
    emit("handleSelect", checkedCurrentRows.value, "全选", tableData.value);
    // emit('handleAllCheckHost', checkedCurrentRows.value);
    // checkedCount.value = checkedRows.value.length;
  } else {
    if (val && val.length > 0) {
      emit("handleSelect", val, "单行复选框");
    } else {
      emit("handleSelect", [], "单行复选框");
    }
  }
};

// 全局选择框处理
const handleCheckChange = (val: string) => {
  checkedCount.value = 0;
  checkedRows.value = [];
  multipleTableRef.value!.clearSelection();
  switch (val) {
    case "不选择":
      checkedCurrentRows.value = [];
      isCheckCurrent.value = false;
      isCheckAll.value = false;
      checkedState.value = false;
      checkedCount.value = 0;
      toggleSelection();
      emit("handleSelect", [], "不选择");
      break;
    case "选择当前页":
      isCheckCurrent.value = true;
      isCheckAll.value = false;
      checkedState.value = true;

      if (checkedCurrentRows.value.length === 0) {
        tableData.value.forEach((item) => {
          checkedCurrentRows.value.push(item);
        });
      } else {
        let checkedRows: Idata[] = JSON.parse(
          JSON.stringify(checkedCurrentRows.value)
        );
        tableData.value.forEach((item) => {
          if (checkedRows.filter((citem) => citem.id === item.id).length == 0) {
            checkedCurrentRows.value.push(item);
          }
        });
      }

      toggleSelection(tableData.value);
      emit("handleSelect", checkedCurrentRows.value, "当前页");
      break;
    default:
      // 选择所有
      checkedCurrentRows.value = [];
      tableData.value.forEach((item) => {
        checkedCurrentRows.value.push(item);
      });
      isCheckCurrent.value = false;
      isCheckAll.value = true;
      checkedState.value = true;
      checkedRows.value = all_hosts.value;
      toggleSelection(tableData.value);
      emit("handleSelect", all_hosts.value, "全选", tableData.value);
      break;
  }
};
// checkbox选中状态事件
const handleCheckState = (state: boolean) => {
  state ? handleCheckChange("选择所有") : handleCheckChange("不选择");
};

/*
 * 变更选中行复选框选中状态
 * @rows[]:需要变更选中状态的行
 * @need_hanlde_rows[]:真正需要处理的行
 *
 */
const toggleSelection = (rows?: any[]) => {
  if (rows) {
    rows.forEach((row) => {
      toggleRowSelection(row, true);
    });
  } else {
    multipleTableRef.value!.clearSelection();
  }
};

// 变更某一行的选择状态
const toggleRowSelection = (row: Idata, isCheck: boolean) => {
  // 记录：数据源更改之后，即使数据一样，也不能操作勾选,需要重新更换数据源才可
  tableData.value.forEach((item) => {
    if (item.id === row.id) {
      multipleTableRef.value!.toggleRowSelection(item, isCheck);
    }
  });
  if (!isCheck && isCheckAll.value) {
    // 全选状态下从已选中的数组删除这一项
    checkedRows.value = checkedRows.value.filter((item) => item.id !== row.id);
  }
  if (!isCheck && isCheckCurrent.value) {
    // 当前页选中状态下从已选中的数组删除这一项
    checkedCurrentRows.value = checkedCurrentRows.value.filter(
      (item) => item.id !== row.id
    );
  }
};

// 下拉框左侧部分点击事件
const handleCheckClick = (event: any, val: any) => {
  if (
    ["operation-select-spanText", "el-button el-button--primary"].includes(
      event.target.className
    )
  ) {
    val ? handleCheckChange("选择所有") : handleCheckChange("不选择");
  }
};

// 变更当前选中数据显示
const changeCheckedCount = (count: number) => {
  if (count > 0) {
    checkedState.value = true;
    checkedCount.value = count;
  } else {
    checkedState.value = false;
  }
};

defineExpose({
  toggleRowSelection,
  changeCheckedCount,
  getTableData,
  toggleSelection,
  handleSearch,
  isCheckAll,
  handleRefresh,
});

// 监听keep-alive模式下的应用状态
window.addEventListener("appstate-change", function (e: any) {
  if (e.detail.appState === "afterhidden") {
    // 子应用已卸载
  } else if (e.detail.appState === "beforeshow") {
    // 子应用即将重新渲染
  } else if (e.detail.appState === "aftershow") {
    // 子应用已经重新渲染
    getTableData();
    getAllData();
  }
});
</script>

<style lang="scss" scoped>
.pm_table {
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;
  width: 100%;
  height: 100%;

  &_header {
    width: 100%;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    &_operation {
      flex: 2;
      display: flex;
      justify-content: flex-start;

      .operation {
        &-select {
          &-spanText {
            font-size: 12px;
            padding-left: 4px;
          }

          &-input {
            margin-left: 10px;
          }
        }
      }
    }

    &_button {
      display: flex;
      justify-content: flex-end;
    }
  }

  &_content {
    height: calc(100% - 64px - 30px - 20px);

    .table {
      width: 100%;
      height: 100%;
    }
  }

  &_page {
    height: 30px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
