import { Button } from "../ui/button";
import { useNavigate } from "react-router-dom";
import uuid from "react-uuid";

const WelcomeBlock = () => {
  const navigate = useNavigate();
  
	return (
		<div className="max-w-[700px] w-[90%] text-center mx-auto md:w-[80%] lg:w-[60%] xl:w-[50%]">
			<div className="text-3xl md:text-4xl lg:text-6xl mb-3 font-medium">
				Узнайте всё о работе Центрального банка РФ
			</div>
			<div className="mb-7">Задайте любой вопрос в диалоговом формате чата</div>
			<Button onClick={()=>{navigate(`chat/${uuid()}`)}}>Начать</Button>
		</div>
	);
};

export default WelcomeBlock;
