export interface Participant {
    id: string;
    name: string;
    isScrumMaster: boolean;
    isConnected: boolean;
}

export interface Task {
    id: string;
    title: string;
    voteCounts: Record<string, number>;
    isCompleted: boolean;
    votes: Record<string, string>;
    votingStatus: string;
}

interface Room {
    id: string;
    name: string;
    participants: Participant[];
    tasks: Task[];
    numberOfTaskCompleted: number;
    numberOfCards: number[] | null;
}

export const votingStatus = {
    PENDING: 'pending',
    VOTING: 'voting',
    COMPLETED: 'completed',
};