import axios, { AxiosError, AxiosResponse } from "axios";
import { ElMessage } from "element-plus";

const server = axios.create({
    baseURL: '/api',
	withCredentials: true,
});

const handleError = (error: AxiosError) => {
	const { response } = error;
	const { status, data } = response as { status: number; data: any };
	// if (data.code && data.message) {
	// 	ElMessage.error(`${data.code}: ${data.message}`)
	// } else {
	// 	ElMessage.error(`Error: ${status}`)
	// }

	if (response) {
		if (data.code && data.message) {
			//ElMessage.error(`${data.code}: ${data.message}`);
		} else {
			ElMessage.error(`Error: ${status}`);
      console.log(status)
		}
	} else {
		ElMessage.error('服务器无响应');
		console.log('Unhandled error:', error.message);
	}
}

const Get = async (url: string, config?: any): Promise<AxiosResponse<any, any>> => {
	try {
		const res = await server.get(url, config)
		return res
	} catch (error) {
		handleError(error);
		return error;
	}
};

const Post = async (url: string, data?: any, config?: any): Promise<AxiosResponse<any, any>> => {
	try {
		const res = await server.post(url, data, config)
		return res
	} catch (error) {
		handleError(error);
		return error;
	}
};

const Put = async (url: string, data?: any, config?: any): Promise<AxiosResponse<any, any>> => {
	try {
		const res = await server.put(url, data, config)
		return res
	} catch (error) {
		handleError(error);
		return error;
	}
}

const Delete = async (url: string, config?: any): Promise<AxiosResponse<any, any>> => {
	try {
		const res = await server.delete(url, config)
		return res
	} catch (error) {
		handleError(error);
		return error;
	}
}

const Patch = async (url: string, data?: any, config?: any): Promise<AxiosResponse<any, any>> => {
	try {
	  const res = await server.patch(url, data, config);
	  return res;
	} catch (error) {
	  handleError(error);
	  return error;
	}
  };

export default { Get, Post, Put, Delete, Patch }
