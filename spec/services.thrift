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
    3: string description
    4: bool enabled
}

/*
 * User entity
 */
struct User {
    1: string id,
    2: string name,
    3: bool enabled,
}

/*
 * Session entity
 */
struct Session {
    1: string id,
    2: Domain domain,
    3: User user,
    4: string userAgent,
    5: string remoteAddr,
    6: string createdOn,
    7: string updatedOn,
    8: string expiresOn
}

/**
 * Exception represents internal server error
 */
exception ServerError {
    1: string code,
    2: string msg
}

/**
 * Exception represents invalid arguments submitted error
 */
exception BadRequestError {
    1: string code,
    2: string msg
}

/**
 * Exception represents forbidden error
 */
exception ForbiddenError {
    1: string code,
    2: string msg
}

/**
 * Authenticator service
 */
service Authenticator {
    # Create new session
    Session createSession(1:string domain,
                          2:string name,
                          3:string password,
                          4:string userAgent,
                          5:string remoteAddr) throws (1:ServerError error1, 2:BadRequestError error2, 3:ForbiddenError error3),

    # Checking existing session by ID
    bool checkSession(1:string sessionID,
                      2:string userAgent,
                      3:string remoteAddr) throws (1:ServerError error1, 2:BadRequestError error2, 3:ForbiddenError error3),

    # Delete existing session by ID
    bool deleteSession(1:string sessionID,
                       2:string userAgent,
                       3:string remoteAddr) throws (1:ServerError error1, 2:BadRequestError error2, 3:ForbiddenError error3)

}

/**
 * ADD MORE SERVICES HERE
 */