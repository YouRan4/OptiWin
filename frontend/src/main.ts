import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import App from './App.vue'
import router from './router'
import zh from './locales/zh.json'
import en from './locales/en.json'

const saved = localStorage.getItem('optiwin_lang')
const systemLang = (navigator.language || '').startsWith('zh') ? 'zh' : 'en'

const i18n = createI18n({
  locale: saved || systemLang,
  messages: { zh, en },
})

createApp(App).use(router).use(i18n).mount('#app')
