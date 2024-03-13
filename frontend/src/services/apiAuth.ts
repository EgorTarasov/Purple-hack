import axios from "axios";
import { BASE_URL } from "../config";
import uuid from "react-uuid";
import storage from '../lib/storage'

interface UserLogin {
	email: string;
	password: string;
}

axios.defaults.withCredentials = true;

const ApiAuth = {
	async loginUser(data: UserLogin) {
		const response = await axios.post(`${BASE_URL}/auth/login`, data);
		if(response) storage.setToken(uuid());
		return response;
	},
};
export default ApiAuth;
