<template>
    <div class="image__preview">
      <el-image
        style="position:absolute; width:100%; height:100%; top:0; left:0;"
        :src="url"
        :preview-src-list="srcList"
        :initial-index="0"
        fit="contain"
      />
    </div>
  </template>
  
  <script lang="ts" setup>
  const url =
    '//fuss10.elemecdn.com/a/3f/3302e58f9a181d2509f3dc0fa68b0jpeg.jpeg'
  const srcList = [
    '//fuss10.elemecdn.com/a/3f/3302e58f9a181d2509f3dc0fa68b0jpeg.jpeg'
  ]
  </script>
  
  <style scoped>
  .image__error .image-slot {
    font-size: 30px;
  }
  .image__error .image-slot .el-icon {
    font-size: 30px;
  }
  .image__error .el-image {
    width: 100%;
    height: 200px;
  }
  </style>