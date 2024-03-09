import bubbleDefault from "../../assets/bubbles-svg/GazpromBubble.svg";

const LogoBubble = ({ bubblePath }: { bubblePath: string }) => {
	return (
		<div className="animate-bounce">
			<img src={bubblePath ? bubblePath : bubbleDefault} alt="company-logo" />
		</div>
	);
};

export default LogoBubble;