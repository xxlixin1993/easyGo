package messageQueue

import (
    "github.com/xxlixin1993/easyGo/configure"
    "github.com/streadway/amqp"
    "math/rand"
    "sync"
    "time"
)

type connections []*amqp.Connection

var rabbitPool *rabbitMqPool

//基于Slice的连接池
type rabbitMqPool struct {
    size        int
    initialized bool
    container   connections
    rw          sync.RWMutex
}

//初始化rabbitMq连接池
func Init(size int) error {
    rabbitPool = &rabbitMqPool{size: size}

    return rabbitPool.initialize()
}

func (p *rabbitMqPool) initialize() error {
    dsn := configure.DefaultString("rabbitMq.dsn", "amqp://guest:guest@localhost:5672")
    for i := 0; i < rabbitPool.size; i++ {
        connection, err := amqp.Dial(dsn)
        if err != nil {
            return  err
        }
        p.put(connection)
    }
    rabbitPool.initialized = true
    return  nil
}

func (p *rabbitMqPool) put(conn *amqp.Connection) int8 {
    p.rw.Lock()
    defer p.rw.Unlock()
    p.container = append(rabbitPool.container, conn)
    return int8(len(p.container) - 1)
}

func (p *rabbitMqPool) removeByPos(pos int8) {
    p.rw.Lock()
    defer p.rw.Unlock()
    p.container = append(p.container[:pos], p.container[pos:]...)
}

// 随机从池中获取一个连接
func (p *rabbitMqPool) getConnection() (*amqp.Connection, int8) {
    rand.Seed(time.Now().Unix())
    index := rand.Intn(p.size - 1)
    p.rw.RLock()
    defer p.rw.RUnlock()
    return p.container[index], int8(index)
}
