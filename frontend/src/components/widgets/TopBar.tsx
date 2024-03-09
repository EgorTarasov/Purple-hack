import { NavLink, Link, useLocation } from "react-router-dom";
import { Button } from "../ui/button";
import logoPath from "../../assets/MediumLogo.svg";
import { topBarLinks } from "@/constants";
import { INavLink } from "@/models";

const TopBar = () => {
	const { pathname } = useLocation();
	return (
		<section className="flex w-full justify-between p-6">
			<Link to="/" className="flex gap-3 items-center">
				<img src={logoPath} alt="logo" width={130} height={325} />
			</Link>

			<ul className="flex gap">
				{topBarLinks.map((link: INavLink) => {
					const isActive = pathname === link.route;

					return (
						<li key={link.label}>
							<NavLink to={link.route} className="flex-center gap-3">
								<Button variant="link" className={`${isActive && "underline"}`}>
									{link.label}
								</Button>
							</NavLink>
						</li>
					);
				})}
			</ul>
		</section>
	);
};

export default TopBar;
