package rpc

import "github.com/oleksandr/idp/errs"

// AssertRole implements RBAC's interface
func (handler *IdentityProviderHandler) AssertRole(sessionID string, roleName string) (r bool, err error) {
	handler.log.Printf("AssertRole(%v, %v)", sessionID, roleName)

	session, err := handler.SessionInteractor.Find(sessionID)
	if err != nil {
		e := err.(*errs.Error)
		return false, errorToServiceError(e)
	}

	ok, err := handler.RBACInteractor.AssertRole(session.User.ID, roleName)
	if err != nil {
		handler.log.Println("ERROR:", err.Error())
		return false, errorToServiceError(err.(*errs.Error))
	}

	if ok {
		return true, nil
	}

	return false, nil
}

// AssertPermission implements RBAC's interface
func (handler *IdentityProviderHandler) AssertPermission(sessionID string, permissioName string) (r bool, err error) {
	handler.log.Printf("AssertPermission(%v, %v)", sessionID, permissioName)

	session, err := handler.SessionInteractor.Find(sessionID)
	if err != nil {
		e := err.(*errs.Error)
		return false, errorToServiceError(e)
	}

	ok, err := handler.RBACInteractor.AssertPermission(session.User.ID, permissioName)
	if err != nil {
		handler.log.Println("ERROR:", err.Error())
		return false, errorToServiceError(err.(*errs.Error))
	}

	if ok {
		return true, nil
	}

	return false, nil
}
