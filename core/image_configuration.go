package core

// ImageConfiguration struct
// Properties used to generate new image
type ImageConfiguration struct {
	// ServerConfiguration *ServerConfiguration
	ID        string
	Width     int
	Height    int
	Filename  string
	Format    string
	Source    string
	Quality   uint
	Namespace string
}

// func (ic *ImageConfiguration) ImageDirectory() string {
// 	// ic.ServerConfiguration.Adapters.Paths
// 	id := strings.Join(ic.IDPartitions(), "/")
//
// 	return fmt.Sprintf("%s/%s", ic.Namespace, id)
// }

// func (ic *ImageConfiguration) IDPartitions() []string {
// 	digits := fmt.Sprintf("%06s", ic.ID)
//
// 	return []string{digits[0:2], digits[2:4], digits[4:6]}
// }
//
// func (ic *ImageConfiguration) LocalDestinationDirectory() string {
// 	return ic.ServerConfiguration.LocalBasePath + "/" + ic.ImageDirectory()
// }
//
// func (ic *ImageConfiguration) LocalOriginalImagePath() string {
// 	return ic.LocalDestinationDirectory() + "/original"
// }
//
// func (ic *ImageConfiguration) LocalResizedImagePath() string {
// 	return ic.LocalDestinationDirectory() + "/" + ic.Filename
// }
//
// func (sc *ServerConfiguration) MantaResizedImagePath(ic *ImageConfiguration) string {
// 	return sc.RemoteBasePath + "/" + ic.ImageDirectory() + "/" + ic.Filename
// }
//
// func (sc *ServerConfiguration) MantaOriginalImagePath(ic *ImageConfiguration) string {
// 	return sc.RemoteBasePath + "/" + ic.ImageDirectory() + "/original"
// }
