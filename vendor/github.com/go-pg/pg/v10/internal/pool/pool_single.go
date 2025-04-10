package pool

import "context"

type SingleConnPool struct {
	pool      Pooler
	cn        *Conn
	stickyErr error
}

var _ Pooler = (*SingleConnPool)(nil)

func NewSingleConnPool(pool Pooler, cn *Conn) *SingleConnPool {
	return &SingleConnPool{
		pool: pool,
		cn:   cn,
	}
}

func (p *SingleConnPool) NewConn(ctx context.Context) (*Conn, error) {
	return p.pool.NewConn(ctx)
}

func (p *SingleConnPool) CloseConn(cn *Conn) error {
	return p.pool.CloseConn(cn)
}

func (p *SingleConnPool) Get(ctx context.Context) (*Conn, error) {
	if p.stickyErr != nil {
		return nil, p.stickyErr
	}
	return p.cn, nil
}

func (p *SingleConnPool) Put(ctx context.Context, cn *Conn) {}

func (p *SingleConnPool) Remove(ctx context.Context, cn *Conn, reason error) {
	p.cn = nil
	p.stickyErr = reason
	// If ctx is cancelled without a reason(error) value,
	// then the ctx.Error is used as the reason for why the p.cn is assigned nil.
	if reason == nil && ctx != nil {
		p.stickyErr = ctx.Err()
	}
}

func (p *SingleConnPool) Close() error {
	p.cn = nil
	p.stickyErr = ErrClosed
	return nil
}

func (p *SingleConnPool) Len() int {
	return 0
}

func (p *SingleConnPool) IdleLen() int {
	return 0
}

func (p *SingleConnPool) Stats() *Stats {
	return &Stats{}
}

func (p *SingleConnPool) GetWriteBuffer() *WriteBuffer {
	return p.pool.GetWriteBuffer()
}

func (p *SingleConnPool) PutWriteBuffer(wb *WriteBuffer) {
	p.pool.PutWriteBuffer(wb)
}

func (p *SingleConnPool) GetReaderContext() *ReaderContext {
	return p.pool.GetReaderContext()
}

func (p *SingleConnPool) PutReaderContext(rd *ReaderContext) {
	p.pool.PutReaderContext(rd)
}
