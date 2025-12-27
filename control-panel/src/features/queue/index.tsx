import ToggleButtonGroup from "@mui/material/ToggleButtonGroup";
import ToggleButton from "@mui/material/ToggleButton";
import { useState } from "react";
const queueTypes = {
	LIVE: "LIVE",
	HISTORY: "HISTORY",
} as const;
export default () => {
	const [selectedQueueType, setSelectedQueueType] = useState(queueTypes.LIVE);
	return (
		<>
			<ToggleButtonGroup
				color="primary"
				value={selectedQueueType}
				exclusive
				onChange={(_, v) => {
					if (v) setSelectedQueueType(v);
				}}
				aria-label="Platform"
			>
				<ToggleButton value={queueTypes.LIVE}>LIVE</ToggleButton>
				<ToggleButton value={queueTypes.HISTORY}>History</ToggleButton>
			</ToggleButtonGroup>
		</>
	);
};
