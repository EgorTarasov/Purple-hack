import { Route, Routes } from "react-router-dom";
import { Home, Chat, ChatHistory, ChatHistoryShared } from "./root/pages";
import RootLayout from "./root/RootLayout";
import { Toaster } from "./components/ui/toaster";

const App = () => {
	return (
		<>
			<main className="flex bg-gradient-to-br from-secondary-medium to-white">
				<Routes>
					{/* public routes */}
					<Route element={<RootLayout />}>
						<Route path="/" element={<Home />} />
						<Route path="/chat/:id" element={<Chat />} />
						<Route path="/history/:id" element={<ChatHistory />} />
						<Route path="/history-shared/:id" element={<ChatHistoryShared />} />
					</Route>
				</Routes>
				<Toaster/>
			</main>
		</>
	);
};

export default App;
