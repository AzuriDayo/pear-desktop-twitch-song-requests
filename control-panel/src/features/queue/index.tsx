import ToggleButtonGroup from "@mui/material/ToggleButtonGroup";
import ToggleButton from "@mui/material/ToggleButton";
import { useState } from "react";
import { useNavigate } from "react-router";
import History from "./History";

const queueTypes = {
	CURRENT: "CURRENT",
	HISTORY: "HISTORY",
} as const;

type EQueueTypes = (typeof queueTypes)[keyof typeof queueTypes];

export default () => {
	const navigate = useNavigate();
	const [selectedQueueType, setSelectedQueueType] = useState<EQueueTypes>(
		queueTypes.CURRENT,
	);
	const handleBack = () => {
		navigate(-1); // Goes back one step in the history stack
	};
	return (
		<div>
			<button onClick={handleBack}>Go Back</button>
			<br />
			<br />
			<ToggleButtonGroup
				color="primary"
				value={selectedQueueType}
				exclusive
				onChange={(_, v) => {
					if (v) setSelectedQueueType(v);
				}}
				aria-label="Platform"
			>
				<ToggleButton value={queueTypes.CURRENT}>Current queue</ToggleButton>
				<ToggleButton value={queueTypes.HISTORY}>Show History</ToggleButton>
			</ToggleButtonGroup>

			<br />
			<br />
			<br />
			{(() => {
				switch (selectedQueueType) {
					case queueTypes.CURRENT:
						return <></>;
					case queueTypes.HISTORY:
						return <History></History>;
					default:
						return <></>;
				}
			})()}
		</div>
	);
};
