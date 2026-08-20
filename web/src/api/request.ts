import axios from 'axios'
import { useUserStore } from '../stores/user'
import router from '../router'

const request = axios.create({ baseURL: '/api/v1', timeout: 15000 })

// 请求拦截：附带 token
request.interceptors.request.use((config) => {
  const userStore = useUserStore()
  if (userStore.accessToken) {
    config.headers.Authorization = `Bearer ${userStore.accessToken}`
  }
  return config
})

// 响应拦截：401 时尝试刷新 token 后重放请求
let refreshing: Promise<string | null> | null = null

request.interceptors.response.use(
  (resp) => resp,
  async (error) => {
    const { response, config } = error
    if (response?.status === 401 && !config._retried) {
      config._retried = true
      refreshing ||= useUserStore().refresh().finally(() => (refreshing = null))
      const token = await refreshing
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
        return request(config)
      }
      router.push('/login')
    }
    return Promise.reject(error)
  },
)

export default request
