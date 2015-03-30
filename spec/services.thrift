/*
 * Namespaces for generate package names
 */
namespace go services
namespace py idp
namespace php idp


/*
 * Domain entity
 */
struct Domain {
    1: string id,
    2: string name,
}

/*
 * User entity
 */
struct User {
    1: string id,
    2: string name,
}

/*
 * Session entity
 */
struct Session {
    1: string id,
    2: Domain domain,
    3: User user,
}

/**
 * Exception represents internal server error
 */
exception ServerError {
    1: string msg
}

/**
 * Exception represents invalid arguments submitted error
 */
exception BadRequest {
    1: string msg
}

/**
 * Exception represents forbidden error
 */
exception Forbidden {
    1: string msg
}

/**
 * Authenticator service
 */
service Authenticator {
    Session createSession(1:string domainID, 2:string name, 3:string password) throws (1:ServerError error1, 2:BadRequest error2, 3:Forbidden error3),
    bool checkSession(1:string sessionID) throws (1:ServerError error1, 2:BadRequest error2, 3:Forbidden error3),
    bool deleteSession(1:string sessionID) throws (1:ServerError error1, 2:BadRequest error2, 3:Forbidden error3)
}