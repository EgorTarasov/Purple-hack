import axios from "axios";
import { BASE_URL } from "../config";

interface UserLogin {
    email: string;
    password: string;
}


const ApiAuth = {
    async loginUser(data: UserLogin) {
        const response = await axios.post(
            `${BASE_URL}/auth/login`,
            data
        );
        
        // const response = await axios.post(`${BASE_URL}/auth/login`, data);
        return response;
    },
};
export default ApiAuth;