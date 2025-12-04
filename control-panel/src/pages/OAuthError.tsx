import { Container, Row, Col, Card, Alert } from "react-bootstrap";
import { Link } from "react-router-dom";
import { XCircleFill } from "react-bootstrap-icons";

export function OAuthError() {
	return (
		<Container className="py-4">
			<Row>
				<Col>
					<Card className="mb-4">
						<Card.Header>
							<h5 className="mb-0">Twitch Authentication Failed</h5>
						</Card.Header>
						<Card.Body className="text-center">
							<XCircleFill size={64} color="#dc3545" className="mb-3" />
							<Alert variant="danger" className="mb-3">
								<strong>Authentication Failed!</strong> There was an error connecting your Twitch account.
							</Alert>
							<p className="text-muted">
								This could be due to network issues, invalid permissions, or an expired authorization.
							</p>
							<div className="d-flex justify-content-center">
								<Link to="/connect" className="btn btn-primary">
									Try Again
								</Link>
							</div>
						</Card.Body>
					</Card>
				</Col>
			</Row>
		</Container>
	);
}
