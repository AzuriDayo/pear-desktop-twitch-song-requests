import { useAppSelector } from "../../app/hooks";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import Avatar from "@mui/material/Avatar";
import Divider from "@mui/material/Divider";
import ListItemText from "@mui/material/ListItemText";
import Typography from "@mui/material/Typography";

export default () => {
	const { song_queue, isLoaded } = useAppSelector(
		(state) => state.songQueueState,
	);

	return (
		<div>
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
									<>
										<ListItem alignItems="flex-start">
											<ListItemAvatar>
												<Avatar alt={`${title} - ${artist}`} src={imageUrl} />
											</ListItemAvatar>
											<ListItemText
												primary={requested_by + (is_ninja ? " 🥷" : "")}
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
									</>
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
