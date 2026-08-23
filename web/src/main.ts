import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { i18n } from './locales'
import { permissionDirective } from './directives/permission'

const app = createApp(App)
app.use(createPinia())
app.use(i18n)
app.use(router)
app.directive('permission', permissionDirective)
app.mount('#app')
