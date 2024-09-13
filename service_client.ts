type CreateJournalEntryCommand = {
    name: string;
    details?: string;
    date?: date;
}

type JournalEntry = {
    tags: string;
    id: uuid;
    name: string;
    details?: string;
    date?: date;
}

type uuid = string;

type date = string;

export interface IServiceClient {
	createJournalEntry(request: CreateJournalEntryCommand): Promise<JournalEntry>;
}

type Fetcher = (method: string, data?: any) => Promise<unknown>;

export class ServiceClient implements IServiceClient {
	private fetcher: Fetcher;
	constructor(fetcher: Fetcher) {
		this.fetcher = fetcher;
	}

	async createJournalEntry(request: CreateJournalEntryCommand): Promise<JournalEntry> {
		const data = await this.fetcher("/create_journal_entry", request);
		return data as JournalEntry;
	}
}
