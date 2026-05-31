import React, { useState } from "react";
import { useAppDispatch, useAppSelector } from "../../app/hooks";
import { removeSongAtIndex } from "../twitchws/songQueueSlice";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import Avatar from "@mui/material/Avatar";
import Divider from "@mui/material/Divider";
import ListItemText from "@mui/material/ListItemText";
import Typography from "@mui/material/Typography";
import IconButton from "@mui/material/IconButton";
import DeleteIcon from "@mui/icons-material/Delete";

const Queue = () => {
	const { song_queue, isLoaded } = useAppSelector(
		(state) => state.songQueueState,
	);
	const playerState = useAppSelector((state) => state.musicPlayerState);
	const dispatch = useAppDispatch();

	const [deletingIdx, setDeletingIdx] = useState<number | null>(null);

	const handleDelete = async (i: number) => {
		if (deletingIdx !== null) return;
		setDeletingIdx(i);
		try {
			// API uses 1-based index, matching !delsong # semantics
			const res = await fetch(`/api/v1/queue/${i + 1}`, {
				method: "DELETE",
			});
			if (res.ok) {
				dispatch(removeSongAtIndex({ index: i }));
			} else {
				console.error("Failed to delete song from queue:", res.status);
			}
		} catch (err) {
			console.error("Error deleting song from queue:", err);
		} finally {
			setDeletingIdx(null);
		}
	};

	return (
		<div>
			<div>
				<h4>Currently playing:</h4>
				<a
					target="_blank"
					href={`${playerState.videoUrl}`}
					style={{
						display: "flex",
						flexDirection: "row",
						alignItems: "center",
						justifyContent: "center",
					}}
				>
					<img
						style={{
							maxHeight: "5vw",
							marginRight: "20px",
							borderRadius: "3px",
						}}
						src={`${playerState.albumArtUrl}`}
					></img>
					<span>{`${playerState.artistName} - ${playerState.songName}`}</span>
				</a>
			</div>
			<br />
			<br />
			{isLoaded ? (
				song_queue.length > 0 ? (
					<List
						sx={{ width: "100%", maxWidth: 360, bgcolor: "background.paper" }}
					>
						{song_queue.map(
							(
								{
									requested_by,
									song: { artist, imageUrl, title, videoId },
									is_ninja,
								},
								i,
							) => {
								return (
									<React.Fragment key={videoId}>
										<ListItem
											alignItems="flex-start"
											secondaryAction={
												<IconButton
													edge="end"
													aria-label={`Remove ${title} from queue`}
													onClick={() => handleDelete(i)}
													disabled={deletingIdx !== null}
												>
													<DeleteIcon fontSize="small" />
												</IconButton>
											}
										>
											<ListItemAvatar>
												<Avatar alt={`${title} - ${artist}`} src={imageUrl} />
											</ListItemAvatar>
											<ListItemText
												primary={
													`#${i + 1} ` + requested_by + (is_ninja ? " 🥷" : "")
												}
												secondary={
													<a
														href={`https://youtu.be/${videoId}`}
														target="_blank"
													>
														<Typography
															component="span"
															variant="body2"
															sx={{ color: "text.primary", display: "inline" }}
														>
															{title}
														</Typography>
														{` — ${artist}`}
													</a>
												}
											/>
										</ListItem>
										{i !== song_queue.length - 1 && (
											<Divider variant="inset" component="li" />
										)}
									</React.Fragment>
								);
							},
						)}
					</List>
				) : (
					<div>Empty queue</div>
				)
			) : (
				<div>Loading...</div>
			)}
		</div>
	);
};

export default Queue;
