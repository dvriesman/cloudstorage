package cloudstorage

import "os"

type Store interface {
	//NewObject creates a new empty object backed by the cloud store
	//  This new object isn't' synced/created in the backing store
	//  until the object is Closed/Sync'ed.
	NewObject(o string) (Object, error)
	//WriteObject convenience method like ioutil.WriteFile
	//  The object will be created if it doesn't already exists.
	//  If the object does exists it will be overwritten.
	WriteObject(o string, meta map[string]string, b []byte) error
	//Get returns the object from the cloud store.   The object
	//  isn't opened already, see Object.Open()
	Get(o string) (Object, error)
	//GetAndOpen is a convenience method that combines Store.Get() and Object.Open() into
	// a single call.
	GetAndOpen(o string, readonly bool) (Object, error)
	//List takes a prefix query and returns an array of unopened objects
	// that have the given prefix.
	List(query Query) (Objects, error)
	//Delete removes the object from the cloud store.   Any Objects which have
	// had Open() called should work as normal.
	Delete(o string) error

	String() string
}

//Objects are just a collection of Object(s).  Used as the results for store.List commands.
type Objects []Object

func (o Objects) Len() int           { return len(o) }
func (o Objects) Less(i, j int) bool { return o[i].Name() < o[j].Name() }
func (o Objects) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

//Object is a handle to a cloud stored file/object.  Calling Open will pull the remote file onto
// your local filesystem for reading/writing.  Calling Sync/Close will push the local copy
// backup to the cloud store.
type Object interface {
	Name() string
	String() string

	MetaData() map[string]string
	SetMetaData(meta map[string]string)

	StorageSource() string
	//Open copies the remote file to a local cache and opens the cached version
	// for read/writing.  Calling Close/Sync will push the copy back to the
	// backing store.
	Open(readonly bool) error
	//Release will remove the locally cached copy of the file.  You most call Close
	// before releasing.  Release will call os.Remove(local_copy_file) so opened
	//filehandles need to be closed.
	Release() error
	//Implement io.ReadWriteCloser Open most be called before using these
	// functions.
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Sync() error
	Close() error

	//CachedCopy returns a pointer to the local cache file.  Changes made to the
	// file will flushed to the remote store when Close/Sync is called.
	CachedCopy() *os.File
}