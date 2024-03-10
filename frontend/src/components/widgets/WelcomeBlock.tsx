import { Button } from "../ui/button";
import { useNavigate } from "react-router-dom";

const WelcomeBlock = () => {
  const navigate = useNavigate();
  
	return (
		<div className="max-w-[700px] w-[70%] text-center">
			<div className="text-6xl mb-3 font-medium">
				Узнайте всё о работе Центрального банка РФ
			</div>
			<div className="mb-7">Задайте любой вопрос в диалоговом формате чата</div>
			<Button onClick={()=>{navigate('chat')}}>Начать</Button>
		</div>
	);
};

export default WelcomeBlock;
