<template>
  <el-form ref="ruleFormRef" :model="form" :rules="rules" label="right" label-width="auto" class="custom-form">
    <el-form-item label="告警名称" prop="alertName">
      <el-input v-model="form.alertName" maxlength="15" placeholder="请输入告警名称" clearable type="text"
        show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="监控主机" :prop="selectHostType">
      <el-radio-group v-model="selectHostType" style="width: 100%" @change="changeHostType">
        <el-radio label="ips">选择主机IP</el-radio>
      </el-radio-group>
      <el-select v-model="form.ips" v-if="selectHostType === 'ips'" multiple filterable placeholder="请选择监控主机"
        collapse-tags @change="changeHost">
        <el-option v-for="item in ips" :key="item.targetIp
          " :label="item.targetIp" :value="item.targetIp" />
        <template #empty>
          <div style="text-align: center; padding: 10px">暂无数据</div>
        </template>
      </el-select>

      <el-select v-model="form.batches" v-if="selectHostType === 'batches'" multiple placeholder="请选择主机批次" collapse-tags
        @change="selectBtaches">
        <el-option v-for="item in batches" :key="item.id" :label="item.name" :value="item.id" />
        <template #empty>
          <div style="text-align: center; padding: 10px">暂无数据</div>
        </template>
      </el-select>
    </el-form-item>
    <el-form-item label="监控指标" prop="metrics">
      <el-select v-model="form.metrics" placeholder="请选择监控指标" @change="handleAlertLevel">
        <el-option v-for="item in metrics" :key="item" :label="item" :value="item" />
        <template #empty>
          <div style="text-align: center; padding: 10px">暂无数据</div>
        </template>
      </el-select>
    </el-form-item>
    <el-form-item label="告警级别" :prop="showLevelInput ? 'input_severity' : 'severity'">
      <el-radio-group v-model="select_level_type" style="width: 100%" @change="changeLevel">
        <el-radio label="select" @click="showLevelInput = false">选择告警级别</el-radio>
        <el-radio label="input" @click="showLevelInput = true">自定义告警级别</el-radio>
      </el-radio-group>
      <el-select v-model="form.severity" filterable v-if="!showLevelInput" placeholder="请选择告警级别">
        <el-option v-for="item in levels" :key="item" :label="item" :value="item" />
        <template #empty>
          <div style="text-align: center; padding: 10px">暂无数据</div>
        </template>
      </el-select>
      <el-input v-model="form.input_severity" v-if="showLevelInput" placeholder="请输入自定义级别"></el-input>
    </el-form-item>
    <el-form-item :label="'告警阈值(' + unit + ')'" prop="threshold">
      <el-input style="width: 240px" v-model="form.threshold" :disabled="form.metrics === '服务器宕机'" placeholder="请输入告警阈值"
        clearable>
      </el-input>
    </el-form-item>
    <el-form-item label="持续触发时间(s)" prop="duration">
      <el-input style="width: 240px" v-model="form.duration" placeholder="请输入持续触发时间，单位秒" clearable></el-input>
    </el-form-item>
    <el-form-item label="自定义描述" prop="desc">
      <el-input v-model="form.desc" type="textarea" placeholder="请输入自定义描述" maxlength="50" show-word-limit></el-input>
    </el-form-item>
  </el-form>
  <el-form class="centered-buttons">
    <el-button @click="onCancel">取消</el-button>
    <el-button type="primary" @click="saveTaskData(ruleFormRef)" :loading="isLoading"
      @mousedown="(e: any) => e.preventDefault()">确定</el-button>
  </el-form>
</template>

<script setup lang="ts">
import { ElMessage, ElForm } from "element-plus";
import { reactive, onMounted, ref } from "vue";
import {
  addConfigRule,
  updateConfigRule,
  getMetrics,
  getExporterList,
  getHostsByBatchId,
} from "@/api/prometheus";
import type { ConfigRule } from "@/types/rule";
import type { Host } from "@/types/host";
import { checkDuration } from "./customRule";

type FormInstance = InstanceType<typeof ElForm>;
type FormRules = InstanceType<typeof ElForm>;

const emit = defineEmits(["close", "update"]);
const props = defineProps({
  row: {
    type: Object,
    required: false,
    default: {},
  },
  isUpdate: {
    type: Boolean,
    required: false,
    defalut: false,
  },
});

const ruleFormRef = ref<FormInstance>();
const form = reactive<ConfigRule>({
  alertName: "",
  batches: [],
  batches_str: "",
  ips: [],
  desc: "",
  alertTargets: [],
  metrics: "",
  threshold: "",
  duration: "15",
  severity: "",
  input_severity: "",
  alertLabel: "",
});
const rules = reactive<FormRules>({
  alertName: [
    { required: true, message: "请输入告警名称", trigger: ["blur", "change"] },
  ],
  threshold: [
    { required: true, message: "请输入告警阈值", trigger: "blur" },
    { validator: checkDuration, trigger: ["blur", "change"] },
  ],
  duration: [
    {
      required: true,
      message: "请输入告警持续触发时间",
      trigger: ["blur", "change"],
    },
    { validator: checkDuration, trigger: ["blur", "change"] },
  ],
  ips: [{ required: true, message: "请选择监控主机", trigger: ["change"] }],
  batches: [{ required: true, message: "请选择主机批次", trigger: "change" }],
  metrics: [{ required: true, message: "请选择监控指标", trigger: "change" }],
  severity: [{ required: true, message: "请选择告警级别", trigger: "change" }],
  input_severity: [
    { required: true, message: "请输入告警级别", trigger: ["blur", "change"] },
  ],
  desc: [
    {
      required: true,
      message: "请输入自定义描述",
      trigger: ["blur", "change"],
    },
  ],
} as any);

onMounted(() => {
  let rowItem: ConfigRule = props.row as ConfigRule;
  if (rowItem && props.isUpdate) {
    selectHostType.value = "ips"
    form.id = rowItem.id;
    form.ips = rowItem.alertTargets.map((item) => item.ip);
    form.alertTargets = rowItem.alertTargets;
    form.alertName = rowItem.alertName;
    form.alertLabel = rowItem.alertLabel;
    form.desc = rowItem.desc;
    form.metrics = rowItem.metrics;
    form.duration = rowItem.duration;
    form.severity = rowItem.severity;
    handleAlertLevel(form.metrics);
    form.threshold = rowItem.threshold;
  }
  getAllHost();
  getAllMetrics();
  console.log('ruleFormRef', ruleFormRef)
});

// 获取所有批次
const selectHostType = ref("ips");
const batches = ref([] as any[]);
// 批次选中事件
const selected_batches = ref([] as number[]);
const selectBtaches = (value: number[]) => {
  if (!value) return;
  selected_batches.value = value;
  form.batches_str = value.join();
};

// 获取批次内的所有主机
const getHostsByBtaches = () => {
  form.alertTargets = [];
  let all_proms = [] as any;
  selected_batches.value.forEach((batch_id) => {
    let promissItem = getHostsByBatchId({ id: batch_id }).then((res) => {
      // 添加主机信息
      if (res.data.code == 200) {
        res.data.data.forEach((host: any) => {
          if (host.agentStatus !== "断开") {
            form.alertTargets.push({
              ip: host.ip,
              hostId: host.id,
            });
          }
        });
      }
    });
    all_proms.push(promissItem);
  });
  return Promise.all(all_proms);
};

// 获取所有监控主机
const ips = ref([] as Host[]);
const getAllHost = () => {
  getExporterList({ paged: false }).then((res) => {
    if (res.data.code === 200) {
      ips.value = res.data.data;
    }
  });
};

// 获取监控指标
// 获取所有的告警级别
const levels = ref([] as any);
const metrics = ref([] as any);
const getAllMetrics = () => {
  getMetrics().then((res) => {
    if (res.data.code === 200) {
      metrics.value = res.data.data.metrics;
      levels.value = res.data.data.ruleLevel;
    }
  });
};
// 处理主机选中
let checked_hosts = ref([] as any[]);
const changeHost = (value: string[]) => {
  if (!value) return;
  checked_hosts.value = [];
  value.forEach((ip) => {
    let ip_filtered_item = ips.value.filter((item) => item.targetIp === ip)[0];
    checked_hosts.value.push({
      hostId: ip_filtered_item.hostId,
      ip: ip_filtered_item.targetIp,
    });
  });

  form.alertTargets = checked_hosts.value;
};

// 点击选择告警级别方式时触发
const changeLevel = () => {
  form.severity = "";
  form.input_severity = "";
  setTimeout(() => {
    ruleFormRef.value.clearValidate("severity");
    ruleFormRef.value.clearValidate("input_severity");
  }, 50);
};
// 点击选择告警级别方式时触发
const changeHostType = () => {
  form.batches = [];
  form.ips = [];
  setTimeout(() => {
    ruleFormRef.value.clearValidate("ips");
    ruleFormRef.value.clearValidate("batches");
  }, 50);
};

// 根据告警指标处理告警阈值
let unit = ref("%");
const handleAlertLevel = (value: string) => {
  form.threshold = "";
  if (["cpu使用率", "内存使用率", "磁盘容量", ""].includes(value))
    unit.value = "%";
  if (["网络流入", "网络流出"].includes(value)) unit.value = "KB";

  if (["TCP连接数"].includes(value)) unit.value = "个";

  if (["服务器宕机"].includes(value)) {
    unit.value = "0";
    form.threshold = "0";
    setTimeout(() => {
      ruleFormRef.value.clearValidate("threshold");
    }, 50);
  }
};

// 新增告警级别
const select_level_type = ref("select");
const showLevelInput = ref(false);

// 新增一条规则
const isLoading = ref(false);
const addRule = () => {
  addConfigRule(form)
    .then((res) => {
      if (res.data.code === 200) {
        isLoading.value = false;
        ElMessage.success(res.data.msg);
        emit("close");
        emit("update", true);
      } else {
        isLoading.value = false;
        ElMessage.error(res.data.msg);
      }
    })
    .catch((err) => {
      ElMessage.error("数据传输失败，请检查", err);
    });
};

// 编辑一条规则
const updateRule = () => {
  updateConfigRule(form)
    .then((res) => {
      if (res.data.code === 200) {
        isLoading.value = false;
        ElMessage.success(res.data.msg);
        emit("close");
        emit("update", false);
      } else {
        isLoading.value = false;
        ElMessage.error(res.data.msg);
      }
    })
    .catch((err) => {
      ElMessage.error("数据传输失败，请检查", err);
    });
};

// 保存任务数据
const saveTaskData = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate((valid: any, _fields: any) => {
    if (valid) {
      isLoading.value = true;
      if (showLevelInput.value) {
        form.severity = form.input_severity!;
      }
      if (selectHostType.value === "batches") {
        getHostsByBtaches().then((res) => {
          props.isUpdate ? updateRule() : addRule();
        });
      } else {
        props.isUpdate ? updateRule() : addRule();
      }
    }
  });
};

const onCancel = () => {
  emit("close");
};
</script>

<style scoped lang="scss">
.custom-form {
  margin-left: 25px;
  margin-right: 20px;
  outline-width: 120px;

  .custom-input {
    white-space: pre-wrap;
    resize: none;
    text-align: left;
    vertical-align: top;
  }
}

.centered-buttons {
  text-align: right;
}
</style>
