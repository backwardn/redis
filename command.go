package redis

import (
	"net"
	"strconv"
	"sync"
	"time"
)

var writeBuffers = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 256)
		return &buf
	},
}

func writeBuffer(prefix string) []byte {
	p := writeBuffers.Get().(*[]byte)
	return append((*p)[:0], prefix...)
}

// GET executes <https://redis.io/commands/get>.
// The return is nil if key does not exist.
func (c *Client) GET(key string) (value []byte, err error) {
	buf := writeBuffer("*2\r\n$3\r\nGET\r\n$")
	buf = appendString(buf, key)
	return c.bulkCmd(buf)
}

// BytesGET executes <https://redis.io/commands/get>.
// The return is nil if key does not exist.
func (c *Client) BytesGET(key []byte) (value []byte, err error) {
	buf := writeBuffer("*2\r\n$3\r\nGET\r\n$")
	buf = appendBytes(buf, key)
	return c.bulkCmd(buf)
}

// SET executes <https://redis.io/commands/set>.
func (c *Client) SET(key string, value []byte) error {
	buf := writeBuffer("*3\r\n$3\r\nSET\r\n$")
	buf = appendStringBytes(buf, key, value)
	return c.okCmd(buf)
}

// BytesSET executes <https://redis.io/commands/set>.
func (c *Client) BytesSET(key, value []byte) error {
	buf := writeBuffer("*3\r\n$3\r\nSET\r\n$")
	buf = appendBytesBytes(buf, key, value)
	return c.okCmd(buf)
}

// SETString executes <https://redis.io/commands/set>.
func (c *Client) SETString(key, value string) error {
	buf := writeBuffer("*3\r\n$3\r\nSET\r\n$")
	buf = appendStringString(buf, key, value)
	return c.okCmd(buf)
}

// DEL executes <https://redis.io/commands/del>.
func (c *Client) DEL(key string) (bool, error) {
	buf := writeBuffer("*2\r\n$3\r\nDEL\r\n$")
	buf = appendString(buf, key)
	removed, err := c.intCmd(buf)
	return removed != 0, err
}

// BytesDEL executes <https://redis.io/commands/del>.
func (c *Client) BytesDEL(key []byte) (bool, error) {
	buf := writeBuffer("*2\r\n$3\r\nDEL\r\n$")
	buf = appendBytes(buf, key)
	removed, err := c.intCmd(buf)
	return removed != 0, err
}

// INCR executes <https://redis.io/commands/incr>.
func (c *Client) INCR(key string) (newValue int64, err error) {
	buf := writeBuffer("*2\r\n$4\r\nINCR\r\n$")
	buf = appendString(buf, key)
	return c.intCmd(buf)
}

// BytesINCR executes <https://redis.io/commands/incr>.
func (c *Client) BytesINCR(key []byte) (newValue int64, err error) {
	buf := writeBuffer("*2\r\n$4\r\nINCR\r\n$")
	buf = appendBytes(buf, key)
	return c.intCmd(buf)
}

// INCRBY executes <https://redis.io/commands/incrby>.
func (c *Client) INCRBY(key string, increment int64) (newValue int64, err error) {
	buf := writeBuffer("*3\r\n$6\r\nINCRBY\r\n$")
	buf = appendStringInt(buf, key, increment)
	return c.intCmd(buf)
}

// BytesINCRBY executes <https://redis.io/commands/incrby>.
func (c *Client) BytesINCRBY(key []byte, increment int64) (newValue int64, err error) {
	buf := writeBuffer("*3\r\n$6\r\nINCRBY\r\n$")
	buf = appendBytesInt(buf, key, increment)
	return c.intCmd(buf)
}

// APPEND executes <https://redis.io/commands/append>.
func (c *Client) APPEND(key string, value []byte) (newLen int64, err error) {
	buf := writeBuffer("*3\r\n$6\r\nAPPEND\r\n$")
	buf = appendStringBytes(buf, key, value)
	return c.intCmd(buf)
}

// BytesAPPEND executes <https://redis.io/commands/append>.
func (c *Client) BytesAPPEND(key, value []byte) (newLen int64, err error) {
	buf := writeBuffer("*3\r\n$6\r\nAPPEND\r\n$")
	buf = appendBytesBytes(buf, key, value)
	return c.intCmd(buf)
}

// APPENDString executes <https://redis.io/commands/append>.
func (c *Client) APPENDString(key, value string) (newLen int64, err error) {
	buf := writeBuffer("*3\r\n$6\r\nAPPEND\r\n$")
	buf = appendStringString(buf, key, value)
	return c.intCmd(buf)
}

// LLEN executes <https://redis.io/commands/llen>.
// The return is 0 if key does not exist.
func (c *Client) LLEN(key string) (int64, error) {
	buf := writeBuffer("*2\r\n$4\r\nLLEN\r\n$")
	buf = appendString(buf, key)
	return c.intCmd(buf)
}

// BytesLLEN executes <https://redis.io/commands/llen>.
// The return is 0 if key does not exist.
func (c *Client) BytesLLEN(key []byte) (int64, error) {
	buf := writeBuffer("*2\r\n$4\r\nLLEN\r\n$")
	buf = appendBytes(buf, key)
	return c.intCmd(buf)
}

// LINDEX executes <https://redis.io/commands/lindex>.
// The return is nil if key does not exist.
// The return is nil if index is out of range.
func (c *Client) LINDEX(key string, index int64) (value []byte, err error) {
	buf := writeBuffer("*3\r\n$6\r\nLINDEX\r\n$")
	buf = appendStringInt(buf, key, index)
	return c.bulkCmd(buf)
}

// BytesLINDEX executes <https://redis.io/commands/lindex>.
// The return is nil if key does not exist.
// The return is nil if index is out of range.
func (c *Client) BytesLINDEX(key []byte, index int64) (value []byte, err error) {
	buf := writeBuffer("*3\r\n$6\r\nLINDEX\r\n$")
	buf = appendBytesInt(buf, key, index)
	return c.bulkCmd(buf)
}

// LRANGE executes <https://redis.io/commands/lrange>.
// The return is empty if key does not exist.
func (c *Client) LRANGE(key string, start, stop int64) (values [][]byte, err error) {
	buf := writeBuffer("*4\r\n$6\r\nLRANGE\r\n$")
	buf = appendStringIntInt(buf, key, start, stop)
	return c.arrayCmd(buf)
}

// BytesLRANGE executes <https://redis.io/commands/lrange>.
// The return is empty if key does not exist.
func (c *Client) BytesLRANGE(key []byte, start, stop int64) (values [][]byte, err error) {
	buf := writeBuffer("*4\r\n$6\r\nLRANGE\r\n$")
	buf = appendBytesIntInt(buf, key, start, stop)
	return c.arrayCmd(buf)
}

// LPOP executes <https://redis.io/commands/lpop>.
// The return is nil if key does not exist.
func (c *Client) LPOP(key string) (value []byte, err error) {
	buf := writeBuffer("*2\r\n$4\r\nLPOP\r\n$")
	buf = appendString(buf, key)
	return c.bulkCmd(buf)
}

// BytesLPOP executes <https://redis.io/commands/lpop>.
// The return is nil if key does not exist.
func (c *Client) BytesLPOP(key []byte) (value []byte, err error) {
	buf := writeBuffer("*2\r\n$4\r\nLPOP\r\n$")
	buf = appendBytes(buf, key)
	return c.bulkCmd(buf)
}

// RPOP executes <https://redis.io/commands/rpop>.
// The return is nil if key does not exist.
func (c *Client) RPOP(key string) (value []byte, err error) {
	buf := writeBuffer("*2\r\n$4\r\nRPOP\r\n$")
	buf = appendString(buf, key)
	return c.bulkCmd(buf)
}

// BytesRPOP executes <https://redis.io/commands/rpop>.
// The return is nil if key does not exist.
func (c *Client) BytesRPOP(key []byte) (value []byte, err error) {
	buf := writeBuffer("*2\r\n$4\r\nRPOP\r\n$")
	buf = appendBytes(buf, key)
	return c.bulkCmd(buf)
}

// LTRIM executes <https://redis.io/commands/ltrim>.
func (c *Client) LTRIM(key string, start, stop int64) error {
	buf := writeBuffer("*4\r\n$5\r\nLTRIM\r\n$")
	buf = appendStringIntInt(buf, key, start, stop)
	return c.okCmd(buf)
}

// BytesLTRIM executes <https://redis.io/commands/ltrim>.
func (c *Client) BytesLTRIM(key []byte, start, stop int64) error {
	buf := writeBuffer("*4\r\n$5\r\nLTRIM\r\n$")
	buf = appendBytesIntInt(buf, key, start, stop)
	return c.okCmd(buf)
}

// LSET executes <https://redis.io/commands/lset>.
func (c *Client) LSET(key string, index int64, value []byte) error {
	buf := writeBuffer("*4\r\n$4\r\nLSET\r\n$")
	buf = appendStringIntBytes(buf, key, index, value)
	return c.okCmd(buf)
}

// LSETString executes <https://redis.io/commands/lset>.
func (c *Client) LSETString(key string, index int64, value string) error {
	buf := writeBuffer("*4\r\n$4\r\nLSET\r\n$")
	buf = appendStringIntString(buf, key, index, value)
	return c.okCmd(buf)
}

// BytesLSET executes <https://redis.io/commands/lset>.
func (c *Client) BytesLSET(key []byte, index int64, value []byte) error {
	buf := writeBuffer("*4\r\n$4\r\nLSET\r\n$")
	buf = appendBytesIntBytes(buf, key, index, value)
	return c.okCmd(buf)
}

// LPUSH executes <https://redis.io/commands/lpush>.
func (c *Client) LPUSH(key string, value []byte) (newLen int64, err error) {
	buf := writeBuffer("*3\r\n$5\r\nLPUSH\r\n$")
	buf = appendStringBytes(buf, key, value)
	return c.intCmd(buf)
}

// BytesLPUSH executes <https://redis.io/commands/lpush>.
func (c *Client) BytesLPUSH(key, value []byte) (newLen int64, err error) {
	buf := writeBuffer("*3\r\n$5\r\nLPUSH\r\n$")
	buf = appendBytesBytes(buf, key, value)
	return c.intCmd(buf)
}

// LPUSHString executes <https://redis.io/commands/lpush>.
func (c *Client) LPUSHString(key, value string) (newLen int64, err error) {
	buf := writeBuffer("*3\r\n$5\r\nLPUSH\r\n$")
	buf = appendStringString(buf, key, value)
	return c.intCmd(buf)
}

// RPUSH executes <https://redis.io/commands/rpush>.
func (c *Client) RPUSH(key string, value []byte) (newLen int64, err error) {
	buf := writeBuffer("*3\r\n$5\r\nRPUSH\r\n$")
	buf = appendStringBytes(buf, key, value)
	return c.intCmd(buf)
}

// BytesRPUSH executes <https://redis.io/commands/rpush>.
func (c *Client) BytesRPUSH(key, value []byte) (newLen int64, err error) {
	buf := writeBuffer("*3\r\n$5\r\nRPUSH\r\n$")
	buf = appendBytesBytes(buf, key, value)
	return c.intCmd(buf)
}

// RPUSHString executes <https://redis.io/commands/rpush>.
func (c *Client) RPUSHString(key, value string) (newLen int64, err error) {
	buf := writeBuffer("*3\r\n$5\r\nRPUSH\r\n$")
	buf = appendStringString(buf, key, value)
	return c.intCmd(buf)
}

// HGET executes <https://redis.io/commands/hget>.
// The return is nil if key does not exist.
func (c *Client) HGET(key, field string) (value []byte, err error) {
	buf := writeBuffer("*3\r\n$4\r\nHGET\r\n$")
	buf = appendStringString(buf, key, field)
	return c.bulkCmd(buf)
}

// BytesHGET executes <https://redis.io/commands/hget>.
// The return is nil if key does not exist.
func (c *Client) BytesHGET(key, field []byte) (value []byte, err error) {
	buf := writeBuffer("*3\r\n$4\r\nHGET\r\n$")
	buf = appendBytesBytes(buf, key, field)
	return c.bulkCmd(buf)
}

// HSET executes <https://redis.io/commands/hset>.
func (c *Client) HSET(key, field string, value []byte) (newField bool, err error) {
	buf := writeBuffer("*4\r\n$4\r\nHSET\r\n$")
	buf = appendStringStringBytes(buf, key, field, value)
	created, err := c.intCmd(buf)
	return created != 0, err
}

// BytesHSET executes <https://redis.io/commands/hset>.
func (c *Client) BytesHSET(key, field, value []byte) (newField bool, err error) {
	buf := writeBuffer("*4\r\n$4\r\nHSET\r\n$")
	buf = appendBytesBytesBytes(buf, key, field, value)
	created, err := c.intCmd(buf)
	return created != 0, err
}

// HSETString executes <https://redis.io/commands/hset>.
func (c *Client) HSETString(key, field, value string) (updated bool, err error) {
	buf := writeBuffer("*4\r\n$4\r\nHSET\r\n$")
	buf = appendStringStringString(buf, key, field, value)
	replaced, err := c.intCmd(buf)
	return replaced != 0, err
}

// HDEL executes <https://redis.io/commands/hdel>.
func (c *Client) HDEL(key, field string) (bool, error) {
	buf := writeBuffer("*3\r\n$4\r\nHDEL\r\n$")
	buf = appendStringString(buf, key, field)
	removed, err := c.intCmd(buf)
	return removed != 0, err
}

// BytesHDEL executes <https://redis.io/commands/hdel>.
func (c *Client) BytesHDEL(key, field []byte) (bool, error) {
	buf := writeBuffer("*3\r\n$4\r\nHDEL\r\n$")
	buf = appendBytesBytes(buf, key, field)
	removed, err := c.intCmd(buf)
	return removed != 0, err
}

func appendBytes(buf, a []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendString(buf []byte, a string) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendBytesBytes(buf, a1, a2 []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a2)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a2...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendBytesInt(buf, a1 []byte, a2 int64) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')

	buf = appendDecimal(buf, a2)

	buf = append(buf, '\r', '\n')
	return buf
}

func appendStringBytes(buf []byte, a1 string, a2 []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a2)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a2...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendStringInt(buf []byte, a1 string, a2 int64) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')

	buf = appendDecimal(buf, a2)

	buf = append(buf, '\r', '\n')
	return buf
}

func appendStringString(buf []byte, a1, a2 string) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a2)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a2...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendBytesBytesBytes(buf, a1, a2, a3 []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a2)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a2...)
	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a3)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a3...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendBytesIntBytes(buf, a1 []byte, a2 int64, a3 []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')

	buf = appendDecimal(buf, a2)

	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a3)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a3...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendBytesIntInt(buf, a1 []byte, a2, a3 int64) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')

	buf = appendDecimal(buf, a2)

	buf = append(buf, '\r', '\n', '$')

	buf = appendDecimal(buf, a3)

	buf = append(buf, '\r', '\n')
	return buf
}

func appendStringIntBytes(buf []byte, a1 string, a2 int64, a3 []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')

	buf = appendDecimal(buf, a2)

	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a3)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a3...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendStringIntInt(buf []byte, a1 string, a2, a3 int64) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')

	buf = appendDecimal(buf, a2)

	buf = append(buf, '\r', '\n', '$')

	buf = appendDecimal(buf, a3)

	buf = append(buf, '\r', '\n')
	return buf
}

func appendStringIntString(buf []byte, a1 string, a2 int64, a3 string) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')

	buf = appendDecimal(buf, a2)

	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a3)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a3...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendStringStringBytes(buf []byte, a1, a2 string, a3 []byte) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a2)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a2...)
	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a3)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a3...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendStringStringString(buf []byte, a1, a2, a3 string) []byte {
	buf = strconv.AppendUint(buf, uint64(len(a1)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a1...)
	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a2)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a2...)
	buf = append(buf, '\r', '\n', '$')
	buf = strconv.AppendUint(buf, uint64(len(a3)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, a3...)
	buf = append(buf, '\r', '\n')
	return buf
}

func appendDecimal(buf []byte, v int64) []byte {
	sizeOffset := len(buf)
	sizeOneDecimal := v > -1e8 && v < 1e9
	if sizeOneDecimal {
		buf = append(buf, 0, '\r', '\n')
	} else {
		buf = append(buf, 0, 0, '\r', '\n')
	}

	intOffset := len(buf)
	buf = strconv.AppendInt(buf, v, 10)
	size := len(buf) - intOffset
	if sizeOneDecimal {
		buf[sizeOffset] = byte(size + '0')
	} else {
		buf[sizeOffset] = byte(size/10 + '0')
		buf[sizeOffset+1] = byte(size%10 + '0')
	}
	return buf
}

func (c *Client) okCmd(buf []byte) error {
	parser := okParsers.Get().(okParser)
	if err := c.send(buf, parser); err != nil {
		return err
	}

	// await response
	err := <-parser
	okParsers.Put(parser)
	return err
}

func (c *Client) intCmd(buf []byte) (int64, error) {
	parser := intParsers.Get().(intParser)
	if err := c.send(buf, parser); err != nil {
		return 0, err
	}

	// await response
	resp := <-parser
	intParsers.Put(parser)
	return resp.Int, resp.Err
}

func (c *Client) bulkCmd(buf []byte) ([]byte, error) {
	parser := bulkParsers.Get().(bulkParser)
	if err := c.send(buf, parser); err != nil {
		return nil, err
	}

	// await response
	resp := <-parser
	bulkParsers.Put(parser)
	return resp.Bytes, resp.Err
}

func (c *Client) arrayCmd(buf []byte) ([][]byte, error) {
	parser := arrayParsers.Get().(arrayParser)
	if err := c.send(buf, parser); err != nil {
		return nil, err
	}

	// await response
	resp := <-parser
	arrayParsers.Put(parser)
	return resp.Array, resp.Err
}

func (c *Client) send(buf []byte, callback parser) error {
	var conn net.Conn
	select {
	case conn = <-c.writeSem:
		break // lock aquired
	case err := <-c.offline:
		return err
	}

	// send command
	if c.timeout != 0 {
		conn.SetWriteDeadline(time.Now().Add(c.timeout))
	}
	if _, err := conn.Write(buf); err != nil {
		// The write semaphore is not released.
		c.writeErr <- struct{}{} // does not block
		return err
	}

	// expect response
	c.queue <- callback

	// release lock
	c.writeSem <- conn

	writeBuffers.Put(&buf)

	return nil
}
