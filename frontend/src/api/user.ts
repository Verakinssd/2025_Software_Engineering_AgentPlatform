import { user } from '@/stores/user';
import server from './server';
import { ElMessage } from 'element-plus';

export const register = async (username: string , password: string , email: string) => {
    return server.Post(
        "/register",
        {
            username: username,
            password: password,
            email: email
        } ,
        { headers: { 'Content-Type': 'multipart/form-data' } }
    )
}

export const login = async (username: string , password: string) => {
    const res = await server.Post(
        "/login",
        {
            username: username,
            password: password
        },
        { headers: { 'Content-Type': 'multipart/form-data' } }
    )
    if (res.status === 200) {
      user.loggedIn = true;
      user.avatar = res.data.avatar;
      user.email = res.data.email;
      user.id = res.data.id;
      user.role = res.data.role;
      user.username = res.data.username;
      user.profile = res.data.profile;
      ElMessage.success(`欢迎回来，${user.username}`)
    }
    return res;
}

export const logout = async () => {
    const res = await server.Post(
        '/user/logout'
    )
    if (res.status === 200) {
      user.loggedIn = false;
      user.avatar = '';
      user.email = '';
      user.id = '';
      user.role = '';
      user.username = '';
      user.profile = '';
    }
    return res;
}

export const getProfile = async () => {
  const res = await server.Get('/user/profile')
  if (res.status === 200) {
    user.loggedIn = true;
    user.avatar = res.data.avatar;
    user.email = res.data.email;
    user.id = res.data.id;
    user.role = res.data.role;
    user.username = res.data.username;
    user.profile = res.data.profile;
  }
  return res;
}

export const uploadUserAvatar = async (file: File) => {
  const formData = new FormData();
  formData.append('avatar', file);

  const res = await server.Post('/user/upload_avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  });

  if (res.status === 200) {
    user.avatar = res.data.avatar;
  }
  return res;
};

export const updateUserProfile = async (data: {
  new_username?: string;
  profile?: string;
  new_email?: string;
  old_password?: string;
  new_password?: string;
}) => {
  const formData = new FormData();
  if (data.new_username) {
    formData.append('new_username', data.new_username);
  }
  if (data.profile) {
    formData.append('profile', data.profile);
  }
  if (data.new_email) {
    formData.append('new_email', data.new_email);
  }
  if (data.old_password) {
    formData.append('old_password', data.old_password);
  }
  if (data.new_password) {
    formData.append('new_password', data.new_password);
  }
  
  const res = await server.Post(
    "/user/update_profile",
    formData,
    { headers: { 'Content-Type': 'multipart/form-data' } }
  )
  if (res.status === 200) {
    user.avatar = res.data.avatar;
    user.email = res.data.email;
    user.id = res.data.id;
    user.role = res.data.role;
    user.username = res.data.username;
    user.profile = res.data.profile;
  }
};
