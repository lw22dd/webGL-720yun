import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import 'tdesign-vue-next/dist/tdesign.css'
import { MessagePlugin } from 'tdesign-vue-next'
import { DialogPlugin } from 'tdesign-vue-next'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import router from './router/router'

const app = createApp(App)
const pinia = createPinia()

pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)
app.use(MessagePlugin)
app.use(DialogPlugin)

app.mount('#app')
