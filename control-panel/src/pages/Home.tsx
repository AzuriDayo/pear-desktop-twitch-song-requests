import { Container, Row, Col } from "react-bootstrap";
import { MusicPlayer } from "../features/musicplayer/MusicPlayer";

export function Home() {
	return (
		<Container className="py-4">
			<Row>
				<Col>
					<MusicPlayer />
				</Col>
			</Row>
		</Container>
	);
}
