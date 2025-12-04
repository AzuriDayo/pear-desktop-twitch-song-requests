import { Navbar, Nav, Container } from "react-bootstrap";
import { Link, useLocation } from "react-router-dom";

export function NavBar() {
	const location = useLocation();

	return (
		<Navbar className="mb-4 fancy-navbar" style={{ borderRadius: "12px" }}>
			<Container>
				<Navbar.Brand className="d-flex align-items-center fancy-navbar-brand mb-4">
					<div className="fancy-brand-icon">🎵</div>
					Twitch Song Requests for Pear Desktop
				</Navbar.Brand>
				<Nav className="ms-auto">
					<Nav.Link
						as={Link}
						to="/"
						className={`fancy-nav-link ${location.pathname === "/" ? "active" : ""}`}
					>
						<span className="d-flex align-items-center">
							<span className="me-2">🏠</span>
							Home
						</span>
					</Nav.Link>
					<Nav.Link
						as={Link}
						to="/connect"
						className={`fancy-nav-link ${location.pathname === "/connect" ? "active" : ""}`}
					>
						<span className="d-flex align-items-center">
							<span className="me-2">🔗</span>
							Connect with Twitch
						</span>
					</Nav.Link>
				</Nav>
			</Container>
		</Navbar>
	);
}
