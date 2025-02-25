package tyumi

// An ID for a loaded resource. Use these to refer to resources
type ResourceID int

var resourceNotFound ResourceID = 0
var invalidResource ResourceID = 0

type resource interface {
	Ready() bool
	Unload()

	getPath() string
}

var resourceCache []resource

func addResourceToCache(res resource) (id ResourceID) {
	if res == nil {
		return
	}

	resourceCache = append(resourceCache, res)

	// remember: resource IDs are offset by 1 so 0 can be the invalid resource id.
	return ResourceID(len(resourceCache))
}

func getResource[T resource](resource_id ResourceID) (res T) {
	if resource_id == invalidResource || int(resource_id) > len(resourceCache) {
		return
	}

	resource := resourceCache[resource_id - 1]
	if r, ok := resource.(T); ok {
		return r
	} else {
		return
	}
}

// scans the cache to see if the resource at path has already been loaded. returns resourceNotFound
// if unsuccessful
func getResourceIDByPath(path string) ResourceID {
	for id, res := range resourceCache {
		if res == nil {
			continue
		}

		if res.getPath() == path {
			return ResourceID(id)
		}
	}

	return resourceNotFound
}

type Resource struct {
	path        string //path to resource on disk, used to prevent duplicate loads
	platform_id int    //the id assigned to this resource by the platform
	ready       bool   //true if resource was successfully loaded and has not been unloaded
}

// Reports if this resource is has been loaded successfully and has not been unloaded.
func (r Resource) Ready() bool {
	return r.ready
}

// Unloads the resource.
func (r *Resource) Unload() {
	if !r.ready {
		return
	}

	r.ready = false
}

func (r Resource) getPath() string {
	return r.path
}
