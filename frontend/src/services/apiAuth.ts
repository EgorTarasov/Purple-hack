import axios from "axios";
import { BASE_URL } from "../config";

interface UserLogin {
    username: string;
    password: string;
}


const ApiAuth = {
    async loginUser(data: UserLogin) {
        const formData = new FormData();
        formData.append("username", data.username);
        formData.append("password", data.password);

        const response = await axios.post(
            `${BASE_URL}/auth/login`,
            formData,
        );
        // const response = await axios.post(`${BASE_URL}/auth/login`, data);
        return response;
    },
};
export default ApiAuth;