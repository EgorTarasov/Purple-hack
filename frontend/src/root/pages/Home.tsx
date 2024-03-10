import WelcomeBlock from "@/components/widgets/WelcomeBlock";
import BgLines from "@/components/widgets/BgLines";

const Home = () => {
	return (
		<>
			<div className="flex justify-center items-center relative z-50 p-20">
				<WelcomeBlock />
			</div>
			<BgLines />
		</>
	);
};

export default Home;
