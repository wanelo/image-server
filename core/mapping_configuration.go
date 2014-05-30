package core

type MapperConfiguration struct {
	NamespaceMappings map[string]string
}

// NamespaceMapping Maps a url namespace with a source path i.e 'p' => 'product/images'
type NamespaceMapping struct {
	Namespace string
	Source    string
}
