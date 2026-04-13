<template>
  <div class="system-settings">
    <t-card title="系统设置" :bordered="false">
      <t-tabs v-model="activeTab">
        <t-tab-panel value="basic" label="基础配置">
          <t-form :data="basicForm" :rules="basicRules" ref="basicFormRef" @submit="handleBasicSubmit">
            <t-form-item label="平台名称" name="platformName">
              <t-input v-model="basicForm.platformName" placeholder="请输入平台名称" />
            </t-form-item>
            <t-form-item label="平台Logo">
              <t-upload
                v-model="basicForm.logo"
                theme="image"
                accept="image/*"
                :auto-upload="false"
                tips="建议尺寸：200x60，支持PNG、JPG格式"
              />
            </t-form-item>
            <t-form-item label="平台描述" name="description">
              <t-textarea 
                v-model="basicForm.description" 
                placeholder="请输入平台描述" 
                :maxlength="200"
                :autosize="{ minRows: 3, maxRows: 5 }"
              />
            </t-form-item>
            <t-form-item label="联系邮箱" name="contactEmail">
              <t-input v-model="basicForm.contactEmail" placeholder="请输入联系邮箱" />
            </t-form-item>
            <t-form-item label="联系电话" name="contactPhone">
              <t-input v-model="basicForm.contactPhone" placeholder="请输入联系电话" />
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" type="submit">保存配置</t-button>
            </t-form-item>
          </t-form>
        </t-tab-panel>

        <t-tab-panel value="map" label="地图配置">
          <t-form :data="mapForm" ref="mapFormRef" @submit="handleMapSubmit">
            <t-form-item label="默认中心经度" name="centerLng">
              <t-input-number v-model="mapForm.centerLng" :decimal-places="4" />
            </t-form-item>
            <t-form-item label="默认中心纬度" name="centerLat">
              <t-input-number v-model="mapForm.centerLat" :decimal-places="4" />
            </t-form-item>
            <t-form-item label="默认缩放级别" name="defaultZoom">
              <t-slider v-model="mapForm.defaultZoom" :min="1" :max="10" />
            </t-form-item>
            <t-form-item label="最小缩放级别" name="minZoom">
              <t-slider v-model="mapForm.minZoom" :min="1" :max="10" />
            </t-form-item>
            <t-form-item label="最大缩放级别" name="maxZoom">
              <t-slider v-model="mapForm.maxZoom" :min="1" :max="10" />
            </t-form-item>
            <t-form-item label="地图主题">
              <t-radio-group v-model="mapForm.theme">
                <t-radio value="dark">深色主题</t-radio>
                <t-radio value="light">浅色主题</t-radio>
              </t-radio-group>
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" type="submit">保存配置</t-button>
            </t-form-item>
          </t-form>
        </t-tab-panel>

        <t-tab-panel value="storage" label="存储配置">
          <t-form :data="storageForm" ref="storageFormRef" @submit="handleStorageSubmit">
            <t-form-item label="存储类型">
              <t-radio-group v-model="storageForm.storageType">
                <t-radio value="local">本地存储</t-radio>
                <t-radio value="oss">阿里云OSS</t-radio>
                <t-radio value="cos">腾讯云COS</t-radio>
              </t-radio-group>
            </t-form-item>
            
            <template v-if="storageForm.storageType === 'local'">
              <t-form-item label="存储路径" name="localPath">
                <t-input v-model="storageForm.localPath" placeholder="/data/uploads" />
              </t-form-item>
            </template>
            
            <template v-else-if="storageForm.storageType === 'oss'">
              <t-form-item label="Access Key ID" name="ossAccessKeyId">
                <t-input v-model="storageForm.ossAccessKeyId" placeholder="请输入Access Key ID" />
              </t-form-item>
              <t-form-item label="Access Key Secret" name="ossAccessKeySecret">
                <t-input v-model="storageForm.ossAccessKeySecret" type="password" placeholder="请输入Access Key Secret" />
              </t-form-item>
              <t-form-item label="Bucket名称" name="ossBucket">
                <t-input v-model="storageForm.ossBucket" placeholder="请输入Bucket名称" />
              </t-form-item>
              <t-form-item label="Endpoint" name="ossEndpoint">
                <t-input v-model="storageForm.ossEndpoint" placeholder="oss-cn-hangzhou.aliyuncs.com" />
              </t-form-item>
            </template>
            
            <template v-else-if="storageForm.storageType === 'cos'">
              <t-form-item label="Secret ID" name="cosSecretId">
                <t-input v-model="storageForm.cosSecretId" placeholder="请输入Secret ID" />
              </t-form-item>
              <t-form-item label="Secret Key" name="cosSecretKey">
                <t-input v-model="storageForm.cosSecretKey" type="password" placeholder="请输入Secret Key" />
              </t-form-item>
              <t-form-item label="Bucket名称" name="cosBucket">
                <t-input v-model="storageForm.cosBucket" placeholder="请输入Bucket名称" />
              </t-form-item>
              <t-form-item label="Region" name="cosRegion">
                <t-input v-model="storageForm.cosRegion" placeholder="ap-guangzhou" />
              </t-form-item>
            </template>
            
            <t-form-item label="最大文件大小" name="maxFileSize">
              <t-input-number v-model="storageForm.maxFileSize" :min="1" :max="500" suffix="MB" />
            </t-form-item>
            <t-form-item label="允许的文件类型" name="allowedTypes">
              <t-input v-model="storageForm.allowedTypes" placeholder="jpg,png,gif,mp4" />
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" type="submit">保存配置</t-button>
              <t-button theme="default" @click="testStorage" style="margin-left: 8px;">测试连接</t-button>
            </t-form-item>
          </t-form>
        </t-tab-panel>

        <t-tab-panel value="security" label="安全配置">
          <t-form :data="securityForm" ref="securityFormRef" @submit="handleSecuritySubmit">
            <t-form-item label="启用验证码">
              <t-switch v-model="securityForm.enableCaptcha" />
            </t-form-item>
            <t-form-item label="登录失败锁定次数" name="maxLoginAttempts">
              <t-input-number v-model="securityForm.maxLoginAttempts" :min="3" :max="10" />
            </t-form-item>
            <t-form-item label="锁定时长(分钟)" name="lockDuration">
              <t-input-number v-model="securityForm.lockDuration" :min="5" :max="60" />
            </t-form-item>
            <t-form-item label="密码最小长度" name="minPasswordLength">
              <t-input-number v-model="securityForm.minPasswordLength" :min="6" :max="20" />
            </t-form-item>
            <t-form-item label="密码强度要求">
              <t-checkbox-group v-model="securityForm.passwordRequirements">
                <t-checkbox value="uppercase">包含大写字母</t-checkbox>
                <t-checkbox value="lowercase">包含小写字母</t-checkbox>
                <t-checkbox value="number">包含数字</t-checkbox>
                <t-checkbox value="special">包含特殊字符</t-checkbox>
              </t-checkbox-group>
            </t-form-item>
            <t-form-item label="Token有效期(小时)" name="tokenExpireHours">
              <t-input-number v-model="securityForm.tokenExpireHours" :min="1" :max="168" />
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" type="submit">保存配置</t-button>
            </t-form-item>
          </t-form>
        </t-tab-panel>
      </t-tabs>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'

const activeTab = ref('basic')

const basicFormRef = ref()
const mapFormRef = ref()
const storageFormRef = ref()
const securityFormRef = ref()

const basicForm = reactive({
  platformName: '全景漫游平台',
  logo: [],
  description: '',
  contactEmail: '',
  contactPhone: ''
})

const basicRules = {
  platformName: [{ required: true, message: '请输入平台名称', trigger: 'blur' }],
  contactEmail: [{ email: true, message: '请输入正确的邮箱格式', trigger: 'blur' }]
}

const mapForm = reactive({
  centerLng: 104.5,
  centerLat: 36,
  defaultZoom: 1.2,
  minZoom: 0.8,
  maxZoom: 6,
  theme: 'dark'
})

const storageForm = reactive({
  storageType: 'local',
  localPath: '/data/uploads',
  ossAccessKeyId: '',
  ossAccessKeySecret: '',
  ossBucket: '',
  ossEndpoint: '',
  cosSecretId: '',
  cosSecretKey: '',
  cosBucket: '',
  cosRegion: '',
  maxFileSize: 50,
  allowedTypes: 'jpg,jpeg,png,gif,mp4'
})

const securityForm = reactive({
  enableCaptcha: true,
  maxLoginAttempts: 5,
  lockDuration: 30,
  minPasswordLength: 8,
  passwordRequirements: ['uppercase', 'lowercase', 'number'],
  tokenExpireHours: 24
})

const handleBasicSubmit = async () => {
  const valid = await basicFormRef.value?.validate()
  if (valid !== true) return
  
  localStorage.setItem('systemSettings_basic', JSON.stringify(basicForm))
  MessagePlugin.success('基础配置保存成功')
}

const handleMapSubmit = () => {
  localStorage.setItem('systemSettings_map', JSON.stringify(mapForm))
  MessagePlugin.success('地图配置保存成功')
}

const handleStorageSubmit = () => {
  localStorage.setItem('systemSettings_storage', JSON.stringify(storageForm))
  MessagePlugin.success('存储配置保存成功')
}

const testStorage = () => {
  MessagePlugin.info('正在测试存储连接...')
  setTimeout(() => {
    MessagePlugin.success('存储连接测试成功')
  }, 1000)
}

const handleSecuritySubmit = () => {
  localStorage.setItem('systemSettings_security', JSON.stringify(securityForm))
  MessagePlugin.success('安全配置保存成功')
}
</script>

<style scoped>
.system-settings {
  padding: 24px;
}

:deep(.t-card__header) {
  border-bottom: 1px solid var(--td-component-border);
}

:deep(.t-form) {
  max-width: 600px;
  margin-top: 24px;
}

:deep(.t-form-item) {
  margin-bottom: 24px;
}

:deep(.t-input-number) {
  width: 200px;
}

:deep(.t-slider) {
  width: 300px;
}
</style>
