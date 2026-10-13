
/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import { createRouter, createWebHistory } from 'vue-router/auto'
import { routes } from 'vue-router/auto-routes'
import {ElMessage} from "element-plus";
import {user} from "@/stores/user";
import {getProfile} from "@/api/user";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Workaround for https://github.com/vitejs/vite/issues/11804
router.onError((err, to) => {
  if (err?.message?.includes?.('Failed to fetch dynamically imported module')) {
    if (!localStorage.getItem('vuetify:dynamic-reload')) {
      console.log('Reloading page to fix dynamic import error')
      localStorage.setItem('vuetify:dynamic-reload', 'true')
      location.assign(to.fullPath)
    } else {
      console.error('Dynamic import error, reloading page did not fix it', err)
    }
  } else {
    console.error(err)
  }
})

router.isReady().then(() => {
  localStorage.removeItem('vuetify:dynamic-reload')
})

// jump to login if user.loggedIn==false
router.beforeEach(async (to, from, next) => {
  await new Promise(resolve => setTimeout(resolve, 100))

  const whiteList = ['/login', '/register']
  if (whiteList.includes(to.path)) {
    return next()
  }

  if (user.loggedIn === true) {
    return next()
  }

  try {
    await getProfile()
    await new Promise(resolve => setTimeout(resolve, 100))
    if (user.loggedIn === true) {
      return next()
    } else {
      return next('/login')
    }
  } catch (error) {
    localStorage.removeItem('token')
    next('/login')
  }
})

export default router
