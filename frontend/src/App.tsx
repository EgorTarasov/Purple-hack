import { Route, Routes } from "react-router-dom";
import { Home, Chat } from "./root/pages";
import RootLayout from "./root/RootLayout";

const App = () => {
	return (
		<main className="flex bg-gradient-to-br from-secondary-medium to-white">
			<Routes>
				{/* public routes */}
				<Route element={<RootLayout />}>
					<Route path="/" element={<Home />} />
					<Route path="/chat/:id" element={<Chat />} />
				</Route>
			</Routes>
		</main>
	);
};

export default App;
