import Paper from "@mui/material/Paper";

import { useEffect, useState } from "react";
import {
	DataGrid,
	GridRenderCellParams,
	type GridColDef,
} from "@mui/x-data-grid";

interface IRequesterData {
	video_id: string;
	twitch_username: string;
	requested_at: string;
	is_ninja: boolean;
	song_title: string;
	artist_name: string;
	image_url: string;
}

const History = () => {
	const [paginationModel, setPaginationModel] = useState({
		page: 0, // page index is 0-based
		pageSize: 10,
	});
	const [rows, setRows] = useState([]);
	const [rowCount, setRowCount] = useState(0);
	const [loading, setLoading] = useState(false);

	const fetchDataFromServer = async (page: number, pageSize: number) => {
		setLoading(true);
		// Replace this with your actual API call (e.g., using axios or fetch)
		const response = await fetch(
			`/api/v1/requesters/history?page=${page}&perPage=${pageSize}`,
		);
		const data = await response.json();
		setRows(data.items); // The data for the current page
		setRowCount(data.max_results); // The total count of all data
		setLoading(false);
	};

	useEffect(() => {
		fetchDataFromServer(paginationModel.page, paginationModel.pageSize);
	}, [paginationModel]); // Re-fetch data whenever pagination model changes

	const columns = [
		{
			headerName: "Requester",
			field: "twitch_username",
			sortable: false,
			filterable: false,
			hideable: false,
			disableColumnMenu: true,
			width: 150,
			renderCell: (params: GridRenderCellParams<IRequesterData>) => (
				<span>{params.value + (params.row.is_ninja ? " 🥷" : "")}</span>
			),
		},
		{
			headerName: "Requested Song",
			field: "song_title",
			minWidth: 500,
			sortable: false,
			filterable: false,
			hideable: false,
			disableColumnMenu: true,
			renderCell: (params: GridRenderCellParams<IRequesterData>) => {
				return (
					<div style={{ display: "flex", flexDirection: "row" }}>
						<div>
							<a
								href={`https://youtu.be/${params.row.video_id}`}
								target="_blank"
							>
								<img
									style={{
										height: "50px",
										width: "50px",
										overflow: "hidden",
										objectFit: "cover",
									}}
									src={params.row.image_url}
									alt={params.row.song_title + " - " + params.row.artist_name}
								/>
							</a>
						</div>
						<div>
							<a
								href={`https://youtu.be/${params.row.video_id}`}
								target="_blank"
							>
								{params.row.song_title + " - " + params.row.artist_name}
							</a>
						</div>
					</div>
				);
			},
		},
		{
			field: "requested_at",
			headerName: "Requested At",
			minWidth: 250,
			sortable: false,
			filterable: false,
			hideable: false,
			disableColumnMenu: true,
			renderCell: (params: GridRenderCellParams<IRequesterData>) => (
				<span>{new Date(params.value).toLocaleString()}</span>
			),
		},
	] as GridColDef<IRequesterData>[];

	return (
		<div>
			<Paper>
				<DataGrid
					virtualizeColumnsWithAutoRowHeight
					rows={rows}
					columns={columns}
					rowCount={rowCount}
					loading={loading}
					paginationMode="server"
					paginationModel={paginationModel}
					onPaginationModelChange={setPaginationModel}
					pageSizeOptions={[10, 20, 50, 75, 100]} // Configure available page sizes
					disableRowSelectionOnClick
				/>
			</Paper>
		</div>
	);
};

export default History;
