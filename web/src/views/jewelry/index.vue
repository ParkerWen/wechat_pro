<template>
  <div class="app-container">
    <el-form ref="elFormRef" :model="formData">
      <div>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="24" :md="10" :lg="10" :xl="10">
            <div class="block" style="height: 600px">
              <el-image v-loading="loading" :src="imgUrl" />
              <el-row :gutter="20" style="margin-top: 6%;" v-if="btnGroup">
                <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                  <el-button @click="upscale(1)" :loading="loading">选择第一张</el-button>
                </el-col>
                <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                  <el-button @click="upscale(2)" :loading="loading">选择第二张</el-button>
                </el-col>
                <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                  <el-button @click="upscale(3)" :loading="loading">选择第三张</el-button>
                </el-col>
                <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                  <el-button @click="upscale(4)" :loading="loading">选择第四张</el-button>
                </el-col>
              </el-row>
              <el-row :gutter="20" style="margin-top: 6%;" v-if="btnGroup">
                <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                  <el-button @click="variation(1)" :loading="loading">更换第一张</el-button>
                </el-col>
                <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                  <el-button @click="variation(2)" :loading="loading">更换第二张</el-button>
                </el-col>
                <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                  <el-button @click="variation(3)" :loading="loading">更换第三张</el-button>
                </el-col>
                <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                  <el-button @click="variation(4)" :loading="loading">更换第四张</el-button>
                </el-col>
              </el-row>
            </div>
          </el-col>

          <el-col :xs="24" :sm="24" :md="14" :lg="14" :xl="14">

            <div style="height: 600px">
              <div class="scrollbar-item">
                <el-tabs v-model="activeName" @tab-click="handleTab">

                  <el-tab-pane label="设计特征" name="feature">
                    <el-row :gutter="10">
                      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                        <el-select v-model="diamond" placeholder="钻石">
                          <el-option v-for="item in diamondOptions" :key="item.value" :label="item.label"
                            :value="item.value" />
                        </el-select>
                      </el-col>
                      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                        <el-select v-model="precious" placeholder="贵宝">
                          <el-option v-for="item in preciousOptions" :key="item.value" :label="item.label"
                            :value="item.value" />
                        </el-select>
                      </el-col>
                      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                        <el-select v-model="colored" placeholder="彩宝">
                          <el-option v-for="item in coloredOptions" :key="item.value" :label="item.label"
                            :value="item.value" />
                        </el-select>
                      </el-col>
                      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                        <el-select v-model="faceted" placeholder="刻面宝石">
                          <el-option v-for="item in facetedOptions" :key="item.value" :label="item.label"
                            :value="item.value" />
                        </el-select>
                      </el-col>
                      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                        <el-select v-model="cabochon" placeholder="蛋面宝石">
                          <el-option v-for="item in cabochonOptions" :key="item.value" :label="item.label"
                            :value="item.value" />
                        </el-select>
                      </el-col>
                      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
                        <el-select v-model="inlay" placeholder="镶嵌">
                          <el-option v-for="item in inlayOptions" :key="item.value" :label="item.label"
                            :value="item.value" />
                        </el-select>
                      </el-col>
                    </el-row>
                  </el-tab-pane>

                  <el-tab-pane label="珠宝图集" name="image">

                    <el-row :gutter="20">
                      <el-scrollbar style="height: 500px;width: 100%;">
                        <el-col v-for="item in imageList" :key="item.id" :xs="24" :sm="24" :md="6" :lg="6" :xl="6"
                          class="scrollbar-item">
                          <el-image v-loading="loading" :src="item.image_url" :preview-src-list="[item.image_url]" />
                        </el-col>
                      </el-scrollbar>

                    </el-row>

                  </el-tab-pane>

                </el-tabs>
              </div>

            </div>

          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="24"><el-button :loading="loading" type="primary" size="large" style="width: 100%;"
              @click="handleSave">生成图片</el-button></el-col>
        </el-row>
      </div>
    </el-form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      formData: {
        prompt: 'jeweler,'
      },
      imageList: [],
      orignImg: '',
      upscaleReq: {
        task_id: '',
        index: 1
      },
      upscaleImg: [],
      variationReq: {
        task_id: '',
        index: 1
      },
      btnGroup: false,
      loading: false,
      redirect: undefined,
      activeName: 'feature',
      diamond: 'Blue diamond',
      diamondOptions: [
        { 'value': 'White diamond', 'label': '白钻' },
        { 'value': 'Yellow diamond', 'label': '黄钻' },
        { 'value': 'Pink diamond', 'label': '粉钻' },
        { 'value': 'Green diamond', 'label': '绿钻' },
        { 'value': 'Black diamond', 'label': '黑钻' },
        { 'value': 'Purple diamond', 'label': '紫钻' },
        { 'value': 'Red diamond', 'label': '红钻 ' },
        { 'value': 'Blue diamond', 'label': '蓝钻 ' }
      ],
      precious: 'sapphire',
      preciousOptions: [
        { 'value': 'ruby', 'label': '红宝石' },
        { 'value': 'sapphire', 'label': '蓝宝石' },
        { 'value': 'emerald', 'label': '祖母绿' },
        { 'value': 'jade', 'label': '翡翠' },
        { 'value': "cat's eye", 'label': '猫眼石' },
        { 'value': 'paparacha', 'label': '帕帕拉恰' },
        { 'value': 'paraiba', 'label': '帕拉伊巴' },
        { 'value': 'opal', 'label': '欧珀' },
        { 'value': 'Mother of Pearl', 'label': '海螺珠' }
      ],
      colored: 'Garnet',
      coloredOptions: [
        { 'value': 'Garnet', 'label': '石榴石' },
        { 'value': 'Peridot', 'label': '橄榄石' },
        { 'value': 'Coral', 'label': '珊瑚' },
        { 'value': 'Pearl', 'label': '珍珠' },
        { 'value': 'Amber', 'label': '琥珀' },
        { 'value': 'Tanzanite', 'label': '坦桑石' },
        { 'value': 'Morganite', 'label': '摩根石' },
        { 'value': 'Lepidolite', 'label': '紫锂辉石' },
        { 'value': 'Hetian Jade', 'label': '和田玉' },
        { 'value': 'Spinel', 'label': '尖晶石' },
        { 'value': 'Hibonite', 'label': '芙蓉石' },
        { 'value': 'Aquamarine', 'label': '碧玺' },
        { 'value': 'Sapphire', 'label': '莎弗莱' }
      ],
      faceted: 'Heart',
      facetedOptions: [
        { 'value': 'Round', 'label': '圆形' },
        { 'value': 'Oval', 'label': '椭圆形' },
        { 'value': 'Marquise', 'label': '马眼形' },
        { 'value': 'Pear', 'label': '梨形' },
        { 'value': 'Heart', 'label': '心形' },
        { 'value': 'Emerald', 'label': '祖母绿形 ' }
      ],
      cabochon: 'Translucent',
      cabochonOptions: [
        { 'value': 'Translucent', 'label': '半透明' },
        { 'value': 'Opaque', 'label': '不透明' },
        { 'value': 'Crushed', 'label': '砂面' },
        { 'value': 'Mirror', 'label': '镜面' }
      ],
      inlay: 'Claw Setting',
      inlayOptions: [
        { 'value': 'Prong Setting', 'label': '爪镶' },
        { 'value': 'Flush Setting', 'label': '包镶' },
        { 'value': 'Claw Setting', 'label': '钉镶' },
        { 'value': 'Pave Setting', 'label': '密钉镶' },
        { 'value': 'Channel Setting', 'label': '卡镶' },
        { 'value': 'Bezel Setting', 'label': '无边镶 ' }
      ],
      imgUrl: 'http://img.itcity.cc/attachments/1109090103432855724/1113652867987935302/wsy_jewelerBlue_diamondsapphirePeridotHeartTranslucentChannel_S_1aa5e53b-5d64-45a2-9127-d3641911a37b.png',
    }
  },
  methods: {
    getImageList() {
      this.$store.dispatch('jewelry/getImageList').then((list) => {
        this.imageList = list
      }).catch(() => {
        this.loading = false
      })
    },
    handleSave() {
      const that = this
      this.loading = true
      this.formData.prompt += [this.diamond, this.precious, this.colored, this.faceted, this.cabochon, this.inlay].join(',')
      this.$store.dispatch('jewelry/create', this.formData).then((task) => {
        const timer = setInterval(function () {
          that.$store.dispatch('jewelry/getInfo', task.id).then((res) => {
            if (res.status === 'SUCCESS' && res.image_url.length > 0) {
              that.upscaleReq.task_id = res.task_id
              that.variationReq.task_id = res.task_id
              that.btnGroup = true
              that.imgUrl = res.image_url
              that.orignImg = res.image_url
              that.loading = false
              clearInterval(timer)
            }
          }).catch(() => {
            clearInterval(timer)
            that.loading = false
          })
        }, 8000)
      }).catch(() => {
        this.loading = false
      })
    },
    upscale(index) {
      var i = this.upscaleImg.findIndex(item => item.index === index)
      if (i !== -1) {
        this.imgUrl = this.upscaleImg[i].imgUrl
        return
      }
      const that = this
      this.loading = true
      this.upscaleReq.index = index
      this.$store.dispatch('jewelry/upscale', this.upscaleReq).then((task) => {
        const timer = setInterval(function () {
          that.$store.dispatch('jewelry/getInfo', task.id).then((res) => {
            if (res.status === 'SUCCESS' && res.image_url.length > 0) {
              const upscaleItem = {
                imgUrl: res.image_url,
                index: index,
              }
              that.upscaleImg.push(upscaleItem)
              that.imgUrl = res.image_url
              that.loading = false
              clearInterval(timer)
            }
          }).catch(() => {
            clearInterval(timer)
            that.loading = false
          })
        }, 8000)
      }).catch(() => {
        this.loading = false
      })
    },
    variation(index) {
      this.upscaleImg = []
      const that = this
      this.loading = true
      this.variationReq.index = index
      this.$store.dispatch('jewelry/variation', this.variationReq).then((task) => {
        that.upscaleReq.task_id = task.task_id
        that.variationReq.task_id = task.task_id
        const timer = setInterval(function () {
          that.$store.dispatch('jewelry/getInfo', task.id).then((res) => {
            if (res.status === 'SUCCESS' && res.image_url.length > 0) {
              that.imgUrl = res.image_url
              that.loading = false
              clearInterval(timer)
            }
          }).catch(() => {
            clearInterval(timer)
            that.loading = false
          })
        }, 8000)
      }).catch(() => {
        this.loading = false
      })
    },
    handleTab(tab, event) {
      if (this.activeName === "image") {
        this.getImageList()
      }
    },
  }
}
</script>

<style>
.el-container {
  background-color: #d3dce6;
}

.el-aside {
  background-color: #d3dce6;
  color: #333;
  text-align: center;
}

.el-main {
  background-color: #e9eef3;
  color: #333;
  text-align: center;
}

.el-image {
  width: 70%;
}

.el-card {
  margin-bottom: 5%;
}

.el-card-active {
  border: "5px solid #4d70ff";
}

.el-select {
  margin-bottom: 5%;
}

.row-bg {
  padding: 10px 0;
  background-color: #f9fafc;
}

.block {
  text-align: center;
}

.image {
  width: 100%;
  display: block;
}

.scrollbar-item {
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  border-radius: 4px;
  padding: 1%;
}
</style>
