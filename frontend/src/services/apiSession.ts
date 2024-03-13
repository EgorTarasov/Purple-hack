import axios from "axios";
import { BASE_URL } from "../config";
import { ISession } from "@/models";


const ApiSession = {
    async getUserSession() {
        const response = await axios.get<ISession[]>(
            `${BASE_URL}/api/sessions/list`
        );
        return response;
    },
    async getHistotySharedSession(uuid: string) {
        const response = await axios.get<ISession[]>(
            `${BASE_URL}/api/sessions/list/${uuid}`,
        );
        return response;
    },
};
export default ApiSession;

// const sessionsExample: ISession[] = [
//     {
//         id: "session1",
//         queries: [
//             {
//                 id: 1,
//                 model: "model1",
//                 body: "query1",
//                 createdAt: "2024-03-13T08:00:00Z"
//             },
//             {
//                 id: 2,
//                 model: "model2",
//                 body: "query2",
//                 createdAt: "2024-03-13T08:10:00Z"
//             }
//         ],
//         responses: [
//             {
//                 id: 3,
//                 context: {"key1": ["value1", "value2"]},
//                 body: "response1",
//                 createdAt: "2024-03-13T08:05:00Z"
//             },
//             {
//                 id: 4,
//                 context: {"key2": ["value3", "value4"]},
//                 body: "response2",
//                 createdAt: "2024-03-13T08:15:00Z"
//             }
//         ],
//         createdAt: "2024-03-13T07:55:00Z"
//     },
//     {
//         id: "session2",
//         queries: [
//             {
//                 id: 5,
//                 model: "model3",
//                 body: "query3",
//                 createdAt: "2024-03-13T08:20:00Z"
//             },
//             {
//                 id: 79,
//                 model: "model4",
//                 body: "query4",
//                 createdAt: "2024-03-13T08:30:00Z"
//             },
//             {
//                 id: 56,
//                 model: "model4",
//                 body: "query5",
//                 createdAt: "2024-03-13T08:30:00Z"
//             },
//         ],
//         responses: [
//             {
//                 id: 7,
//                 context: {"key3": ["value5", "value6"]},
//                 body: "response3",
//                 createdAt: "2024-03-13T08:25:00Z"
//             },
//             {
//                 id: 98,
//                 context: {"key3": ["value5", "value6"]},
//                 body: "response5",
//                 createdAt: "2024-03-13T08:25:00Z"
//             },
//         ],
//         createdAt: "2024-03-13T08:15:00Z"
//     }
// ];


// // Модуль для работы с API
// import { ISession } from "@/models";

// export interface IApiSession {
//   getUserSession: () => Promise<ISession[]>;
// }

// const ApiSession: IApiSession = {
//   async getUserSession() {
//     // Вместо реального вызова axios используем тестовый объект
//     return new Promise<ISession[]>((resolve) => {
//       // Здесь должна быть логика для обработки запроса, замокаем ее
//       resolve(sessionsExample); // Возвращаем моковые данные
//     });
//   }
// };

// export default ApiSession;
