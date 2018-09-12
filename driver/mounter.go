package driver

// Mounter is responsible for formatting and mounting volumes
type Mounter interface {
	// Format formats the source with the given filesystem type
	Format(source, fsType string) error

	// Mount mounts source to target with the given fstype and options.
	Mount(source, target, fsType string, options ...string) error

	// Unmount unmounts the given target
	Unmount(target string) error

	// IsFormatted checks whether the source device is formatted or not. It
	// returns true if the source device is already formatted.
	IsFormatted(source string) (bool, error)

	// IsMounted checks whether the source device is mounted to the target
	// path. Source can be empty. In that case it only checks whether the
	// device is mounted or not.
	// It returns true if it's mounted.
	IsMounted(source, target string) (bool, error)
}

type mounter struct{}

func (m *mounter) Format(source, fsType string) error {

	return nil
}

func (m *mounter) Mount(source, target, fsType string, opts ...string) error {
	return nil
}

func (m *mounter) Unmount(target string) error {
	return nil
}

func (m *mounter) IsFormatted(source string) (bool, error) {
	return false, nil
}

func (m *mounter) IsMounted(source, target string) (bool, error) {
	return false, nil
}
