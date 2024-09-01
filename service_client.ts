type CreateUserResponse = { 
   userId: uuid;
};

type User = { 
   email?: string;
   userId: int;
   username: string;
};

type UserPage = { 
   end: int;
   start?: int;
   users: User;
};

interface IServiceClient { 
    CreateUser(request: User): Promise<CreateUserResponse>
    GetUsers(): Promise<UserPage>
    PingUser(request: int): Promise<void>
}

class ServiceClient implements IServiceClient {
    
}