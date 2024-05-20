namespace go example

enum TErrorCode{
    EGood = 0,
    ENotFound = -1,
    EUnknown = -2 ,
    EDataExisted = -3
}

struct User{
    1: string id,
    2: string name,
    3: i32 age,
}

struct TDataResult{
    1: TErrorCode errorCode,
    2: optional User data
}

struct TListDataResult{
    1: TErrorCode errorCode,
    2: optional list<User> data,
}


service UserService{
    TDataResult postUser(1:User user), 
    TDataResult putUser(1:string key, 2: User data),
    TListDataResult getListUser(1:list<string> data ),
    TErrorCode removeUser(1:string key)
}

service UserStorageService extends UserService {

}