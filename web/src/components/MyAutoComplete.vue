<template>
  <div>
    <el-autocomplete
      ref="autocomplete"
      v-model="searchKey"
      :fetch-suggestions="querySearch"
      clearable
      :placeholder="props.placeholder"
      @change="handleChange"
      @select="handleChange"
      @keyup.enter.stop="searchInputKey"
    >
      <template #suffix>
        <el-icon class="el-input__icon">
          <ArrowDown />
        </el-icon>
      </template>
    </el-autocomplete>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

const searchKey = ref("");
const autocomplete: any = ref(null);

const props = defineProps({
  all_data: {
    type: Array,
    required: true,
    default: [],
  },
  placeholder: {
    type: String,
    required: true,
    default: "请输入关键字",
  },
});

const emit = defineEmits<{
  (e: "change", params: SearchItem, input_key: string): void;
  (e: "input", input_key: string): void;
}>();

watch(
  () => props.all_data,
  (newV, oldV) => {
    if (newV) {
      all_data.value = newV;
    }
  }
);

interface SearchItem {
  value: string;
  [key: string]: string;
}
const all_data = ref([] as any);
const querySearch = (queryString: string, cb: any) => {
  const results = queryString
    ? all_data.value.filter(createFilter(queryString))
    : all_data.value;
  cb(results);
};
const createFilter = (queryString: string) => {
  return (keyword: SearchItem) => {
    return keyword.value.indexOf(queryString) === 0;
  };
};

// 搜索事件
const handleChange = (value: SearchItem) => {
  emit("change", value, searchKey.value);
  if (!value) {
    $blur();
  }
};
const searchInputKey = () => {
  emit("input", searchKey.value);
};

// 失焦事件
const $blur = () => {
  autocomplete.value.blur();
};

defineExpose({
  searchKey,
  $blur,
});
</script>

<style scoped></style>
