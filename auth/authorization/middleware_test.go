package auth_authorization

import "testing"

// Public site pages must answer anonymous visitors; a page missing here is
// redirected to "/" in production even though its own module tests pass.
func TestPublicSitePagesNeedNoLogin(t *testing.T) {
	for _, path := range []string{"/", "/about", "/consulting", "/solutions", "/products", "/lab", "/blog", "/llms.txt"} {
		if !public.matchPath(path) {
			t.Errorf("%s is not public: anonymous visitors are redirected to /", path)
		}
	}
}
