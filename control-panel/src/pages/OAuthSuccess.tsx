import { Container, Row, Col, Card, Alert } from "react-bootstrap";
import { Link } from "react-router-dom";
import { CheckCircleFill } from "react-bootstrap-icons";

export function OAuthSuccess() {
	return (
		<Container className="py-4">
			<Row>
				<Col>
					<Card className="mb-4">
						<Card.Header>
							<h5 className="mb-0">Twitch Authentication Successful</h5>
						</Card.Header>
						<Card.Body className="text-center">
							<CheckCircleFill size={64} color="#28a745" className="mb-3" />
							<Alert variant="success" className="mb-3">
								<strong>Success!</strong> Your Twitch account has been successfully connected.
							</Alert>
							<p className="text-muted">
								You can now use song request features in your Twitch chat.
							</p>
							<div className="d-flex justify-content-center">
								<Link to="/" className="btn btn-primary">
									Go to Homepage
								</Link>
							</div>
						</Card.Body>
					</Card>
				</Col>
			</Row>
		</Container>
	);
}
