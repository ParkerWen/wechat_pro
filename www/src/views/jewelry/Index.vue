<template>
  <div v-masonry="containerId" transition-duration="0.3s" item-selector=".item" >
    <a-button @click="onShow" style="">开始设计珠宝</a-button>
    <div class="ant-pro-pages-list-projects-cardList" >
      <a-list :loading="loading" :data-source="data" :grid="{ gutter: 24, xl: 4, lg: 3, md: 3, sm: 2, xs: 1 }" class="card-list">
        <a-list-item slot="renderItem" slot-scope="item" :loading="true" v-masonry-tile class="item">
          <template>
            <a-card class="ant-pro-pages-list-projects-card" hoverable>
              <img slot="cover" :src="item.image_url?item.image_url:loadingImg" :alt="item.prompt" style="width: 100%; height: 300px" @click="previewImage(item.image_url)" />
              <template slot="actions" v-if="item.action === 'UPSCALE'">
                <a-tooltip title="下载">
                  <a-icon type="download" @click="downloadHandle(item.image_url, item.task_id)" />
                </a-tooltip>
                <a-tooltip title="分享">
                  <a-icon type="share-alt" @click="copyData(item.image_url)" />
                </a-tooltip>
              </template>
              <template slot="actions" v-else>
                <a-dropdown>
                  <a class="ant-dropdown-link">
                    选择图片
                  </a>
                  <a-menu slot="overlay" style="text-align: center;">
                    <a-menu-item>
                      <a href="javascript:;" @click="upscaleHandle(1, item.task_id)">第一张</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a href="javascript:;" @click="upscaleHandle(2, item.task_id)">第二张</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a href="javascript:;" @click="upscaleHandle(3, item.task_id)">第三张</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a href="javascript:;" @click="upscaleHandle(4, item.task_id)">第四张</a>
                    </a-menu-item>
                  </a-menu>
                </a-dropdown>
                <a-dropdown>
                  <a class="ant-dropdown-link">
                    轮换图片
                  </a>
                  <a-menu slot="overlay" style="text-align: center;">
                    <a-menu-item>
                      <a href="javascript:;" @click="variationHandle(1, item.task_id)">第一张</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a href="javascript:;" @click="variationHandle(1, item.task_id)">第二张</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a href="javascript:;" @click="variationHandle(1, item.task_id)">第三张</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a href="javascript:;" @click="variationHandle(1, item.task_id)">第四张</a>
                    </a-menu-item>
                  </a-menu>
                </a-dropdown>
              </template>
            </a-card>
          </template>
        </a-list-item>
      </a-list>
    </div>

    <a-drawer
      title="AI 设计"
      width="500"
      :placement="placement"
      :closable="false"
      :visible="visible"
      @close="onClose"
    >
      <div class="setting-drawer-index-content">
        <a-form>
          <div :style="{ marginBottom: '25px' }">
            <h3 class="setting-drawer-index-title">模型</h3>

            <div class="setting-drawer-index-blockChecbox">
              <a-tooltip>
                <template slot="title">
                  侧边栏导航
                </template>
                <div class="setting-drawer-index-item">
                  <img src="/default.png" alt="sidemenu">
                  <div class="setting-drawer-index-selectIcon">
                    <a-icon type="check"/>
                  </div>
                </div>
              </a-tooltip>
            </div>

            <div :style="{ marginTop: '25px' }">
              <h3 class="setting-drawer-index-title">样式</h3>
              <a-list :split="false">
                <a-list-item>
                  <a-select size="small" v-model="diamond" :options="diamondOptions" style="width: 100px"></a-select>
                  <a-select size="small" v-model="precious" :options="preciousOptions" style="width: 100px"></a-select>
                  <a-select size="small" v-model="colored" :options="coloredOptions" style="width: 100px"></a-select>
                  <a-select size="small" v-model="faceted" :options="facetedOptions" style="width: 100px"></a-select>
                </a-list-item>
                <a-list-item>
                  <a-select size="small" v-model="cabochon" :options="cabochonOptions" style="width: 100px"></a-select>
                  <a-select size="small" v-model="inlay" :options="inlayOptions" style="width: 100px"></a-select>
                  <a-select size="small" v-model="category" :options="categoryOptions" style="width: 100px"></a-select>
                  <a-select size="small" v-model="jewelry" :options="jewelryOptions" style="width: 100px"></a-select>

                </a-list-item>
              </a-list>
            </div>
          </div>
        </a-form>
        <div :style="{ marginBottom: '25px' }">
          <a-button
            icon="sketch"
            block
            @click="imagineHandle"
          >立即生成</a-button>
        </div>
      </div>
      <div class="setting-drawer-index-handle" @click="toggle" slot="handle">
        <a-icon type="left" v-if="!visible"/>
        <a-icon type="close" v-else/>
      </div>
    </a-drawer>
  </div>
</template>

<script>
import moment from 'moment'
import { TagSelect, StandardFormRow, Ellipsis, AvatarList } from '@/components'
import { list, imagine, upscale, variation } from '@/api/jewelry'
import Vue from 'vue'
import VueClipboard from 'vue-clipboard2'

Vue.use(VueClipboard)
const TagSelectOption = TagSelect.Option
const AvatarListItem = AvatarList.Item

export default {
  components: {
    AvatarList,
    AvatarListItem,
    Ellipsis,
    TagSelect,
    TagSelectOption,
    StandardFormRow
  },
  data () {
    return {
      data: [],
      form: this.$form.createForm(this),
      loading: true,
      placement: 'right',
      visible: false,
      top: 180,
      layoutMode: '',
      loadingImg: '/image-loader.gif',
      containerId: '.card-list',
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
      category: 'Gold Jewelry',
      categoryOptions: [
        { 'value': 'Gold Jewelry', 'label': '金饰' },
        { 'value': 'Silver Jewelry', 'label': '银饰' },
        { 'value': 'Light Luxury', 'label': '轻奢' },
        { 'value': 'Classic', 'label': '经典' },
        { 'value': 'Enamel', 'label': '珐琅' },
        { 'value': 'Middle', 'label': '中端' }
      ],
      jewelry: 'Ring',
      jewelryOptions: [
        { 'value': 'Earrings', 'label': '耳钉' },
        { 'value': 'Pendant', 'label': '吊坠' },
        { 'value': 'Ring', 'label': '戒指' },
        { 'value': 'Necklace', 'label': '项链' },
        { 'value': 'Bracelet', 'label': '手镯' }
      ]
    }
  },
  filters: {
    fromNow (date) {
      return moment(date).fromNow()
    }
  },
  mounted () {
    this.getList()
    setInterval(this.getList, 30000)
  },
  methods: {
    onClose () {
      this.visible = false
    },
    onShow () {
      this.visible = true
    },
    previewImage (img) {
      if (img) {
        window.open(img)
      }
    },
    copyData (image) {
      this.$copyText(image).then(e => {
          alert('复制图片链接成功')
      }, function (e) {
          alert('复制图片链接失败')
          console.log(e)
      })
    },
    downloadHandle (image, name) {
      this.downloadByBlob(image, name)
    },
    downloadByBlob (url, name) {
      const image = new Image()
      image.setAttribute('crossOrigin', 'anonymous')
      image.src = url
      image.onload = () => {
        const canvas = document.createElement('canvas')
        canvas.width = image.width
        canvas.height = image.height
        const ctx = canvas.getContext('2d')
        ctx.drawImage(image, 0, 0, image.width, image.height)
        canvas.toBlob((blob) => {
          const url = URL.createObjectURL(blob)
          this.download(url, name)
          URL.revokeObjectURL(url)
        })
      }
    },
    download (href, name) {
      const eleLink = document.createElement('a')
      eleLink.download = name
      eleLink.href = href
      eleLink.click()
      eleLink.remove()
    },
    getList () {
      list().then(res => {
        if (res.data) {
          this.data = res.data
        }
        this.loading = false
      }).catch(error => {
        console.log('list fail:', error)
      })
    },
    imagineHandle () {
      this.loading = true
      const imagineParam = {
        'prompt': ''
      }
      imagineParam.prompt += [this.diamond, this.precious, this.colored, this.faceted, this.cabochon, this.inlay, this.category, this.jewelry].join(',')
      imagine(imagineParam).then((res) => {
        if (res.code === 200) {
        }
        this.visible = false
        this.getList()
        this.loading = false
      })
    },
    upscaleHandle (index, taskId) {
      this.loading = true
      const upscaleParam = {
        'index': index,
        'task_id': taskId
      }
      upscale(upscaleParam).then((res) => {
        if (res.code === 200) {
          const data = res.data
          if (data.image_url) {
            window.open(data.image_url)
          } else {
            this.getList()
            this.$redrawVueMasonry()
          }
          this.loading = false
        }
      })
    },
    variationHandle (index, taskId) {
      this.loading = true
      const variationParam = {
        'index': index,
        'task_id': taskId
      }
      variation(variationParam).then((res) => {
        if (res.code === 200) {
          this.getList()
          this.loading = false
        }
      })
    },
    toggle () {
      this.visible = !this.visible
    }
  }
}
</script>

<style lang="less" scoped>
@import "~@/components/index.less";
.ant-pro-components-tag-select {
  :deep(.ant-pro-tag-select .ant-tag) {
    margin-right: 24px;
    padding: 0 8px;
    font-size: 14px;
  }
}
.ant-pro-pages-list-projects-cardList {
  margin-top: 24px;

  :deep(.ant-card-cover) {
    margin-bottom: 4px;
  }

  :deep(.ant-card-actions) {
    height: 44px;
    overflow: hidden;
    line-height: 22px;
  }

  .cardItemContent {
    display: flex;
    height: 20px;
    margin-top: 16px;
    margin-bottom: -4px;
    line-height: 20px;

    > span {
      flex: 1 1;
      color: rgba(0,0,0,.45);
      font-size: 12px;
    }

    :deep(.ant-pro-avatar-list) {
      flex: 0 1 auto;
    }
  }
}

.ant-card-actions {
    background: #f7f9fa;

    li {
      float: left;
      text-align: center;
      margin: 12px 0;
      color: rgba(0, 0, 0, 0.45);
      width: 50%;

      &:not(:last-child) {
        border-right: 1px solid #e8e8e8;
      }

      a {
        color: rgba(0, 0, 0, .45);
        line-height: 22px;
        display: inline-block;
        width: 100%;
        &:hover {
          color: @primary-color;
        }
      }
    }
  }

  .new-btn {
    background-color: #fff;
    border-radius: 2px;
    width: 100%;
    height: 188px;
  }

  .setting-drawer-index-content {

.setting-drawer-index-blockChecbox {
  display: flex;

  .setting-drawer-index-item {
    margin-right: 16px;
    position: relative;
    border-radius: 4px;
    cursor: pointer;

    img {
      width: 60px;
    }

    .setting-drawer-index-selectIcon {
      position: absolute;
      top: 0;
      right: 0;
      width: 100%;
      padding-top: 15px;
      padding-left: 24px;
      height: 100%;
      color: #1890ff;
      font-size: 14px;
      font-weight: 700;
    }
  }
}
.setting-drawer-theme-color-colorBlock {
  width: 20px;
  height: 20px;
  border-radius: 2px;
  float: left;
  cursor: pointer;
  margin-right: 8px;
  padding-left: 0px;
  padding-right: 0px;
  text-align: center;
  color: #fff;
  font-weight: 700;

  i {
    font-size: 14px;
  }
}
}

.setting-drawer-index-handle {
position: absolute;
top: 240px;
background: #1890ff;
width: 48px;
height: 48px;
right: 500px;
display: flex;
justify-content: center;
align-items: center;
cursor: pointer;
pointer-events: auto;
z-index: 1001;
text-align: center;
font-size: 16px;
border-radius: 4px 0 0 4px;

i {
  color: rgb(255, 255, 255);
  font-size: 20px;
}
}

.masonry-container {
  text-align: center;
}
</style>
