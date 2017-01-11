package duktape

import . "gopkg.in/check.v1"

func (s *DuktapeSuite) TestPevalString(c *C) {
	s.ctx.EvalString(`"Golang love Duktape!"`)
	c.Assert(s.ctx.IsString(-1), Equals, true)
	c.Assert(s.ctx.GetString(-1), Equals, "Golang love Duktape!")
}

func (s *DuktapeSuite) TestPevalString_Error(c *C) {
	err := s.ctx.PevalString("var = 'foo';")
	c.Assert(err.(*Error).Type, Equals, "SyntaxError")
}

func (s *DuktapeSuite) TestPevalFile_Error(c *C) {
	err := s.ctx.PevalFile("foo.js")
	c.Assert(err.(*Error).Message, Equals, "no sourcecode")
}

func (s *DuktapeSuite) TestPcompileString(c *C) {
	err := s.ctx.PcompileString(CompileFunction, "foo")
	c.Assert(err.(*Error).Type, Equals, "SyntaxError")
	c.Assert(err.(*Error).LineNumber, Equals, 1)
}

func (s *DuktapeSuite) TestPushErrorObject(c *C) {
	s.ctx.PushErrorObject(ErrType, "Got an error thingy: %v", 5)
	s.assertErrorInCtx(c, ErrType, "TypeError: Got an error thingy: 5")
}

func (s *DuktapeSuite) TestPushErrorObjectVa(c *C) {
	s.ctx.PushErrorObjectVa(ErrURI, "Got an error thingy: %x %s %s", 0xdeadbeef, "is", "tasty")
	s.assertErrorInCtx(c, ErrURI, "URIError: Got an error thingy: deadbeef is tasty")
}

func (s *DuktapeSuite) assertErrorInCtx(c *C, code int, msg string) {
	c.Assert(s.ctx.IsError(-1), Equals, true)
	c.Assert(s.ctx.GetErrorCode(-1), Equals, code)
	c.Assert(s.ctx.SafeToString(-1), Equals, msg)
}
