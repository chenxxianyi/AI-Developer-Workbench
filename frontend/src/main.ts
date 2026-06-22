import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import { i18n } from './i18n'

// Import global styles (includes Tailwind)
import './styles/globals.css'

const app = createApp(App)

// Use Pinia for state management
app.use(createPinia())

// Use Vue Router
app.use(router)

// Use vue-i18n for localized UI copy
app.use(i18n)

app.mount('#app')
