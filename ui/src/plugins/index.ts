import router from '../router';
import { createPinia } from 'pinia';
import piniaPersistedState from 'pinia-plugin-persistedstate'
/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Types
import type { App } from 'vue'

// Plugins
import vuetify from './vuetify'

export function registerPlugins(app: App) {
    const pinia = createPinia()
    pinia.use(piniaPersistedState)
    app.use(vuetify)
    app.use(pinia);
    app.use(router);
}