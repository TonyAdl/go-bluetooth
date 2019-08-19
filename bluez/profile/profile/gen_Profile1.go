package profile



import (
   "sync"
   "github.com/muka/go-bluetooth/bluez"
   "github.com/godbus/dbus"
)

var Profile1Interface = "org.bluez.Profile1"


// NewProfile1 create a new instance of Profile1
//
// Args:
// - servicePath: unique name
// - objectPath: freely definable
func NewProfile1(servicePath string, objectPath dbus.ObjectPath) (*Profile1, error) {
	a := new(Profile1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: Profile1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	return a, nil
}


/*
Profile1 Profile hierarchy

*/
type Profile1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*Profile1Properties
}

// Profile1Properties contains the exposed properties of an interface
type Profile1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

//Lock access to properties
func (p *Profile1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *Profile1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *Profile1) Close() {
	
	a.client.Disconnect()
}

// Path return Profile1 object path
func (a *Profile1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return Profile1 dbus client
func (a *Profile1) Client() *bluez.Client {
	return a.client
}

// Interface return Profile1 interface
func (a *Profile1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *Profile1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

	if a.objectManagerSignal == nil {
		if a.objectManager == nil {
			om, err := bluez.GetObjectManager()
			if err != nil {
				return nil, nil, err
			}
			a.objectManager = om
		}

		s, err := a.objectManager.Register()
		if err != nil {
			return nil, nil, err
		}
		a.objectManagerSignal = s
	}

	cancel := func() {
		if a.objectManagerSignal == nil {
			return
		}
		a.objectManagerSignal <- nil
		a.objectManager.Unregister(a.objectManagerSignal)
		a.objectManagerSignal = nil
	}

	return a.objectManagerSignal, cancel, nil
}




/*
Release 
			This method gets called when the service daemon
			unregisters the profile. A profile can use it to do
			cleanup tasks. There is no need to unregister the
			profile, because when this method gets called it has
			already been unregistered.


*/
func (a *Profile1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

/*
NewConnection 
			This method gets called when a new service level
			connection has been made and authorized.

			Common fd_properties:

			uint16 Version		Profile version (optional)
			uint16 Features		Profile features (optional)

			Possible errors: org.bluez.Error.Rejected
			                 org.bluez.Error.Canceled


*/
func (a *Profile1) NewConnection(device dbus.ObjectPath, fd int32, fd_properties map[string]interface{}) error {
	
	return a.client.Call("NewConnection", 0, device, fd, fd_properties).Store()
	
}

/*
RequestDisconnection 
			This method gets called when a profile gets
			disconnected.

			The file descriptor is no longer owned by the service
			daemon and the profile implementation needs to take
			care of cleaning up all connections.

			If multiple file descriptors are indicated via
			NewConnection, it is expected that all of them
			are disconnected before returning from this
			method call.

			Possible errors: org.bluez.Error.Rejected
			                 org.bluez.Error.Canceled

*/
func (a *Profile1) RequestDisconnection(device dbus.ObjectPath) error {
	
	return a.client.Call("RequestDisconnection", 0, device).Store()
	
}

