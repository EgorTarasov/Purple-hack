import { Button } from "../ui/button"


const WelcomeBlock = () => {
  return (
    <div className="max-w-[700px] w-[70%] text-center">
      <div className="text-5xl mb-3 font-medium">
      Узнайте всё о работе Центрального банка РФ
      </div>
      <div className="mb-7">
      Задайте любой вопрос в диалоговом формате чата
      </div>
      <Button>
        Начать
      </Button>
    </div>
  )
}

export default WelcomeBlock
