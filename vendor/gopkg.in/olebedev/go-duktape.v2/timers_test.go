package duktape

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *DuktapeSuite) TestCheckTheStack(c *C) {
	s.ctx.PushInt(1)
	err := s.ctx.PushTimers()
	c.Assert(err, IsNil)
	c.Assert(s.ctx.GetTop(), Equals, 1)
	err = s.ctx.PushTimers()
	c.Assert(err.Error(), Equals, "Timers are already defined")
}

func (s *DuktapeSuite) TestSetTimeOut(c *C) {
	ch := make(chan struct{})
	s.ctx.PushTimers()
	s.ctx.PushGlobalGoFunction("test", func(ctx *Context) int {
		ctx.PushNumber(2)
		ch <- struct{}{}
		return 1
	})
	s.ctx.PevalString(`setTimeout(test, 0);`)
	c.Assert(s.ctx.SafeToString(-1), Equals, "1")
	s.ctx.Pop()
	<-ch
	c.Assert(s.ctx.SafeToString(-1), Equals, "2")
	s.ctx.PopN(s.ctx.GetTop())
	c.Succeed()
}

func (s *DuktapeSuite) TestCrashProcess(c *C) {
	s.ctx.PushTimers()
	s.ctx.PushGlobalGoFunction("test", func(_ *Context) int {
		return 0
	})
	s.ctx.PevalString(`
		var id = setTimeout(test, 2);
	`)
}

func (s *DuktapeSuite) TestClearTimeOut(c *C) {
	ch := make(chan struct{}, 1) // buffered channel
	s.ctx.PushTimers()
	s.ctx.PushGlobalGoFunction("test", func(_ *Context) int {
		ch <- struct{}{}
		return 0
	})
	s.ctx.PevalString(`
		var id = setTimeout(test, 0);
		clearTimeout(id);
	`)
	<-time.After(2 * time.Millisecond)
	select {
	case <-ch:
		c.Fail()
	default:
		c.Succeed()
	}
}

func (s *DuktapeSuite) TestSetInterval(c *C) {
	ch := make(chan struct{}, 5)
	s.ctx.PushTimers()
	s.ctx.PushGlobalGoFunction("test", func(_ *Context) int {
		ch <- struct{}{}
		return 0
	})

	s.ctx.PevalString(`var id = setInterval(test, 0);`)
	s.ctx.Pop()
	<-ch
	<-ch
	<-ch
	s.ctx.PevalString(`clearInterval(id);`)
	s.ctx.Pop() // pop undefined

	<-time.After(4 * time.Millisecond)
	select {
	case <-ch:
		c.Fail()
	default:
		c.Succeed()
	}
}

func (s *DuktapeSuite) TestFlushTimers(c *C) {
	s.ctx.PushTimers()
	s.ctx.PevalString(`setInterval(test, 2);`)
	id := s.ctx.GetNumber(-1)
	s.ctx.FlushTimers()
	s.ctx.putTimer(id)

	c.Assert(s.ctx.GetType(-1).IsUndefined(), Equals, true)
	<-time.After(3 * time.Millisecond)
}
