import axios from "axios";
import { BASE_URL } from "../config";
import { ISession } from "@/models";


const ApiSession = {
    async getUserSession() {
        const response = await axios.get<ISession>(
            `${BASE_URL}/api/sessions/list`
        );
        
        // const response = await axios.post(`${BASE_URL}/auth/login`, data);
        return response;
    },
};
export default ApiSession;