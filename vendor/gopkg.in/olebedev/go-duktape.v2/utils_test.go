package duktape

import . "gopkg.in/check.v1"

func (s *DuktapeSuite) TestThatContextExists(c *C) {
	defer func() {
		r := recover()
		c.Assert(r, Equals, "[duktape] Context does not exists!\nYou cannot call any contexts methods after `DestroyHeap()` was called.")
	}()
	ctx := New()
	ctx.DestroyHeap()
	ctx.Must()
}
