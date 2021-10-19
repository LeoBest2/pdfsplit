<template>
  <div id="app">
    <el-container style="align-items: center;">
      <el-header>
        <h1 style="color: #409EFF;">PDF切割工具</h1>
      </el-header>
      <el-main v-loading="loading">
        <el-upload
          class="upload-demo"
          ref="upload"
          action="#"
          :file-list="fileList"
          accept=".pdf"
          :auto-upload="false"
          :limit="1"
          :http-request="uploadFile"
          :on-error="errorHandler"
          :on-success="successHandler"
        >
          <el-button slot="trigger" size="small" type="primary">选取文件</el-button>
          <el-button style="margin-left: 10px;" size="small" type="warning" @click="clearFiles">清除文件</el-button>
          <el-input-number v-model="num" :min="1" :max="1000" label="每隔几页分割"></el-input-number>
          <el-button
            style="margin-left: 10px;"
            size="small"
            type="success"
            @click="submitUpload"
          >点击切割</el-button>
          <div slot="tip" class="el-upload__tip">请上传pdf后缀文件，仅支持单个文件, 默认每10页切割</div>
        </el-upload>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data: function () {
    return {
      fileList: [],
      num: 10,
      loading: false
    }
  },
  methods: {
    submitUpload() {
      if (this.$refs.upload.uploadFiles.length === 0) {
        this.$message.warning('请先添加文件！')
        return
      }
      this.$refs.upload.submit();
    },
    uploadFile(param) {
      this.loading = true;
      var form = new FormData();
      form.append('file', param.file);
      axios({
        method: 'post',
        url: '/upload?span=' + this.num,
        // url: 'http://localhost:8090/upload?span=' + this.num,
        data: form
      }).then(
        res => { console.log(res); this.$message.success(res.data); this.loading = false; }
      ).catch(error => { this.$message.error(error.response.data); this.loading = false; });

    },
    clearFiles(file, fileList) {
      this.$refs.upload.clearFiles();
    },
    errorHandler(err, file, fileList) {
      console.log(err, file, fileList)
      this.$message.error('上传文件' + file.name + '失败' + err)
    },
    successHandler(response, file, fileList) {
      this.$message.success('上传文件' + file + '成功' + response)
    }
  }
}
</script>

<style>
html,
body {
  height: 100%;
  margin: 0;
}
#app {
  height: 80%;
  display: flex;
  align-items: center;
}
</style>
