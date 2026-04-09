export interface ShellState {
	title: string;
}

export function createShellState(): ShellState {
	return {
		title: "Architecture Skeleton"
	};
}
