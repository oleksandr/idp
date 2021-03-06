// Autogenerated by Thrift Compiler (0.9.2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package services

import (
	"bytes"
	"fmt"
	"git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var GoUnusedProtection__ int

type Domain struct {
	Id          string `thrift:"id,1" json:"id"`
	Name        string `thrift:"name,2" json:"name"`
	Description string `thrift:"description,3" json:"description"`
	Enabled     bool   `thrift:"enabled,4" json:"enabled"`
}

func NewDomain() *Domain {
	return &Domain{}
}

func (p *Domain) GetId() string {
	return p.Id
}

func (p *Domain) GetName() string {
	return p.Name
}

func (p *Domain) GetDescription() string {
	return p.Description
}

func (p *Domain) GetEnabled() bool {
	return p.Enabled
}
func (p *Domain) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *Domain) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Id = v
	}
	return nil
}

func (p *Domain) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Name = v
	}
	return nil
}

func (p *Domain) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 3: %s", err)
	} else {
		p.Description = v
	}
	return nil
}

func (p *Domain) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return fmt.Errorf("error reading field 4: %s", err)
	} else {
		p.Enabled = v
	}
	return nil
}

func (p *Domain) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Domain"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *Domain) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("id", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:id: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Id)); err != nil {
		return fmt.Errorf("%T.id (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:id: %s", p, err)
	}
	return err
}

func (p *Domain) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:name: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return fmt.Errorf("%T.name (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:name: %s", p, err)
	}
	return err
}

func (p *Domain) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("description", thrift.STRING, 3); err != nil {
		return fmt.Errorf("%T write field begin error 3:description: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Description)); err != nil {
		return fmt.Errorf("%T.description (3) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 3:description: %s", p, err)
	}
	return err
}

func (p *Domain) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("enabled", thrift.BOOL, 4); err != nil {
		return fmt.Errorf("%T write field begin error 4:enabled: %s", p, err)
	}
	if err := oprot.WriteBool(bool(p.Enabled)); err != nil {
		return fmt.Errorf("%T.enabled (4) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 4:enabled: %s", p, err)
	}
	return err
}

func (p *Domain) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Domain(%+v)", *p)
}

type User struct {
	Id      string `thrift:"id,1" json:"id"`
	Name    string `thrift:"name,2" json:"name"`
	Enabled bool   `thrift:"enabled,3" json:"enabled"`
}

func NewUser() *User {
	return &User{}
}

func (p *User) GetId() string {
	return p.Id
}

func (p *User) GetName() string {
	return p.Name
}

func (p *User) GetEnabled() bool {
	return p.Enabled
}
func (p *User) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *User) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Id = v
	}
	return nil
}

func (p *User) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Name = v
	}
	return nil
}

func (p *User) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return fmt.Errorf("error reading field 3: %s", err)
	} else {
		p.Enabled = v
	}
	return nil
}

func (p *User) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("User"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *User) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("id", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:id: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Id)); err != nil {
		return fmt.Errorf("%T.id (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:id: %s", p, err)
	}
	return err
}

func (p *User) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:name: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return fmt.Errorf("%T.name (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:name: %s", p, err)
	}
	return err
}

func (p *User) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("enabled", thrift.BOOL, 3); err != nil {
		return fmt.Errorf("%T write field begin error 3:enabled: %s", p, err)
	}
	if err := oprot.WriteBool(bool(p.Enabled)); err != nil {
		return fmt.Errorf("%T.enabled (3) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 3:enabled: %s", p, err)
	}
	return err
}

func (p *User) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("User(%+v)", *p)
}

type Session struct {
	Id         string  `thrift:"id,1" json:"id"`
	Domain     *Domain `thrift:"domain,2" json:"domain"`
	User       *User   `thrift:"user,3" json:"user"`
	UserAgent  string  `thrift:"userAgent,4" json:"userAgent"`
	RemoteAddr string  `thrift:"remoteAddr,5" json:"remoteAddr"`
	CreatedOn  string  `thrift:"createdOn,6" json:"createdOn"`
	UpdatedOn  string  `thrift:"updatedOn,7" json:"updatedOn"`
	ExpiresOn  string  `thrift:"expiresOn,8" json:"expiresOn"`
}

func NewSession() *Session {
	return &Session{}
}

func (p *Session) GetId() string {
	return p.Id
}

var Session_Domain_DEFAULT *Domain

func (p *Session) GetDomain() *Domain {
	if !p.IsSetDomain() {
		return Session_Domain_DEFAULT
	}
	return p.Domain
}

var Session_User_DEFAULT *User

func (p *Session) GetUser() *User {
	if !p.IsSetUser() {
		return Session_User_DEFAULT
	}
	return p.User
}

func (p *Session) GetUserAgent() string {
	return p.UserAgent
}

func (p *Session) GetRemoteAddr() string {
	return p.RemoteAddr
}

func (p *Session) GetCreatedOn() string {
	return p.CreatedOn
}

func (p *Session) GetUpdatedOn() string {
	return p.UpdatedOn
}

func (p *Session) GetExpiresOn() string {
	return p.ExpiresOn
}
func (p *Session) IsSetDomain() bool {
	return p.Domain != nil
}

func (p *Session) IsSetUser() bool {
	return p.User != nil
}

func (p *Session) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.ReadField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.ReadField6(iprot); err != nil {
				return err
			}
		case 7:
			if err := p.ReadField7(iprot); err != nil {
				return err
			}
		case 8:
			if err := p.ReadField8(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *Session) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Id = v
	}
	return nil
}

func (p *Session) ReadField2(iprot thrift.TProtocol) error {
	p.Domain = &Domain{}
	if err := p.Domain.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.Domain, err)
	}
	return nil
}

func (p *Session) ReadField3(iprot thrift.TProtocol) error {
	p.User = &User{}
	if err := p.User.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.User, err)
	}
	return nil
}

func (p *Session) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 4: %s", err)
	} else {
		p.UserAgent = v
	}
	return nil
}

func (p *Session) ReadField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 5: %s", err)
	} else {
		p.RemoteAddr = v
	}
	return nil
}

func (p *Session) ReadField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 6: %s", err)
	} else {
		p.CreatedOn = v
	}
	return nil
}

func (p *Session) ReadField7(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 7: %s", err)
	} else {
		p.UpdatedOn = v
	}
	return nil
}

func (p *Session) ReadField8(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 8: %s", err)
	} else {
		p.ExpiresOn = v
	}
	return nil
}

func (p *Session) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Session"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := p.writeField5(oprot); err != nil {
		return err
	}
	if err := p.writeField6(oprot); err != nil {
		return err
	}
	if err := p.writeField7(oprot); err != nil {
		return err
	}
	if err := p.writeField8(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *Session) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("id", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:id: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Id)); err != nil {
		return fmt.Errorf("%T.id (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:id: %s", p, err)
	}
	return err
}

func (p *Session) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("domain", thrift.STRUCT, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:domain: %s", p, err)
	}
	if err := p.Domain.Write(oprot); err != nil {
		return fmt.Errorf("%T error writing struct: %s", p.Domain, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:domain: %s", p, err)
	}
	return err
}

func (p *Session) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("user", thrift.STRUCT, 3); err != nil {
		return fmt.Errorf("%T write field begin error 3:user: %s", p, err)
	}
	if err := p.User.Write(oprot); err != nil {
		return fmt.Errorf("%T error writing struct: %s", p.User, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 3:user: %s", p, err)
	}
	return err
}

func (p *Session) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("userAgent", thrift.STRING, 4); err != nil {
		return fmt.Errorf("%T write field begin error 4:userAgent: %s", p, err)
	}
	if err := oprot.WriteString(string(p.UserAgent)); err != nil {
		return fmt.Errorf("%T.userAgent (4) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 4:userAgent: %s", p, err)
	}
	return err
}

func (p *Session) writeField5(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("remoteAddr", thrift.STRING, 5); err != nil {
		return fmt.Errorf("%T write field begin error 5:remoteAddr: %s", p, err)
	}
	if err := oprot.WriteString(string(p.RemoteAddr)); err != nil {
		return fmt.Errorf("%T.remoteAddr (5) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 5:remoteAddr: %s", p, err)
	}
	return err
}

func (p *Session) writeField6(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("createdOn", thrift.STRING, 6); err != nil {
		return fmt.Errorf("%T write field begin error 6:createdOn: %s", p, err)
	}
	if err := oprot.WriteString(string(p.CreatedOn)); err != nil {
		return fmt.Errorf("%T.createdOn (6) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 6:createdOn: %s", p, err)
	}
	return err
}

func (p *Session) writeField7(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("updatedOn", thrift.STRING, 7); err != nil {
		return fmt.Errorf("%T write field begin error 7:updatedOn: %s", p, err)
	}
	if err := oprot.WriteString(string(p.UpdatedOn)); err != nil {
		return fmt.Errorf("%T.updatedOn (7) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 7:updatedOn: %s", p, err)
	}
	return err
}

func (p *Session) writeField8(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("expiresOn", thrift.STRING, 8); err != nil {
		return fmt.Errorf("%T write field begin error 8:expiresOn: %s", p, err)
	}
	if err := oprot.WriteString(string(p.ExpiresOn)); err != nil {
		return fmt.Errorf("%T.expiresOn (8) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 8:expiresOn: %s", p, err)
	}
	return err
}

func (p *Session) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Session(%+v)", *p)
}

type ServerError struct {
	Msg   string `thrift:"msg,1" json:"msg"`
	Cause string `thrift:"cause,2" json:"cause"`
}

func NewServerError() *ServerError {
	return &ServerError{}
}

func (p *ServerError) GetMsg() string {
	return p.Msg
}

func (p *ServerError) GetCause() string {
	return p.Cause
}
func (p *ServerError) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *ServerError) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Msg = v
	}
	return nil
}

func (p *ServerError) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Cause = v
	}
	return nil
}

func (p *ServerError) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ServerError"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *ServerError) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("msg", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:msg: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Msg)); err != nil {
		return fmt.Errorf("%T.msg (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:msg: %s", p, err)
	}
	return err
}

func (p *ServerError) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("cause", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:cause: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Cause)); err != nil {
		return fmt.Errorf("%T.cause (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:cause: %s", p, err)
	}
	return err
}

func (p *ServerError) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ServerError(%+v)", *p)
}

func (p *ServerError) Error() string {
	return p.String()
}

type BadRequestError struct {
	Msg   string `thrift:"msg,1" json:"msg"`
	Cause string `thrift:"cause,2" json:"cause"`
}

func NewBadRequestError() *BadRequestError {
	return &BadRequestError{}
}

func (p *BadRequestError) GetMsg() string {
	return p.Msg
}

func (p *BadRequestError) GetCause() string {
	return p.Cause
}
func (p *BadRequestError) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *BadRequestError) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Msg = v
	}
	return nil
}

func (p *BadRequestError) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Cause = v
	}
	return nil
}

func (p *BadRequestError) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("BadRequestError"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *BadRequestError) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("msg", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:msg: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Msg)); err != nil {
		return fmt.Errorf("%T.msg (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:msg: %s", p, err)
	}
	return err
}

func (p *BadRequestError) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("cause", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:cause: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Cause)); err != nil {
		return fmt.Errorf("%T.cause (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:cause: %s", p, err)
	}
	return err
}

func (p *BadRequestError) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BadRequestError(%+v)", *p)
}

func (p *BadRequestError) Error() string {
	return p.String()
}

type UnauthorizedError struct {
	Msg   string `thrift:"msg,1" json:"msg"`
	Cause string `thrift:"cause,2" json:"cause"`
}

func NewUnauthorizedError() *UnauthorizedError {
	return &UnauthorizedError{}
}

func (p *UnauthorizedError) GetMsg() string {
	return p.Msg
}

func (p *UnauthorizedError) GetCause() string {
	return p.Cause
}
func (p *UnauthorizedError) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *UnauthorizedError) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Msg = v
	}
	return nil
}

func (p *UnauthorizedError) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Cause = v
	}
	return nil
}

func (p *UnauthorizedError) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("UnauthorizedError"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *UnauthorizedError) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("msg", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:msg: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Msg)); err != nil {
		return fmt.Errorf("%T.msg (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:msg: %s", p, err)
	}
	return err
}

func (p *UnauthorizedError) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("cause", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:cause: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Cause)); err != nil {
		return fmt.Errorf("%T.cause (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:cause: %s", p, err)
	}
	return err
}

func (p *UnauthorizedError) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UnauthorizedError(%+v)", *p)
}

func (p *UnauthorizedError) Error() string {
	return p.String()
}

type ForbiddenError struct {
	Msg   string `thrift:"msg,1" json:"msg"`
	Cause string `thrift:"cause,2" json:"cause"`
}

func NewForbiddenError() *ForbiddenError {
	return &ForbiddenError{}
}

func (p *ForbiddenError) GetMsg() string {
	return p.Msg
}

func (p *ForbiddenError) GetCause() string {
	return p.Cause
}
func (p *ForbiddenError) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *ForbiddenError) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Msg = v
	}
	return nil
}

func (p *ForbiddenError) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Cause = v
	}
	return nil
}

func (p *ForbiddenError) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ForbiddenError"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *ForbiddenError) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("msg", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:msg: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Msg)); err != nil {
		return fmt.Errorf("%T.msg (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:msg: %s", p, err)
	}
	return err
}

func (p *ForbiddenError) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("cause", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:cause: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Cause)); err != nil {
		return fmt.Errorf("%T.cause (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:cause: %s", p, err)
	}
	return err
}

func (p *ForbiddenError) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ForbiddenError(%+v)", *p)
}

func (p *ForbiddenError) Error() string {
	return p.String()
}

type NotFoundError struct {
	Msg   string `thrift:"msg,1" json:"msg"`
	Cause string `thrift:"cause,2" json:"cause"`
}

func NewNotFoundError() *NotFoundError {
	return &NotFoundError{}
}

func (p *NotFoundError) GetMsg() string {
	return p.Msg
}

func (p *NotFoundError) GetCause() string {
	return p.Cause
}
func (p *NotFoundError) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *NotFoundError) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Msg = v
	}
	return nil
}

func (p *NotFoundError) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Cause = v
	}
	return nil
}

func (p *NotFoundError) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("NotFoundError"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *NotFoundError) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("msg", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:msg: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Msg)); err != nil {
		return fmt.Errorf("%T.msg (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:msg: %s", p, err)
	}
	return err
}

func (p *NotFoundError) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("cause", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:cause: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Cause)); err != nil {
		return fmt.Errorf("%T.cause (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:cause: %s", p, err)
	}
	return err
}

func (p *NotFoundError) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("NotFoundError(%+v)", *p)
}

func (p *NotFoundError) Error() string {
	return p.String()
}
