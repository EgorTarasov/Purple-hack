import { NavLink, Link, useLocation } from "react-router-dom";
import { Button } from "../ui/button";
import logoPath from "../../assets/MediumLogo.svg";
import { topBarLinks } from "@/constants";
import { INavLink } from "@/models";
import uuid from "react-uuid";

const TopBar = () => {
	const { pathname } = useLocation();
	return (
		<section className="flex w-full justify-between p-6">
			<Link to="/" className="flex gap-3 items-center">
				<img src={logoPath} alt="logo" width={130} height={325} />
			</Link>

			<ul className="flex gap">
				{topBarLinks.map((link: INavLink) => {
					const index = pathname.lastIndexOf("/");
					const concatedPathname = pathname.substring(0, index);
					const isActive = pathname === link.route || concatedPathname === link.route;
					let linkTo = link.route;
					if(link.route === '/chat') linkTo += `/${uuid()}`;

					return (
						<li key={link.label}>
							<NavLink to={linkTo} className="flex-center gap-3">
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
