import { reactive } from 'vue'

interface IUser {
  loggedIn: boolean
  id: string
  email: string
  role: string
  username: string
  avatar: string
  profile: string
}

export const user = reactive<IUser>({
  loggedIn: false,
  id: '',
  email: '',
  role: '',
  username: '',
  avatar: '',
  profile: ''
})