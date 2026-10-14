<template>
  <el-card>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="字典类型" name="type">
        <DictType />
      </el-tab-pane>
      <el-tab-pane label="字典数据" name="data">
        <DictData />
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import DictType from './components/DictType.vue'
import DictData from './components/DictData.vue'

const activeTab = ref('type')
</script>
